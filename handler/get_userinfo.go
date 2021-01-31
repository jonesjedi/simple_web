package handler

import (
	"net/http"
	logger "onbio/logger"
	"onbio/model"
	"onbio/utils/errcode"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HandleGetUserInfoRequest(c *gin.Context) {

	//没有参数
	//userID := c.GetInt("user_id")
	userName := c.DefaultQuery("user_name", "")
	//根据用户ID获取用户信息

	err, user := model.GetUserInfo("", userName, 0)

	if err != nil {
		logger.Info("get user info by username failed ", zap.String("user name ", userName))
		c.Error(errcode.ErrEmail)
		return
	}

	if user.UserLink == "" {
		user.UserLink = model.UserLinkPre + user.UserName
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"user_id":      user.ID,
			"user_name":    user.UserName,
			"user_avatar":  user.UserAvatar,
			"is_confirmed": user.IsConfirmed,
			"user_link":    user.UserLink,
			"user_email":   user.Email,
		},
	})

	return

}
