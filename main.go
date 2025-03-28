package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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
	if resp.StatusCode != http.StatusOK {
		return "ERROR: Status request to URL failed with status code: " + strconv.Itoa(resp.StatusCode), err
	}
	defer resp.Body.Close()

	if contentType := resp.Header.Get("Content-Type"); contentType != "text/html" {
		return "ERROR: response Header is not of Content-Type `text/html`", err
	}
	htmlResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "ERROR: Failed to read HTML body", err
	}

	return string(htmlResp), nil
}
