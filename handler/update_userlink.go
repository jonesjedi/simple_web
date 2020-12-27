package handler

import (
	"net/http"
	logger "onbio/logger"
	"onbio/model"
	"onbio/utils/errcode"

	"github.com/gin-gonic/gin"
)

type UpdateUserLinkParam struct {
	ID       int    `json:"id" binding:"required"`
	LinkUrl  string `json:"link_url"`
	LinkImg  string `json:"link_img"`
	LinkDesc string `json:"link_desc"`
}

func HandleUpdateUserLinkRequest(c *gin.Context) {

	var params UpdateUserLinkParam
	err := c.Bind(&params)
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

	userID := (uint64)(c.GetInt("user_id"))

	var link model.Link
	if params.LinkDesc != "" {
		link.LinkDesc = params.LinkDesc
	}
	if params.LinkUrl != "" {
		link.LinkUrl = params.LinkUrl
	}
	if params.LinkImg != "" {
		link.LinkImg = params.LinkImg
	}
	linkID := (uint64)(params.ID)
	err = model.UpdateLinkByID(linkID, userID, link)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{},
	})

}
