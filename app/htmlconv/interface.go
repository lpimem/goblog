package htmlconv

import (
	"errors"
	"strings"
)

type HtmlConv interface {
	ToHTML(bytes []byte) string
}

func GetConv(ftype string) (HtmlConv, error) {
	ftype = strings.ToLower(ftype)
	switch ftype {
	case ".md":
		return MarkdownConv, nil
	}
	return nil, errors.New("No HtmlConv for file type " + ftype)
}

func init() {

}
