package handler

import (
	"net/http"
	logger "onbio/logger"
	"onbio/model"
	"onbio/utils/errcode"
	"onbio/utils/htmlparser"
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
	if params.LinkUrl == ""{
		logger.Error("params err ")
		c.Error(errcode.ErrParam)
		return
	}

	userID := (uint64)(c.GetInt("user_id"))

	//根据传入的URL拉取对应的信息
	title ,desc,img ,err := htmlparser.ParseUrl(params.LinkUrl)

	if err != nil {
		logger.Error("ParseUrl failed ",zap.String("url",params.LinkUrl),zap.Error(err))
		c.Error(errcode.ErrInternal)
		return 
	}

	err = model.CreateLink(userID,params.Position, params.LinkUrl, desc, img,title)

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
