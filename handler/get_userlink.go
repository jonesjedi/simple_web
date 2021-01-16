package handler

import (
	"net/http"
	"onbio/logger"
	"onbio/model"
	"onbio/utils/errcode"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GetUserLinkReq struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}

func HandleGetUserLinkRequest(c *gin.Context) {

	//没有参数
	userID := (uint64)(c.GetInt("user_id"))

	var params GetUserLinkReq

	if err := c.ShouldBind(&params); err != nil {
		logger.Info("err params ")
		c.Error(errcode.ErrParam)
		return
	}

	//获取该用户的链接
	linkList, count, err := model.GetUserLinkList(userID, params.Page, params.PageSize)

	if err != nil {
		logger.Error("get user link list failed ", zap.Error(err))
		c.Error(errcode.ErrInternal)
		return
	}
	var dataList []gin.H

	if len(linkList) > 0 {
		for _, link := range linkList {
			dataList = append(dataList, gin.H{
				"id":         link.ID,
				"position":   link.Position,
				"use_flag":   link.UseFlag,
				"is_special": link.IsSpecial,
				"link_url":   link.LinkUrl,
				"link_desc":  link.LinkDesc,
				"link_img":   link.LinkImg,
				"link_title": link.LinkTitle,
			})
		}
	}

	//

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			//"page":      params.Page,
			//"page_size": params.PageSize,
			"count": count,
			"list":  dataList,
		},
	})

}
