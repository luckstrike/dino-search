package crawler

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gocolly/colly/v2"
)

func isURLValid(urlStr *url.URL) error {
	if urlStr.Scheme != "http" && urlStr.Scheme != "https" {
		return fmt.Errorf("only HTTP/HTTPS URLs are allowed")
	}

	return nil
}

func Crawl(userURL string) error {
	parsedURL, err := url.Parse(userURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}

	// Adding a scheme if it's missing
	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "https"
	}

	if err := isURLValid(parsedURL); err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}

	robots := newRobotsChecker()

	c := colly.NewCollector(
		colly.MaxDepth(2), // limits the crawl depth
	)

	// Rate Limiting
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 5 * time.Second,
	})

	c.SetRequestTimeout(10 * time.Second)

	// Before making a request print "Visiting..."
	c.OnRequest(func(r *colly.Request) {
		if !robots.isAllowed(r.URL.Scheme, r.URL.Host, r.URL.Path, c.UserAgent) {
			r.Abort()
			fmt.Printf("Skipping %s - disallowed by robots.txt\n", r.URL.String())
		} else {
			fmt.Println("Visiting", r.URL.String())
		}
	})

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		absoluteLink := e.Request.AbsoluteURL(link)
		if absoluteLink != "" {
			fmt.Printf("Link found: %q -> %s\n", e.Text, absoluteLink)
			e.Request.Visit(absoluteLink)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error visiting %s: %v\n", r.Request.URL, err)
	})

	err = c.Visit(parsedURL.String())
	if err != nil {
		fmt.Printf("Visit error: %v\n", err)
		return err
	}

	return nil
}
