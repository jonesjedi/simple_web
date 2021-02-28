package htmlparser

import (
	"errors"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
	"net/http"
	"onbio/logger"
	"strings"
)

func ParseUrl(url string) (title, desc, img string, err error) {

	if url == "" {
		err = errors.New("invalid url ")
		return
	}

	res, err := http.Get(url)
	if err != nil {
		logger.Error("http get err", zap.String("url", url), zap.Error(err))
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		logger.Error("http get err", zap.String("url", url), zap.Error(err))
		return
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logger.Error("go query failed ", zap.Error(err))
		return
	}

	title = doc.Find("head").Find("title").Text()

	desc, _ = doc.Find("head").Find("meta[name=description]").Attr("content")

	img, _ = doc.Find("body").Find("img").Attr("src")

	find := strings.Contains(img, "http")
	if !find {
		img = "https:" + img
	}
	//img += "https:"

	logger.Info("parse result", zap.String("title", title), zap.String("desc", desc), zap.String("img", img))

	return
}
