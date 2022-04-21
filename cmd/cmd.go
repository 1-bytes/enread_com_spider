package cmd

import (
	"context"
	"encoding/json"
	"enread_com/bootstrap"
	"enread_com/parser"
	elasticsearch "enread_com/pkg/elastic"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/olivere/elastic/v7"
	"net"
	"net/http"
	"strconv"
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
	// 请求发起之前要处理的一些事件
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		r.Headers.Set("Referer", "http://www.baidu.com/?from=home")
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")
		r.Headers.Set("Cookie", "security_session_mid_verify=692bca3343d933a33e46033436ac1c65;")
	})

	// 抓取新的页面
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		if strings.Index(url, "http://") != 0 {
			url = "http://www.enread.com" + url
		}

		// 链接中存在某些关键字的直接跳过
		skipURLKeywordsMap := []string{
			"#", "com/tags", "com/skill", "com/exam", "com/job",
		}
		for _, keyword := range skipURLKeywordsMap {
			if strings.Index(url, keyword) != -1 {
				return
			}
		}
		e.Request.Visit(url)
	})

	// 处理请求结果
	c.OnResponse(func(r *colly.Response) {
		url := r.Request.URL.String()
		if strings.Index(url, "index.html") != -1 || strings.Index(url, "html") == -1 {
			return
		}

		body := r.Body
		id := parser.ID(url)
		title := parser.Title(body)
		author := parser.Author()
		category := parser.Category(url)
		releaseDate := parser.ReleaseDate(body)
		paragraphs, err := parser.Content(body)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		for _, paragraph := range paragraphs {
			//fmt.Printf("ID: %d\n", id)
			//fmt.Printf("Title: %s\n", title)
			//fmt.Printf("Author: %s\n", author)
			//fmt.Printf("Category: %s\n", category)
			//fmt.Printf("ReleaseDate: %s\n", releaseDate)
			//fmt.Printf("EN: %s\n", paragraph["EN"])
			//fmt.Printf("CN: %s\n", paragraph["CN"])
			//fmt.Println()

			data := parser.JsonData{
				ID:          strconv.Itoa(id),
				Title:       title,
				Author:      author,
				ReleaseDate: releaseDate,
				Category:    category,
				SourceURL:   url,
				Paragraph:   paragraph,
			}
			if err = SaveDataToElastic("enread_com", "", data); err != nil {
				fmt.Printf("SaveData error: %v\n", err)
			}
		}
		err = SaveDataToMySQL("dict_article_test", &parser.DictArticleModel{
			ID:                  id,
			Title:               title,
			Author:              author,
			ReleaseDate:         releaseDate,
			MostRecentlyUpdated: "",
		})
		if err != nil {
			fmt.Printf("SaveData error: %v\n", err)
		}
	})

	// 错误处理
	c.OnError(func(resp *colly.Response, err error) {
		err = resp.Request.Retry()
		if err != nil {
			fmt.Println("Request URL:", resp.Request.URL, "failed with response:", resp, "\nError:", err)
		}
	})
}

// SaveDataToElastic 存储数据至 ES
func SaveDataToElastic(index string, id string, data interface{}) error {
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

// SaveDataToMySQL 存储数据至 mysql
func SaveDataToMySQL(tables string, data *parser.DictArticleModel) error {
	db := bootstrap.DB
	tx := db.Table(tables).Create(data)
	if err := tx.Error; err != nil {
		return err
	}
	return nil
}
