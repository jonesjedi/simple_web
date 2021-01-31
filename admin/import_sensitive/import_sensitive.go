package import_sensitive

import (
	"fmt"
	"onbio/logger"
	"onbio/redis"

	"github.com/tealeg/xlsx"
	"go.uber.org/zap"
)

const (
	USER_SENSITIVE_WORD_KEY = "onbio_user_sensitive_word"
)

//导入敏感词
func ImportSensitiveWordFromExcel(excelPath string) (err error) {
	xlFile, err := xlsx.OpenFile(excelPath)
	if err != nil {
		fmt.Printf("open failed: %s\n", err)
		return
	}

	conn := redis.GetConn("onbio")
	defer conn.Close()

	key := USER_SENSITIVE_WORD_KEY

	for _, sheet := range xlFile.Sheets {
		fmt.Printf("Sheet Name: %s\n", sheet.Name)
		for _, row := range sheet.Rows {
			for _, cell := range row.Cells {
				text := cell.String()
				//直接塞到redis 的set里
				_, err = conn.Do("sadd", key, text)
				if err != nil {
					logger.Error("err set redis ", zap.String("key", key), zap.Error(err))
					continue
				}
			}
		}
	}
	return
}
