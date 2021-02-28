package handler

import (
	"net/http"
	logger "onbio/logger"
	"onbio/model"
	"onbio/utils"
	"onbio/utils/errcode"
	"onbio/utils/htmlparser2"
	"onbio/utils/uploader"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CreateUserLinkParam struct {
	LinkUrl  string `json:"link_url" binding:"required"`
	Position uint64 `json:"position" binding:"required"`
	//LinkImg  string `json:"link_img" binding:"required"`
	//LinkDesc string `json:"link_desc" binding:"required"`
}

func HandleCreateUserLinkRequest(c *gin.Context) {

	var params CreateUserLinkParam
	err := c.Bind(&params)
	if err != nil {
		logger.Error("params err ")
		c.Error(errcode.ErrParam)
		return
	}
	if params.LinkUrl == "" {
		logger.Error("params err ")
		c.Error(errcode.ErrParam)
		return
	}

	userID := (uint64)(c.GetInt("user_id"))

	//isConfirmed := c.GetInt("is_confirmed")

	userName := c.GetString("user_name")
	/** 没确认可以创建链接
	if isConfirmed == 0 {
		logger.Info("you are not confirmed before")
		c.Error(errcode.ErrEmailNotConfirmed)
		return
	}*/
	linkUrl := params.LinkUrl
	//判断下前缀，如果没传http或者https，就默认是https
	if !strings.Contains(linkUrl, "http") && !strings.Contains(linkUrl, "tel") &&
		!strings.Contains(linkUrl, "sms") && !strings.Contains(linkUrl, "mailto") {
		linkUrl = "https://" + linkUrl
	}

	//根据传入的URL拉取对应的信息
	title, desc, img, err := htmlparser2.ParseURL(linkUrl)

	if err != nil {
		logger.Error("ParseUrl failed ", zap.String("url", linkUrl), zap.Error(err))
		c.Error(errcode.ErrInternal)
		return
	}
	var remoteUrl string
	//根据获取到的URL地址，下载图片，上传到ucloud ,拿到新的地址
	err, picPath := utils.DownloadPiscToTmp(img, userName) //下载到本地服务器

	//如果下载不下来，就还用原来的链接，算盗链
	if err != nil {
		logger.Error("download pic failed ", zap.Error(err))
		remoteUrl = ""

	} else {
		logger.Info("download succ ", zap.String("picPath", picPath))

		for {
			//判断文件大小
			fileOnTmp, err := os.Stat(picPath)
			if err != nil {
				logger.Info("err stat file ", zap.String("file path", picPath))
				remoteUrl = ""
				break
			}
			fileSize := fileOnTmp.Size()

			if fileSize < 2000 {
				logger.Info("file is too small as a pic", zap.Int64("fileSize", fileSize))
				remoteUrl = ""
				break
			}
			//再上传到ucloud,超过5kb的才上传
			if fileSize >= 2000 {
				remoteUrl, err = uploader.UploadIMGToUcloud(picPath)

				if err != nil {
					logger.Error("upload to ucloud failed ", zap.Error(err))
				}
			}
		}

	}

	id, err := model.CreateLink(userID, params.Position, linkUrl, desc, remoteUrl, title)

	if err != nil {
		logger.Error("create link failed ", zap.Any("params", params), zap.Uint64("userID", userID))
		c.Error(errcode.ErrInternal)
		return
	}

	//把刚写入的记录查出来，作为回显
	link, err := model.GetUserLinkByID(id)
	if err != nil {
		logger.Error("create link failed ", zap.Any("id", id))
		c.Error(errcode.ErrInternal)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"id":         link.ID,
			"position":   link.Position,
			"use_flag":   link.UseFlag,
			"is_special": link.IsSpecial,
			"link_url":   link.LinkUrl,
			"link_desc":  link.LinkDesc,
			"link_img":   link.LinkImg,
			"link_title": link.LinkTitle,
		},
	})

}
