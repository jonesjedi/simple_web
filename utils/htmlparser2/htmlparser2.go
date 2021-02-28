package htmlparser2

import (
	"errors"
	"onbio/logger"
	"onbio/utils/goscraper"

	"go.uber.org/zap"
)

// ParseURL 使用goscraper获取链接预览
func ParseURL(uri string) (title, desc, img string, err error) {
	if uri == "" {
		err = errors.New("empty url ")
		return
	}
	// user go scraper get url link preview
	s, serr := goscraper.Scrape(uri, 5)
	if serr != nil {
		logger.Info("parse failed", zap.Error(serr))
	}
	title = s.Preview.Title
	desc = s.Preview.Description
	if len(s.Preview.Images) != 0 {
		img = s.Preview.Images[0]
	} else {
		img = s.Preview.Icon
	}
	logger.Info("parse result", zap.String("title", title), zap.String("desc", desc), zap.String("img", img))

	return
}
