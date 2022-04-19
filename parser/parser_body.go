package parser

import (
	"regexp"
)

var bodyRe = regexp.MustCompile(`<div id="dede_content">([\S\s]+)<div class="dede_pages">`)

// ContentFromBody 从网页中提取正文
func ContentFromBody(body []byte) []byte {
	content := bodyRe.FindAllSubmatch(body, -1)
	if len(content) == 0 {
		return nil
	}
	return content[0][1]
}
