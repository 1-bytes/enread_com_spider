package parser

import (
	"enread_com/pkg/filters"
	"regexp"
)

// Title 用于解析页面中的标题
func Title(body []byte) string {
	re := `<tbody><tr><td height="72"><div[^>]*?><font[^>]*?>(.+)</font></div></td> </tr>`
	var titleRe = regexp.MustCompile(re)
	title := titleRe.FindAllSubmatch(body, -1)
	if len(title) == 0 {
		return ""
	}
	return string(filters.HtmlFilter(title[0][1]))
}
