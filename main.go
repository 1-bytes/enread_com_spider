package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func main() {
	a := regexp.MustCompile(`/[A-Za-z]+/\d+\.html$`)
	text := "/tags.php?/work/"
	fmt.Println(a.Match([]byte(text)))

	c := colly.NewCollector(
		//colly.Debugger(&debug.LogDebugger{}),
		colly.Async(false),
		//colly.AllowedDomains("www.enread.com", "enread.com"),
		colly.DetectCharset(),
		colly.URLFilters(
			regexp.MustCompile(`/[A-Za-z/]+/\d+\.html$`),
			regexp.MustCompile(`/[A-Za-z/]+/index\.html$`),
			regexp.MustCompile(`http://www\.enread\.com/$`),
		),
	)

	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		r.Headers.Set("Referer", "http://www.baidu.com/?from=home")
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")
		r.Headers.Set("Cookie", "security_session_mid_verify=692bca3343d933a33e46033436ac1c65;")
	})

	//Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		if strings.Index(url, "http://") != 0 {
			url = "http://www.enread.com" + url
		}

		// 链接中存在某些关键字的直接跳过
		skipURLKeywordsMap := []string{
			"#", "skill",
		}
		for _, keyword := range skipURLKeywordsMap {
			if strings.Index(url, keyword) != -1 {
				return
			}
		}
		e.Request.Visit(url)
	})

	c.OnError(func(resp *colly.Response, err error) {
		fmt.Println("Request URL:", resp.Request.URL, "failed with response:", resp, "\nError:", err)
	})
	c.Visit("http://www.enread.com/")

	//bootstrap.Setup()
	//paragraphs, err := parser.FetchAndContent("http://www.enread.com/science/116639.html")
	////paragraphs, err := parser.FetchAndContent("http://www.enread.com/science/116816.html")
	//if err != nil {
	//	fmt.Printf("Error: %s\n", err)
	//}
	//
	//for _, paragraph := range paragraphs {
	//	fmt.Printf("EN: %s\n", paragraph["EN"])
	//	fmt.Printf("CN: %s\n", paragraph["CN"])
	//	fmt.Println()
	//}
}
