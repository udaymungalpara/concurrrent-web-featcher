package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Result struct {
	url      string
	status   int
	duration time.Duration
	Title    string
}

func fetchurl(url string, wg *sync.WaitGroup, res chan<- Result) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil || resp == nil {
		res <- Result{url: url, Title: "Error"}
		return
	}
	defer resp.Body.Close()
	start := time.Now()

	if err != nil {

	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)

	var title string

	if err != nil {
		title = "title not found"
	} else {

		title = doc.Find("title").Text()

	}

	res <- Result{url: url, status: resp.StatusCode, duration: time.Since(start), Title: title}

}

func main() {

	urls := []string{
		"https://golang.org",
		"https://google.com",
		"https://www.geeksforgeeks.org",
	}

	res := make(chan Result, len(urls))

	var wg sync.WaitGroup

	for _, i := range urls {
		wg.Add(1)
		go fetchurl(i, &wg, res)
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	for i := range res {
		fmt.Println("URL:", i.url)
		fmt.Println("Status:", i.status)
		fmt.Println("Title:", i.Title)
		fmt.Println("Time:", i.duration)
		fmt.Println()

	}

}
