package markdown

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/pingcap/errors"
)

func mdToHTML(md []byte) []byte {
	// create Markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

// GetPreviewContent 截取正文内容 (简介)
func GetPreviewContent(md string, length int) (preview string, err error) {
	htmlContent := mdToHTML([]byte(md))
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlContent))
	if err != nil {
		err = errors.New("goquery解析失败")
		return
	}
	text := doc.Text()
	if len(text) >= length {
		preview = string([]rune(text[:length]))
	}
	return
}
