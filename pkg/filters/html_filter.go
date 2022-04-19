package filters

import "regexp"

var (
	commonRe = regexp.MustCompile(`<[^>]+>`)
	nbspRe   = regexp.MustCompile(`\s*?&nbsp; \s*?`)
)

// HtmlFilter 用于将信息中的 html 标签过滤掉0
func HtmlFilter(content []byte) []byte {
	content = commonRe.ReplaceAll(content, []byte(""))
	content = nbspRe.ReplaceAll(content, []byte(" "))
	return content
}
