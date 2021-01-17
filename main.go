package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/gocolly/colly"
)

type House struct {
	location		string
	price			string
}

type Olx struct {
	scraper *colly.Collector
	entries map[string]*House
}

func NewOlx() *Olx {

	c := colly.NewCollector(
		colly.URLFilters(
			regexp.MustCompile("http[s]://www.olx.pt/.*$"),
		),
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
		make(map[string]*House),
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
		url := e.Request.URL.String()

		price := goquery.Find("strong").Text()
		
		if e, ok := o.entries[url]; ok {
			e.price = price
		} else {
			newEntry := House{}
			newEntry.price = price
			o.entries[url] = &newEntry
		}


	})

	o.scraper.OnHTML("div[class='offer-user__location']", func (e *colly.HTMLElement) {
		location := e.DOM.Children().Find("p").Text()
		url := e.Request.URL.String()

		if entry, ok := o.entries[url]; ok {
			entry.location = location
		} else {
			newEntry := House{}
			newEntry.location = location
			o.entries[url] = &newEntry
		}
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
