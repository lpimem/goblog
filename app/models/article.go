package models

import "time"

type ArticleInfo struct {
	Tag         string
	Title       string
	Path        string
	Content     []byte
	HTMLContent string
	MTime       time.Time
}

func (a *ArticleInfo) MTimeRepr() string {
	return a.MTime.Format("01/02/2006")
}
