package xss

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

// XssFilter 防止xss注入
func XssFilter(content string) (newContent string) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return
	}
	doc.Find("script").Remove()
	doc.Find("img").Remove()
	doc.Find("iframe").Remove()
	newContent = doc.Text()
	return
}
