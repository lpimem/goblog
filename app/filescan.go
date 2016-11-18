package app

import (
	"goblog/app/models"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"goblog/app/htmlconv"

	"github.com/bradfitz/slice"
	"github.com/revel/revel"
)

var reloadDelay = 5000
var ArticleCache map[string]*models.ArticleInfo
var ArticleList []models.ArticleInfo

func LoadContent(a *models.ArticleInfo) error {
	ftype := filepath.Ext(a.Path)
	conv, err := htmlconv.GetConv(ftype)
	if err != nil {
		return err
	}
	article, err := ioutil.ReadFile(a.Path)
	if err != nil {
		return err
	}
	a.Content = article
	rawHtml := conv.ToHTML(article)
	a.HTMLContent = htmlconv.PostProcessHTML(rawHtml)
	return nil
}

func reloadArticleList() {
	for {
		// revel.INFO.Printf("Reloading article list ...")
		var l []models.ArticleInfo = make([]models.ArticleInfo, 0, 1000)
		var infoCacheNew = make(map[string]*models.ArticleInfo)
		err := filepath.Walk(DocBaseDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				revel.ERROR.Printf("Cannot access file %s, %s", path, err)
				return nil
			}
			if info.IsDir() {
				return nil
			}
			name := strings.Split(info.Name(), ".")[0]
			inCache := ArticleCache[name]
			if inCache == nil || inCache.MTime.Before(info.ModTime()) {
				article := models.ArticleInfo{Title: name, Path: path, MTime: info.ModTime()}
				err := LoadContent(&article)
				if err != nil {
					revel.ERROR.Println("Cannot load article content from ", article.Path, "Reason:", err)
				}
				infoCacheNew[name] = &article
			} else {
				infoCacheNew[name] = inCache
			}
			l = append(l, *infoCacheNew[name])
			return nil
		})
		if err != nil {
			revel.ERROR.Printf("Error reloading article list: %s", err)
		} else {
			ArticleCache = infoCacheNew
			slice.Sort(l, func(i, j int) bool {
				return l[i].MTime.After(l[j].MTime)
			})
			ArticleList = l
			// revel.INFO.Printf("Article list reloaded.")
		}
		time.Sleep(time.Duration(reloadDelay) * time.Millisecond)
	}
}

func reloadArticlesRoutine() {
	go reloadArticleList()
}
