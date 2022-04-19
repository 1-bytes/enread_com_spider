package parser

import (
	"enread_com/pkg/filters"
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

// Content 解析正文
func Content(bytes []byte) []paragraph {
	contents := strings.Split(string(bytes), "\t&nbsp; ")

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
				temp["TranslationCN"] = split[k]
				continue
			}
			temp["English"] = split[k]
		}
		paragraphs = append(paragraphs, temp)
	}
	return paragraphs
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
