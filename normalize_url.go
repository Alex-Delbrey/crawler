package main

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func normalizeURL(inputURL string) (string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}

	resultURL := inputURL[len(parsedURL.Scheme)+3:]
	if resultURL[len(resultURL)-1] == '/' {
		resultURL = resultURL[:len(resultURL)-1]
	}

	return resultURL, nil
}

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	var str []string

	rootURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return str, err
	}
	rootScheme := rootURL.Scheme

	htmlReader := strings.NewReader(htmlBody)
	treeNode, err := html.Parse(htmlReader)
	if err != nil {
		return str, err
	}

	for n := range treeNode.Descendants() {
		if n.Data == "a" && len(n.Attr) != 0 && n.Attr[0].Key == "href" {
			if !strings.Contains(n.Attr[0].Val, rootScheme) {
				str = append(str, rawBaseURL+n.Attr[0].Val)
			} else {
				str = append(str, n.Attr[0].Val)
			}
		}
	}

	return str, nil
}
