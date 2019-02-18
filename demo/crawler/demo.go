package main

import (
	"fmt"
	"sync"

	"github.com/gocolly/colly"
)

func main() {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		//colly.AllowedDomains("blog.xiongyingqi.com", "xiongyingqi.com"),
	)
	wg := sync.WaitGroup{}

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

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://blog.xiongyingqi.com/")

	wg.Wait()
}
