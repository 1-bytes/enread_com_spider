package filters

import (
	"regexp"
)

var (
	reMap = map[string]string{
		`<[^>]+>`: ` `,
		` {2,10}`: ` `,
		`&amp;`:   `&`,
		`&lt;`:    `<`,
		`&gt;`:    `>`,
		//`&nbsp;`:   ` `,
		`&quot;`:   `"`,
		`&#39;`:    `'`,
		`&middot;`: `·`,
		`&hellip;`: `…`,
		`&mdash;`:  `—`,
		`&ldquo;`:  `“`,
		`&rdquo;`:  `”`,
		`&lsquo;`:  `‘`,
		`&rsquo;`:  `’`,
		`&laquo;`:  `«`,
		`&raquo;`:  `»`,
		`&#8212;`:  `—`,
		`&#8221;`:  `”`,
		`&#8217;`:  `’`,
		`&#8220;`:  `“`,
		`&#8211;`:  `–`,
		`&#8216;`:  `‘`,
	}
)

// HtmlFilter 用于将信息中的 html 标签过滤掉
func HtmlFilter(content string) string {
	for re, value := range reMap {
		reg := regexp.MustCompile(re)
		content = string(reg.ReplaceAll([]byte(content), []byte(value)))
	}
	return content
}
