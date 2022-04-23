package fetcher

import (
	"bufio"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"net/http"
)

// Fetch 用于获取网页内容
func Fetch(url string) ([]byte, error) {
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	request.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	request.Header.Set("Cache-Control", "max-age=0")
	request.Header.Set("Cookie", "srcurl=687474703a2f2f7777772e656e726561642e636f6d2f; security_session_mid_verify=ea4e4321c96759befa268f2016f61716; __51cke__=; PHPSESSID=mjc97n2f3rviailkfvkd0ij4h3; yunsuo_session_verify=7d0331068ff89ddb2355f994338f4c9e; __tins__1636281=^%^7B^%^22sid^%^22^%^3A^%^201650614448898^%^2C^%^20^%^22vd^%^22^%^3A^%^208^%^2C^%^20^%^22expires^%^22^%^3A^%^201650617486307^%^7D; __51laig__=15")
	request.Header.Set("DNT", "1")
	request.Header.Set("Host", "www.enread.com")
	request.Header.Set("If-Modified-Since", "Thu, 21 Apr 2022 08:49:22 GMT")
	request.Header.Set("If-None-Match", `W/"62611a92-9947"`)
	request.Header.Set("Referer", "http://www.enread.com/science/116639.html?security_verify_data=313730372c393630")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")

	resp, err := client.Do(request)
	bodyReader := bufio.NewReader(resp.Body)
	e := determineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}

// determineEncoding 自动判断网页编码.
func determineEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		log.Printf("Fetcher error: %v", err)
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
