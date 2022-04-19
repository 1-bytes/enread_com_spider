package main

import (
	"enread_com/parser"
	"enread_com/pkg/fetcher"
	"enread_com/pkg/filters"
	"fmt"
)

func main() {
	//c := colly.NewCollector()
	//
	//// Find and visit all links
	//c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	//	e.Request.Visit(e.Attr("href"))
	//})
	//
	//c.OnRequest(func(r *colly.Request) {
	//	fmt.Println("Visiting", r.URL)
	//})
	//
	//c.Visit("http://go-colly.org/")
	//bytes, err := fetcher.Fetch("http://www.enread.com/science/116816.html")
	bytes, err := fetcher.Fetch("http://www.enread.com/science/116639.html")
	if err != nil {
		panic(err)
	}
	bytes = parser.ContentFromBody(bytes)
	bytes = filters.HtmlFilter(bytes)
	fmt.Println(string(bytes))
}
