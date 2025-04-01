package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no website provided")
		return
	} else if len(os.Args) > 2 {
		fmt.Println("too many arguments provided")
		return
	} else {
		fmt.Println("starting crawler of:", os.Args[1])
	}

	BASE_URL := os.Args[1]
	cfg := config{
		pages:   make(map[string]int),
		baseURL: &url.URL{Host: BASE_URL},
		mu:      &sync.Mutex{},
		wg:      &sync.WaitGroup{},
	}
	cfg.crawlPage(cfg.baseURL.Host)
	fmt.Println(cfg.pages)
}

func getHTML(rawURL string) (string, error) {
	resp, err := http.Get(rawURL)
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

// code is getting stuck because WaitGroup.wait is still waiting on WaitGroup to be zero.
// it is not zero ecause I add to wg before goroutine yet I only wg.Done when a url exists
// I need to implement wg.Done function in line 92, tricky part is that first call from
// function can not, do wg.Done because it is zero in the first call
func (cfg *config) crawlPage(rawCurrentURL string) {
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

	if _, ok := cfg.pages[normCurrentURL]; ok {
		cfg.mu.Lock()
		cfg.pages[normCurrentURL]++
		cfg.wg.Done()
		cfg.mu.Unlock()
		return
	} else {
		cfg.pages[normCurrentURL] = 1
	}

	currentURLhtml, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("CURRENTLY IN THIS URL'S HTML: ", currentURLhtml)

	allURLinCurrent, err := getURLsFromHTML(currentURLhtml, cfg.baseURL.Host)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, valUrl := range allURLinCurrent {
		fmt.Println("url to traverse: ", valUrl)
		cfg.wg.Add(1)
		go cfg.crawlPage(valUrl)
	}
	cfg.wg.Wait()
}
