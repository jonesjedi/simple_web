package utils

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"onbio/logger"
	"os"
	"path"

	"go.uber.org/zap"
)

const (
	DstPath = "/tmp/"
)

func Md5Sum(src string) (dst string, err error) {

	data := []byte(src)
	has := md5.Sum(data)
	dst = fmt.Sprintf("%x", has) //将[]byte转成16进制

	return
}

func DownloadPiscToTmp(imgUrl string, userName string) (err error, picPath string) {

	md5Str, _ := Md5Sum(path.Base(imgUrl))

	fileName := userName + "_" + md5Str

	res, err := http.Get(imgUrl)
	if err != nil {
		logger.Error("Get url failed ", zap.Error(err))
		return
	}
	defer res.Body.Close()
	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32*1024)
	picPath = DstPath + fileName
	file, err := os.Create(picPath)
	if err != nil {
		logger.Error("Create failed", zap.Error(err))
		return
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)

	_, _ = io.Copy(writer, reader)

	return

}
