package htmlconv

import (
	"regexp"
	"strings"
)

var ptn_youtube = regexp.MustCompile("(?:(?:https://youtu.be/)|(?:https://www.youtube.com/watch\\?v=))([^\\s><]*)")
var prefix_youtube = "<div class='yt_container'><iframe src=\"https://www.youtube.com/embed/"
var suffix_youtube = "\" frameborder=\"0\" allowfullscreen class='yt_video'></iframe></div>"

func processYouTubeLink(html string) string {
	return ptn_youtube.ReplaceAllStringFunc(html, func(subs string) string {
		return ptn_youtube.ReplaceAllString(subs, prefix_youtube+"$1"+suffix_youtube)
	})
}

func processImgLink(html string) string {
	html = strings.Replace(html, "<img src=\"img", "<img src=\"/public/img", -1)
	html = strings.Replace(html, "<img src=\"../img", "<img src=\"/public/img", -1)
	return html
}

func PostProcessHTML(html string) string {
	processors := []func(string) string{
		processImgLink,
		processYouTubeLink}
	for _, p := range processors {
		html = p(html)
	}
	return html
}
