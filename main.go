package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

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
	pages := make(map[string]int)
	crawlPage(BASE_URL, BASE_URL, pages)
	fmt.Println(pages)
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

// TODO
// - properly recurse through child URLs
func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	rawBase, err := url.Parse(rawBaseURL)
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

	if _, ok := pages[normCurrentURL]; ok {
		pages[normCurrentURL]++
		return
	} else {
		pages[normCurrentURL] = 1
	}

	currentURLhtml, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("CURRENTLY IN THIS URL'S HTML: ", currentURLhtml)

	htmlReader := strings.NewReader(currentURLhtml)
	treeNode, err := html.Parse(htmlReader)
	if err != nil {
		fmt.Println(err)
		return
	}

	if treeNode.FirstChild != nil && treeNode.LastChild.Data == "a" {
		crawlPage(rawBaseURL, treeNode.FirstChild.Attr[0].Val, pages)
	} else if treeNode.FirstChild != nil && treeNode.LastChild.Data == "a" {
		crawlPage(rawBaseURL, treeNode.LastChild.Attr[0].Val, pages)
	}
}
