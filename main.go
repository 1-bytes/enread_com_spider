package main

import (
	"enread_com/bootstrap"
	"enread_com/parser"
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

	bootstrap.Setup()
	//paragraphs := parser.FetchAndContent("http://www.enread.com/science/116639.html")
	paragraphs, err := parser.FetchAndContent("http://www.enread.com/science/116816.html")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	for _, paragraph := range paragraphs {
		fmt.Printf("EN: %s\n", paragraph["EN"])
		fmt.Printf("CN: %s\n", paragraph["CN"])
		fmt.Println()
	}
}
