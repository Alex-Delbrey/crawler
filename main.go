package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"sync"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no website provided")
		return
	} else if len(os.Args) > 4 {
		fmt.Println("too many arguments provided")
		return
	} else {
		fmt.Println("starting crawler of:", os.Args[1])
	}

	BASE_URL := os.Args[1]
	MAX_CONCURRENCY, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	MAX_PAGES, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}
	cfg := config{
		pages:              make(map[string]int),
		baseURL:            &url.URL{Host: BASE_URL},
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, MAX_CONCURRENCY),
		wg:                 &sync.WaitGroup{},
		maxPages:           MAX_PAGES,
	}
	cfg.crawlPage(cfg.baseURL.Host)
	cfg.wg.Wait()
	fmt.Println(cfg.pages)
	printReport(cfg.pages, cfg.baseURL.Host)
}

func getHTML(rawURL string) (string, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return "ERROR: getting url -" + rawURL, err
	}
	if resp.StatusCode != http.StatusOK {
		return "ERROR: Status request to URL failed with status code: " + strconv.Itoa(resp.StatusCode), err
	}
	defer resp.Body.Close()

	if contentType := resp.Header.Get("Content-Type"); contentType != "text/html; charset=utf-8" {
		return "ERROR: response Header is not of Content-Type `text/html`", err
	}
	htmlResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "ERROR: Failed to read HTML body", err
	}

	return string(htmlResp), nil
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	fmt.Println("url to traverse: ", rawCurrentURL)
	if len(cfg.pages) >= cfg.maxPages {
		fmt.Println("REACHED MAX PAGES")
		return
	}
	rawBase, err := url.Parse(cfg.baseURL.Host)
	if err != nil {
		fmt.Println("URL stdlib was not able to parse rawBaseURL: ", err)
		return
	}
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Println("URL stdlib was not able to parse rawCurrentURL: ", err)
		return
	}

	if rawBase.Host != currentURL.Host {
		return
	}

	normCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
	}

	cfg.mu.Lock()
	if cfg.addPageVisit(normCurrentURL) {
		cfg.pages[normCurrentURL]++
		cfg.mu.Unlock()
		return
	} else {
		cfg.pages[normCurrentURL] = 1
		cfg.mu.Unlock()
	}

	currentURLhtml, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	allURLinCurrent, err := getURLsFromHTML(currentURLhtml, cfg.baseURL.Host)
	if err != nil {
		fmt.Println(err)
		return
	}
	cfg.wg.Add(1)
	go func() {
		defer cfg.wg.Done()
		cfg.concurrencyControl <- struct{}{}
		for _, valUrl := range allURLinCurrent {
			cfg.crawlPage(valUrl)
		}
		<-cfg.concurrencyControl
	}()
}

func (cfg *config) addPageVisit(normalizedURL string) bool {
	if _, ok := cfg.pages[normalizedURL]; ok {
		return true
	} else {
		return false
	}
}

func printReport(pages map[string]int, baseURL string) {
	fmt.Printf(`
		==============================
		REPORT for %s
		==============================`, baseURL)
	fmt.Println()
	for key, val := range pages {
		fmt.Printf(`Found %d internal links to %v`, val, key)
		fmt.Println()
	}
}
