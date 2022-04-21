package main

import (
	"enread_com/bootstrap"
	"enread_com/cmd"
	"enread_com/parser"
	"enread_com/pkg/fetcher"
	"fmt"
	"github.com/gocolly/colly/v2"
	"regexp"
	"strconv"
)

func main() {
	bootstrap.Setup()
	c := cmd.NewCollector(
		//colly.Debugger(&debug.LogDebugger{}),
		colly.Async(true),
		colly.AllowedDomains("www.enread.com", "enread.com"),
		colly.DetectCharset(),
		colly.URLFilters(
			regexp.MustCompile(`/[A-Za-z/]+/\d+\.html$`),
			regexp.MustCompile(`/[A-Za-z/]+/index\.html$`),
			regexp.MustCompile(`http://www\.enread\.com/$`),
		))
	cmd.SpiderCallbacks(c)
	c.Visit("http://www.enread.com/")
	c.Wait()

	//testCase()
}

// testCase 测试用例
func testCase() {
	bootstrap.Setup()
	//url := "http://www.enread.com/novel/43490.html"
	url := "http://www.enread.com/science/116639.html"
	bytes, err := fetcher.Fetch(url)
	if err != nil {
		panic(err)
	}
	id := parser.ID(url)
	title := parser.Title(bytes)
	author := parser.Author()
	category := parser.Category(url)
	releaseDate := parser.ReleaseDate(bytes)
	paragraphs, err := parser.Content(bytes)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	for _, paragraph := range paragraphs {
		fmt.Printf("ID: %d\n", id)
		fmt.Printf("Title: %s\n", title)
		fmt.Printf("Author: %s\n", author)
		fmt.Printf("Category: %s\n", category)
		fmt.Printf("ReleaseDate: %s\n", releaseDate)
		fmt.Printf("EN: %s\n", paragraph["EN"])
		fmt.Printf("CN: %s\n", paragraph["CN"])
		fmt.Println()

		data := parser.JsonData{
			ID:        strconv.Itoa(id),
			Category:  category,
			SourceURL: url,
			Paragraph: paragraph,
		}
		if err = cmd.SaveDataToElastic("enread_com", "", data); err != nil {
			fmt.Printf("SaveData error: %v\n", err)
		}
	}
}
