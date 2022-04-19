package parser

import (
	"enread_com/pkg/filters"
	"regexp"
	"strings"
	"unicode"
)

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
func Content(bytes []byte) []byte {
	contents := strings.Split(string(bytes), "\t&nbsp; ")

	var result []map[string]string
	for key, value := range contents {
		if value == "" {
			continue
		}
		contents[key] = strings.TrimSpace(value)
		temp := strings.Split(value, "\r\n")
		str := map[string]string{}
		for k, v := range temp {
			temp[k] = strings.TrimSpace(v)
			if temp[k] == "" {
				continue
			}
			if hasChinese(temp[k]) {
				str["TranslationCN"] = temp[k]
				continue
			}
			str["English"] = temp[k]
		}
		result = append(result, str)
	}
	return nil
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
