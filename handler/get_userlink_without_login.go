package handler

import (
	"net/http"
	"onbio/logger"
	"onbio/model"
	"onbio/utils/errcode"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HandleGetUserLinkWithoutLoginRequest(c *gin.Context) {

	userName := c.DefaultQuery("user_name", "")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if pageSize > 100 {
		pageSize = 10
	}
	if page <= 0 {
		page = 1
	}

	//处理参数
	err, user := model.GetUserInfo("", userName, 0)

	if err != nil {
		logger.Info("get user info by username failed ", zap.String("user name ", userName))
		c.Error(errcode.ErrEmail)
		return
	}

	userID := user.ID

	linkList, count, err := model.GetUserLinkListWithPage(userID, page, pageSize)

	if err != nil {
		logger.Error("get user link list failed ", zap.Error(err))
		c.Error(errcode.ErrInternal)
		return
	}
	var dataList []gin.H

	if len(linkList) > 0 {
		for _, link := range linkList {
			dataList = append(dataList, gin.H{
				"id":       link.ID,
				"position": link.Position,
				//"use_flag":   link.UseFlag,
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
			"page":      page,
			"page_size": pageSize,
			"count":     count,
			"list":      dataList,
		},
	})
}
