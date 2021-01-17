package handler

import (
	"net/http"
	logger "onbio/logger"
	"onbio/model"
	"onbio/utils"
	"onbio/utils/errcode"
	"onbio/utils/htmlparser"
	"onbio/utils/uploader"

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

	isConfirmed := c.GetInt("is_confirmed")

	userName := c.GetString("user_name")

	if isConfirmed == 0 {
		logger.Info("you are not confirmed before")
		c.Error(errcode.ErrEmailNotConfirmed)
		return
	}

	//根据传入的URL拉取对应的信息
	title, desc, img, err := htmlparser.ParseUrl(params.LinkUrl)

	if err != nil {
		logger.Error("ParseUrl failed ", zap.String("url", params.LinkUrl), zap.Error(err))
		c.Error(errcode.ErrInternal)
		return
	}
	var remoteUrl string
	//根据获取到的URL地址，下载图片，上传到ucloud ,拿到新的地址
	err, picPath := utils.DownloadPiscToTmp(img, userName) //下载到本地服务器

	//如果下载不下来，就还用原来的链接，算盗链
	if err != nil {
		logger.Error("download pic failed ", zap.Error(err))
		remoteUrl = img

	} else {
		logger.Info("download succ ", zap.String("picPath", picPath))

		//再上传到ucloud
		remoteUrl, err = uploader.UploadIMGToUcloud(picPath)

		if err != nil {
			logger.Error("upload to ucloud failed ", zap.Error(err))
		}

	}

	err = model.CreateLink(userID, params.Position, params.LinkUrl, desc, remoteUrl, title)

	if err != nil {
		logger.Error("crate link failed ", zap.Any("params", params), zap.Uint64("userID", userID))
		c.Error(errcode.ErrInternal)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{},
	})

}
