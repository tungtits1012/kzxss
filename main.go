package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	concurrency := 8
	timeout := 8
	retries := 3
	retrySleep := 1

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			Dial:                (&net.Dialer{Timeout: time.Duration(timeout) * time.Second}).Dial,
			TLSHandshakeTimeout: time.Duration(timeout) * time.Second,
		},
	}

	work := make(chan string)
	wg := &sync.WaitGroup{}

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range work {
				checkURL(url, client, retries, retrySleep)
			}
		}()
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			work <- line
		}
	}
	close(work)
	wg.Wait()
}

func checkURL(urlStr string, client *http.Client, retries, retrySleep int) {
	payload := "kzxss"

	// Parse URL
	parsed, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println("[!] Invalid URL:", urlStr, err)
		return
	}

	// Copy original query
	origQuery := parsed.Query()

	// For each parameter, replace its value with the payload
	for param := range origQuery {
		modifiedQuery := parsed.Query()
		modifiedQuery.Set(param, payload)
		parsed.RawQuery = modifiedQuery.Encode()
		testURL := parsed.String()

		// --- Test GET ---
		for attempt := 0; attempt <= retries; attempt++ {
			req, err := http.NewRequest("GET", testURL, nil)
			if err != nil {
				fmt.Println("[!] GET request creation failed:", testURL, err)
				break
			}
			req.Header.Set("Connection", "close")

			resp, err := client.Do(req)
			if err != nil {
				if attempt < retries {
					time.Sleep(time.Duration(retrySleep) * time.Second)
					continue
				}
				fmt.Println("[!] GET request failed:", testURL, err)
				break
			}

			body, _ := io.ReadAll(io.LimitReader(resp.Body, 1_000_000))
			resp.Body.Close()

			if strings.Contains(string(body), payload) {
				fmt.Printf("[REFLECTION:GET] %s (param: %s)\n", testURL, param)
			}
			break
		}

		// --- Test POST ---
		postData := url.Values{}
		for k := range origQuery {
			if k == param {
				postData.Set(k, payload)
			} else {
				postData.Set(k, origQuery.Get(k))
			}
		}

		for attempt := 0; attempt <= retries; attempt++ {
			req, err := http.NewRequest("POST", parsed.Scheme+"://"+parsed.Host+parsed.Path, strings.NewReader(postData.Encode()))
			if err != nil {
				fmt.Println("[!] POST request creation failed:", urlStr, err)
				break
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("Connection", "close")

			resp, err := client.Do(req)
			if err != nil {
				if attempt < retries {
					time.Sleep(time.Duration(retrySleep) * time.Second)
					continue
				}
				fmt.Println("[!] POST request failed:", urlStr, err)
				break
			}

			body, _ := io.ReadAll(io.LimitReader(resp.Body, 1_000_000))
			resp.Body.Close()

			if strings.Contains(string(body), payload) {
				fmt.Printf("[REFLECTION:POST] %s (param: %s)\n", urlStr, param)
			}
			break
		}
	}
}
