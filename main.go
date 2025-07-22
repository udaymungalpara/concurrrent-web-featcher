package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

type Result struct {
	url      string
	status   int
	duration time.Duration
	Title    string
	attempt  int
}

func retry(url string) (*http.Response, error, int) {

	var resp *http.Response
	var err error
	at := 0
	for i := 0; i < 3; i++ {
		resp, err = http.Get(url)
		at++

		if err == nil {
			return resp, err, at
		}
		time.Sleep(time.Second * 2)

	}
	return nil, err, at
}

func fetchurl(url string, wg *sync.WaitGroup, res chan<- Result) {
	defer wg.Done()
	start := time.Now()
	var at int
	resp, err, at := retry(url)
	if err != nil {
		res <- Result{url: url, Title: "error with retry", attempt: at}
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	var title string

	title = "title not found"
	if err == nil {
		title = strings.TrimSpace(doc.Find("title").Text())
	}

	res <- Result{url: url, status: resp.StatusCode, duration: time.Since(start), Title: title, attempt: at}

}

func main() {

	urls := []string{
		"https://golang.org",
		"https://google.com",
		"https://www.geeksforgeeks.org",
	}
	pbar := progressbar.NewOptions(len(urls), progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

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
	var result []Result

	for i := range res {
		
		result = append(result, i)
	}

	for _, i := range result {

		fmt.Println("\n-------Result------")

		fmt.Println("\nURL:", i.url)
		fmt.Println("Status:", i.status)
		fmt.Println("Title:", i.Title)
		fmt.Println("Time:", i.duration)
		fmt.Println("Attempt:", i.attempt)

	}

}
