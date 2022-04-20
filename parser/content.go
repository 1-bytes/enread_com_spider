package parser

import (
	"enread_com/cmd"
	"enread_com/pkg/fetcher"
	"enread_com/pkg/filters"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode"
)

type paragraph map[string]string

var bodyRe = regexp.MustCompile(`<div id="dede_content">([\S\s]+)<div class="dede_pages">`)

// GetContentFromBody 从网页中提取正文
func GetContentFromBody(body []byte) []byte {
	content := bodyRe.FindAllSubmatch(body, -1)
	if len(content) == 0 {
		return nil
	}
	return filters.HtmlFilter(content[0][1])
}

// FetchAndContent 解析正文
func FetchAndContent(u string) ([]paragraph, error) {
	bytes, err := fetcher.Fetch(u)
	if err != nil {
		return nil, err
	}
	bytes = GetContentFromBody(bytes)
	contents := strings.Split(string(bytes), "\t&nbsp; ")
	urlParse, err := url.Parse(u)
	// 取分类
	urlPath := strings.Trim(urlParse.Path, "/")
	urlPathSplit := strings.Split(urlPath, "/")
	category := urlPathSplit[0]

	// 取 ID
	urlPath = strings.Trim(urlPathSplit[1], "/")
	articleID := strings.Split(urlPath, ".")[0]

	var paragraphs []paragraph
	for key, value := range contents {
		if value == "" {
			continue
		}
		contents[key] = strings.TrimSpace(value)
		split := strings.Split(value, "\r\n")
		temp := paragraph{}
		for k, v := range split {
			split[k] = strings.TrimSpace(v)
			if split[k] == "" {
				continue
			}
			if hasChinese(split[k]) {
				temp["CN"] = split[k]
				continue
			}
			temp["EN"] = split[k]
		}
		if temp["EN"] == "" {
			continue
		}
		data := JsonData{
			ArticleID: articleID,
			Category:  category,
			SourceURL: u,
			Paragraph: temp,
		}
		if err = cmd.SaveData("enread_com", "", data); err != nil {
			fmt.Printf("SaveData error: %v\n", err)
		}
	}
	return paragraphs, nil
}

// hasChinese 判断字符串中是否包含中文
func hasChinese(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return true
		}
	}
	return false
}
