package htmlconv

import (
	"regexp"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

func makeUGCPolicy() *bluemonday.Policy {
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class").Matching(
		regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
	return p
}

type markdownconv struct {
	policy *bluemonday.Policy
}

func (c markdownconv) ToHTML(bytes []byte) string {
	extensions := 0
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	// extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS

	htmlFlags := 0
	htmlFlags |= blackfriday.HTML_USE_XHTML
	htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_DASHES
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES

	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")
	unsafeHtml := blackfriday.Markdown(bytes, renderer, extensions)
	html := c.policy.SanitizeBytes(unsafeHtml)
	return string(html)
}

var MarkdownConv = markdownconv{policy: makeUGCPolicy()}
