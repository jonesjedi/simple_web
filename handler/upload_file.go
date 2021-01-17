package handler

import (
	"net/http"
	"onbio/logger"
	"onbio/utils/errcode"
	"onbio/utils/uploader"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	TempFilePath = "/tmp/"
)

func HandleUploadFile(c *gin.Context) {

	//获取file
	file, err := c.FormFile("upload")

	if err != nil {
		logger.Error("upload file err ", zap.Error(err))
		c.Error(errcode.ErrUploadFileFailed)
		return
	}

	fileName := file.Filename
	fileOnServer := TempFilePath + fileName
	err = c.SaveUploadedFile(file, fileOnServer)

	if err != nil {
		logger.Error("SaveUploadedFile file err ", zap.Error(err))
		c.Error(errcode.ErrInternal)
		return
	}

	//把文件上传到ucloud
	remoteUrl, err := uploader.UploadIMGToUcloud(fileOnServer)

	if err != nil {
		logger.Error("UploadIMGToUcloud  err ", zap.Error(err))
		c.Error(errcode.ErrInternal)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"url":  remoteUrl,
	})

}
