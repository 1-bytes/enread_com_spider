package main

import (
	"enread_com/bootstrap"
	"enread_com/parser"
	"fmt"
)

func main() {
	//c := cmd.NewCollector(
	//	//colly.Debugger(&debug.LogDebugger{}),
	//	colly.Async(false),
	//	colly.AllowedDomains("www.enread.com", "enread.com"),
	//	colly.DetectCharset(),
	//	colly.URLFilters(
	//		regexp.MustCompile(`/[A-Za-z/]+/\d+\.html$`),
	//		regexp.MustCompile(`/[A-Za-z/]+/index\.html$`),
	//		regexp.MustCompile(`http://www\.enread\.com/$`),
	//	))
	//cmd.SpiderCallbacks(c)
	//c.Visit("http://www.enread.com/")

	bootstrap.Setup()
	paragraphs, err := parser.FetchAndContent("http://www.enread.com/science/116639.html")
	//paragraphs, err := parser.FetchAndContent("http://www.enread.com/science/116816.html")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	for _, paragraph := range paragraphs {
		fmt.Printf("EN: %s\n", paragraph["EN"])
		fmt.Printf("CN: %s\n", paragraph["CN"])
		fmt.Println()
	}
}
