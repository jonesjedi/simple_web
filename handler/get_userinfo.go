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
	userID := c.GetInt("user_id")

	//根据用户ID获取用户信息

	err, user := model.GetUserInfo("", "", uint64(userID))

	if err != nil {
		logger.Info("get user info by email failed ", zap.Int("user ID ", userID))
		c.Error(errcode.ErrEmail)
		return
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
		},
	})

	return

}
