package parser

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// ID 解析 ID
func ID(u string) int {
	idRe := regexp.MustCompile(`/(\d+)\.html`)
	match := idRe.FindAllStringSubmatch(u, -1)
	if len(match) > 0 {
		id, _ := strconv.Atoi(match[0][1])
		return id
	}
	return -1
}

// Author 作者
func Author() string {
	return "EnzoLiu"
}

// Category 解析分类
func Category(u string) string {
	urlParse, err := url.Parse(u)
	if err != nil {
		return ""
	}
	// 取分类
	urlPath := strings.Trim(urlParse.Path, "/")
	urlPathSplit := strings.Split(urlPath, "/")
	return strings.TrimSpace(urlPathSplit[0])
}

// ReleaseDate 发布日期
func ReleaseDate(body []byte) string {
	re := regexp.MustCompile(`发布时间：(.+)\s{1,4}字体`)
	match := re.FindAllStringSubmatch(string(body), -1)
	if len(match) == 0 {
		return ""
	}
	return strings.TrimSpace(match[0][1])
}
