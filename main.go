package main

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

func main() {

	c := colly.NewCollector(
	//colly.Debugger(&debug.LogDebugger{}),
	)

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Opening", r.URL.String())
	})

	c.OnHTML("td[class='title-cell ']", func(e *colly.HTMLElement) {
		goquery := e.DOM

		attr, _ := goquery.Children().Children().Find("a").Attr("href")

		e.Request.Visit(attr)
	})

	c.OnHTML("div[class='pricelabel']", func(e *colly.HTMLElement) {
		goquery := e.DOM

		attr := goquery.Find("strong").Text()
		fmt.Println(attr)
	})

	c.Visit("https://www.olx.pt/imoveis/")

}
