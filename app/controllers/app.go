package controllers

import (
	"errors"
	"goblog/app"
	"goblog/app/routes"

	"strings"

	"github.com/revel/revel"
)

var siteTitle string

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	articles := app.ArticleList
	visitorCount := app.VisitorCount
	return c.Render(articles, visitorCount)
}

func (c App) Doc(query string) revel.Result {
	var (
		err      error
		mtime    string
		html     string
		pageUrl  string
		title    string
		titleUrl string
		tag      string
	)
	_, query = app.SplitNameTag(query)
	titleUrl = app.RemoveSpace(query)
	if app.ArticleCache[titleUrl] == nil {
		err = errors.New("article '" + query + "' not found.")
		title = "Opps"
	} else {
		articleInfo := app.ArticleCache[titleUrl]
		title = articleInfo.Title
		html = articleInfo.HTMLContent
		mtime = articleInfo.MTimeRepr()
		tag = articleInfo.Tag
	}
	if err != nil {
		html = err.Error()
	} else {
		app.ReaderCounts[title] += 1
		pageUrl = routes.App.Doc(titleUrl)
	}
	visitorCount := app.ReaderCounts[title]
	return c.Render(title, html, mtime, visitorCount, pageUrl, tag)
}

type Page struct {
	Title string
	Url   string
}

func init() {
	revel.InterceptFunc(app.RecordVisit, revel.BEFORE, &App{})
	revel.TemplateFuncs["SiteTitle"] = func() string {
		return app.SiteTitle
	}
	revel.TemplateFuncs["Pages"] = func() []Page {
		return []Page{
			{"Journal", "/"},
			{"About", "/doc/About-Me"},
		}
	}
	revel.TemplateFuncs["HasPrefix"] = func(a, b string) bool {
		return strings.HasPrefix(a, b)
	}
	revel.TemplateFuncs["RemoveSpace"] = app.RemoveSpace
}
