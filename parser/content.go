package parser

import (
	"enread_com/cmd"
	"enread_com/pkg/fetcher"
	"enread_com/pkg/filters"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type paragraph map[string]string

var bodyRe = regexp.MustCompile(`<div id="dede_content">([\S\s]+)<div class="dede_pages">`)
var titleRe = regexp.MustCompile(`<tbody><tr><td height="72"><div[^>]*?><font[^>]*?>(.+)</font></div></td> </tr>`)

// GetTitleFromBody 从网页中提取标题
func GetTitleFromBody(body string) string {
	title := titleRe.FindAllStringSubmatch(body, -1)
	if len(title) == 0 {
		return ""
	}
	return filters.HtmlFilter(title[0][1])
}

// GetContentFromBody 从网页中提取正文
func GetContentFromBody(body string) string {
	content := bodyRe.FindAllStringSubmatch(body, -1)
	if len(content) == 0 {
		return ""
	}
	return filters.HtmlFilter(content[0][1])
}

// FetchAndContent 解析正文
func FetchAndContent(u string) ([]paragraph, error) {
	// 发起请求 获取网页内容
	bytes, err := fetcher.Fetch(u)
	if err != nil {
		return nil, err
	}
	body := string(bytes)
	title := GetTitleFromBody(body)
	fmt.Println(title)
	// 从完整的网页中获取文章内容
	content := GetContentFromBody(body)
	// 文章拆段落(此时段落中可能中英文混合)
	contents := strings.Split(content, "\t&nbsp; ")
	urlParse, err := url.Parse(u)
	// 取分类
	urlPath := strings.Trim(urlParse.Path, "/")
	urlPathSplit := strings.Split(urlPath, "/")
	category := urlPathSplit[0]
	// 取 ID
	urlPath = strings.Trim(urlPathSplit[1], "/")
	articleIDSplit := strings.Split(urlPath, ".")[0]

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
			// 判断当前段落是中文还是英文
			if hasChinese(split[k]) {
				temp["CN"] = split[k]
				continue
			}
			temp["EN"] = split[k]
		}
		// 如果这个段落的英文是空的，那也没必要存了
		if temp["EN"] == "" {
			continue
		}
		articleID, _ := strconv.Atoi(articleIDSplit)
		data := JsonData{
			ID:        strconv.Itoa(articleID + 1000000),
			Title:     title,
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
