package scraper

import (
	"github.com/gocolly/colly/v2"
	"strings"
)

type SearchableContent struct {
	URL       string
	Title     string
	Text      string   // Main text content
	Headlines []string // h1, h2, h3 headers
	Keywords  []string // meta keywords + important terms
}

type Scraper struct {
	collector *colly.Collector
}

func NewScraper() *Scraper {
	return &Scraper{
		collector: colly.NewCollector(),
	}
}

func (s *Scraper) Scrape(url string) (*SearchableContent, error) {
	content := &SearchableContent{
		URL: url,
	}

	c := s.collector.Clone()

	// We'll store text content from different elements
	var textParts []string
	var headlines []string

	c.OnHTML("html", func(e *colly.HTMLElement) {
		// Get the title
		content.Title = e.ChildText("title")

		// Get meta keywords if available
		e.ForEach("meta[name=keywords]", func(_ int, el *colly.HTMLElement) {
			keywords := strings.Split(el.Attr("content"), ",")
			for _, k := range keywords {
				k = strings.TrimSpace(k)
				if k != "" {
					content.Keywords = append(content.Keywords, k)
				}
			}
		})

		// Get all headings
		e.ForEach("h1, h2, h3", func(_ int, el *colly.HTMLElement) {
			text := strings.TrimSpace(el.Text)
			if text != "" {
				headlines = append(headlines, text)
			}
		})

		// Common content containers
		contentSelectors := []string{
			"article", "main", "[role='main']",
			".content", "#content",
			".post-content", ".entry-content",
			".article", ".post",
			"section",
		}

		// Elements to skip (common noise)
		skipSelectors := []string{
			"header", "footer", "nav",
			".header", ".footer", ".nav",
			".sidebar", ".comments", ".menu",
			".advertisement", ".ads", ".social-share",
			"style", "script", "noscript",
		}

		// First try to find main content area
		for _, selector := range contentSelectors {
			e.ForEach(selector, func(_ int, el *colly.HTMLElement) {
				// Clone element to avoid modifying original
				clone := el.DOM.Clone()

				// Remove noise elements
				for _, skipSelector := range skipSelectors {
					clone.Find(skipSelector).Remove()
				}

				text := strings.TrimSpace(clone.Text())
				if text != "" {
					textParts = append(textParts, text)
				}
			})
		}

		// If no content found, fallback to body with noise removed
		if len(textParts) == 0 {
			clone := e.DOM.Find("body").Clone()
			for _, skipSelector := range skipSelectors {
				clone.Find(skipSelector).Remove()
			}
			text := strings.TrimSpace(clone.Text())
			if text != "" {
				textParts = append(textParts, text)
			}
		}
	})

	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	// Clean and combine all text content
	content.Text = s.cleanText(strings.Join(textParts, " "))
	content.Headlines = headlines

	// If no keywords were found in meta tags, extract important terms
	if len(content.Keywords) == 0 {
		content.Keywords = s.extractKeywords(content.Text)
	}

	return content, nil
}

func (s *Scraper) cleanText(text string) string {
	// Remove extra whitespace
	text = strings.Join(strings.Fields(text), " ")

	// Remove very short lines (likely menu items)
	var validParts []string
	for _, part := range strings.Split(text, ".") {
		if len(part) > 30 { // Arbitrary threshold for meaningful content
			validParts = append(validParts, part)
		}
	}

	return strings.Join(validParts, ". ")
}

func (s *Scraper) extractKeywords(text string) []string {
	// Very basic keyword extraction
	// In a real implementation, you might want to use NLP
	// or TF-IDF to find important terms
	words := strings.Fields(strings.ToLower(text))

	// Remove common words and very short words
	stopwords := map[string]bool{
		"the": true, "and": true, "or": true,
		"a": true, "an": true, "in": true,
		"to": true, "of": true, "for": true,
		// Add more stopwords as needed
	}

	var keywords []string
	seen := make(map[string]bool)

	for _, word := range words {
		if !stopwords[word] && len(word) > 3 && !seen[word] {
			keywords = append(keywords, word)
			seen[word] = true
			if len(keywords) >= 10 { // Limit number of keywords
				break
			}
		}
	}

	return keywords
}

