package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
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
	htmlBody, err := getHTML(BASE_URL)
	if err != nil {
		fmt.Println("error with html")
	}
	fmt.Println(htmlBody)
}

func getHTML(rawURL string) (string, error) {
	resp, err := http.Get(rawURL)
	if resp.StatusCode >= 400 {
		return "", err
	}
	defer resp.Body.Close()

	if contentType := resp.Header.Get("Content-Type"); contentType != "text/html" {
		return "", err
	}
	htmlResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(htmlResp), nil
}
