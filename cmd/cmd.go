package cmd

import (
	"context"
	"encoding/json"
	elasticsearch "enread_com/pkg/elastic"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/olivere/elastic/v7"
	"net"
	"net/http"
	"strings"
	"time"
)

// NewCollector 传入配置信息，创建并返回一个 colly 的 collector 实例
func NewCollector(options ...colly.CollectorOption) *colly.Collector {
	c := colly.NewCollector(options...)
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
	return c
}

// SpiderCallbacks colly 的回调函数
func SpiderCallbacks(c *colly.Collector) {
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
}

// SaveData 存储数据
func SaveData(index string, id string, data interface{}) error {
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var e *elastic.IndexService
	e = elasticsearch.GetInstance().Index()
	if id != "" {
		e.Id(id)
	}
	do, err := e.Index(index).BodyJson(string(j)).Do(context.Background())
	fmt.Printf("%+v: %+v\n", do.Result, do.Id)
	return err
}
