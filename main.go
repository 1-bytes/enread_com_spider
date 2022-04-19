package main

import (
	"enread_com/pkg/filters"
	"fmt"
	"io/ioutil"
	"net/http"
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
	resp, err := http.Get("http://www.enread.com/science/116816.html")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)
	bytes = filters.HtmlFilter(bytes)
	fmt.Println(string(bytes))
}
