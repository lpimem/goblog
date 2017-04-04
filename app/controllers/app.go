package controllers

import (
	"errors"
	"goblog/app"
	"goblog/app/routes"

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

func (c App) Doc(title string) revel.Result {
	var err error = nil
	var mtime string
	var html string
	var pageUrl string
	if app.ArticleCache[title] == nil {
		err = errors.New("article '" + title + "' not found.")
		title = "Opps"
	} else {
		articleInfo := app.ArticleCache[title]
		html = articleInfo.HTMLContent
		mtime = articleInfo.MTimeRepr()
	}
	if err != nil {
		html = err.Error()
	} else {
		app.ReaderCounts[title] += 1
		pageUrl = routes.App.Doc(title)
	}
	visitorCount := app.ReaderCounts[title]
	return c.Render(title, html, mtime, visitorCount, pageUrl)
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
			{"About", "/doc/[Author] About Me"},
			{"Portfolio", "/doc/[Author] Portfolio"},
		}
	}
}
