package parser

import (
	"enread_com/pkg/filters"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type content struct {
	url    string
	title  string
	author string
	date   string
	text   string
}

type paragraph map[string]string

// Content 解析正文
func Content(body []byte) ([]paragraph, error) {
	re := `<div id="dede_content">([\S\s]+)<div class="dede_pages">`
	var contentRe = regexp.MustCompile(re)
	// 从完整的网页中获取文章内容
	content := contentRe.FindAllSubmatch(body, -1)
	if len(content) == 0 {
		return []paragraph{}, fmt.Errorf("failed to parse content")
	}
	body = filters.HtmlFilter(content[0][1])
	// 文章拆段落(此时段落中可能中英文混合)
	contentSplit := strings.Split(string(body), "\r\n")
	var paragraphs []paragraph
	for key, value := range contentSplit {
		if value == "" {
			continue
		}
		temp := paragraph{}
		contentSplit[key] = strings.TrimSpace(value)
		if contentSplit[key] == "" {
			continue
		}
		// 判断当前段落是中文还是英文
		if hasChinese(contentSplit[key]) {
			continue
		}
		temp["EN"] = contentSplit[key]
		paragraphs = append(paragraphs, temp)
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
