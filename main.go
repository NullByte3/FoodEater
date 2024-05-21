package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

func main() {
	file, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	c := colly.NewCollector(colly.Async(true))

	c.OnHTML("article[data-test-id='product-card']", func(e *colly.HTMLElement) {
		productName := e.ChildText("span.sc-facd2606-1")
		unitPrice := cleanPrice(e.ChildText("span[data-test-id='product-price__unitPrice']"))
		comparisonPrice := e.ChildText("div[data-test-id='product-card__productPrice__comparisonPrice']")

		output := fmt.Sprintf("Product: %s, Unit Price: %s, Comparison Price: %s\n", productName, unitPrice, comparisonPrice)
		fmt.Print(output)
		file.WriteString(output)

		productDetailURL := e.Request.AbsoluteURL(e.ChildAttr("a[data-test-id='product-card-link']", "href"))
		c.Visit(productDetailURL)
	})

	c.OnHTML("a[href^='/tuotteet/']", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		visitLink(c, link, e.Request.AbsoluteURL)
	})

	c.OnRequest(func(r *colly.Request) {
		if r.URL.String() == "https://www.s-kaupat.fi/tuotteet" {
			fmt.Println("Visiting main page:", r.URL.String())
			c.Visit("https://www.s-kaupat.fi/tuotteet")
			return
		}
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnHTML("div.sc-162b6339-3", func(e *colly.HTMLElement) {
		productName := e.ChildText("h1[data-test-id='product-name']")
		unitPrice := cleanPrice(e.ChildText("span[data-test-id='product-price__unitPrice']"))
		comparisonPrice := e.ChildText("div[data-test-id='product-page-price__comparisonPrice']")

		output := fmt.Sprintf("Product: %s, Unit Price: %s, Comparison Price: %s\n", productName, unitPrice, comparisonPrice)
		fmt.Print(output)
		file.WriteString(output)
	})

	c.OnHTML("a[href^='/tuote']", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		visitLink(c, link, e.Request.AbsoluteURL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Error:", err)
	})

	c.Visit("https://www.s-kaupat.fi/tuotteet")
	c.Wait()
}

func cleanPrice(price string) string {
	return strings.TrimSpace(strings.TrimPrefix(price, "~"))
}

func visitLink(c *colly.Collector, link string, getAbsoluteURL func(string) string) {
	if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
		fmt.Println("Visiting product page:", link)
		c.Visit(link)
	} else {
		absURL := getAbsoluteURL(link)
		if absURL != "" {
			fmt.Println("Visiting product page:", absURL)
			c.Visit(absURL)
		}
	}
}
