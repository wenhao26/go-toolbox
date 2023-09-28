package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

var (
	targetURL = "https://www.runoob.com/go/go-tutorial.html"
)

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.runoob.com"),
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"),
	)

	c.OnHTML("div #leftcolumn > a", func(e *colly.HTMLElement) {
		title := strings.Replace(e.Text, "\n", "", -1)
		title = strings.Replace(title, "\r", "", -1)
		title = strings.Replace(title, "\t", "", -1)
		title = strings.Replace(title, " ", "", -1)
		href := "https://www.runoob.com" + e.Attr("href")
		fmt.Printf("[title]:%s - [href]:%s \n", title, href)
	})

	err := c.Visit(targetURL)
	if err != nil {
		panic(err)
	}
	c.Wait()
}
