package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	c := colly.NewCollector()
	c.WithTransport(t)

	pages := []string{"https://blog.xiongyingqi.com"}


	c.OnHTML("h1", func(e *colly.HTMLElement) {
		pages = append(pages, e.Text)
	})

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		wg.Add(1)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		go func() {
			c.Visit(e.Request.AbsoluteURL(link))
			wg.Done()
		}()
	})

	c.OnHTML("a", func(e *colly.HTMLElement) {
		c.Visit("file://" + dir + "/html" + e.Attr("href"))
	})

	fmt.Println("file://" + dir + "/html/index.html")
	c.Visit("file://" + dir + "/html/index.html")
	c.Wait()
	for i, p := range pages {
		fmt.Printf("%d : %s\n", i, p)
	}
	wg.Wait()
}

