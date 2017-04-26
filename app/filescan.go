package app

import (
	"goblog/app/models"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"goblog/app/htmlconv"

	"regexp"

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

var _RE_NAME_TAG_SPLIT = regexp.MustCompile(`\[(.*)\][\s\t]*(.*)`)

func SplitNameTag(raw string) (tag, name string) {
	matched := _RE_NAME_TAG_SPLIT.FindStringSubmatch(raw)
	if len(matched) > 0 {
		tag = matched[1]
		name = matched[2]
	} else {
		name = strings.Trim(raw, " \t")
		tag = "random"
	}
	return
}

func RemoveSpace(raw string) string {
	return strings.Replace(raw, " ", "-", -1)
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
			var (
				name string
				tag  string
				url  string
			)
			name = strings.Split(info.Name(), ".")[0]
			tag, name = SplitNameTag(name)
			url = RemoveSpace(name)
			inCache := ArticleCache[url]
			if inCache == nil || inCache.MTime.Before(info.ModTime()) {
				article := models.ArticleInfo{Tag: tag, Title: name, Path: path, MTime: info.ModTime()}
				err := LoadContent(&article)
				if err != nil {
					revel.ERROR.Println("Cannot load article content from ", article.Path, "Reason:", err)
				}
				infoCacheNew[url] = &article
			} else {
				infoCacheNew[url] = inCache
			}
			l = append(l, *infoCacheNew[url])
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
		}
		time.Sleep(time.Duration(reloadDelay) * time.Millisecond)
	}
}

func reloadArticlesRoutine() {
	go reloadArticleList()
}
