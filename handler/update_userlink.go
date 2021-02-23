package handler

import (
	"net/http"
	logger "onbio/logger"
	"onbio/model"
	"onbio/utils/errcode"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UpdateUserLinkParam struct {
	ID        int    `json:"id" binding:"required"`
	LinkUrl   string `json:"link_url"`
	LinkImg   string `json:"link_img"`
	LinkDesc  string `json:"link_desc"`
	Position  int    `json:"position"`
	Title     string `json:"title"`
	UseFlag   int    `form:"use_flag,default=-1" json:"use_flag,default=-1"`
	IsSpecial int    `form:"is_special,default=-1" json:"is_special,default=-1"`
}

func HandleUpdateUserLinkRequest(c *gin.Context) {

	var params UpdateUserLinkParam
	params.UseFlag = -1
	params.IsSpecial = -1
	params.Title = "default"
	params.LinkDesc = "default"
	params.LinkImg = "default"
	err := c.ShouldBindJSON(&params)
	if err != nil {
		logger.Error("params err ")
		c.Error(errcode.ErrParam)
		return
	}
	if params.ID == 0 {
		logger.Error("params err ")
		c.Error(errcode.ErrParam)
		return
	}
	logger.Info("params", zap.Any("params", params))
	userID := (uint64)(c.GetInt("user_id"))

	var link model.Link
	//if params.LinkDesc != "default" {
	link.LinkDesc = params.LinkDesc
	//}
	if params.LinkUrl != "" {
		link.LinkUrl = params.LinkUrl
	}
	//if params.LinkImg != "default" {
	link.LinkImg = params.LinkImg
	//}
	//if params.Title != "default" {
	link.LinkTitle = params.Title
	//}
	if params.Position != 0 {
		link.Position = (uint64)(params.Position)
	}
	//if params.IsSpecial != -1 {
	link.IsSpecial = params.IsSpecial
	//}
	//if params.UseFlag != -1 {
	link.UseFlag = params.UseFlag
	//}
	logger.Info("link info", zap.Any("link", link))
	linkID := (uint64)(params.ID)
	err = model.UpdateLinkByID(linkID, userID, link)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{},
	})

}
