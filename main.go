package main

import (
	"log"
	"fmt"

	"github.com/gocolly/colly"
)

type House struct {
	price string
}

type Olx struct {
	scraper *colly.Collector
	entries map[string]House
}

func NewOlx() *Olx {

	c := colly.NewCollector(
	//colly.Debugger(&debug.LogDebugger{}),
	)

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Opening", r.URL.String())
	})

	return &Olx{
		c,
		make(map[string]House),
	}
}

func (o *Olx) Run(url string) {

	o.scraper.OnHTML("td[class='title-cell ']", func(e *colly.HTMLElement) {
		goquery := e.DOM

		attr, _ := goquery.Children().Children().Find("a").Attr("href")

		e.Request.Visit(attr)
	})

	o.scraper.OnHTML("div[class='pricelabel']", func(e *colly.HTMLElement) {
		goquery := e.DOM

		attr := goquery.Find("strong").Text()

		o.entries[e.Request.URL.String()] = House{attr}

	})

	o.scraper.Visit(url)

}
func main() {

	olx := NewOlx()
	url := "https://www.olx.pt/imoveis/"

	olx.Run(url)

	for key, value := range olx.entries {
		fmt.Println("Key:", key, "Value:", value)
	}

}
