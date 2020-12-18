package handler

import (
	"net/http"
	"onbio/logger"
	"onbio/model"
	"onbio/utils/errcode"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RegisterParam struct {
	UserName   string `json:"user_name" binding:"required"`
	UserPwd    string `json:"user_pwd" binding:"required"`
	UserAvatar string `json:"user_avatar" binding:"required"`
	Email      string `json:"email" binding:"required"`
}

func HandleRegisteRequest(c *gin.Context) {

	var params RegisterParam
	err := c.Bind(&params)
	if err != nil {
		logger.Error("params err ")
		c.Error(errcode.ErrParam)
		return
	}
	userName := params.UserName
	userPwd := params.UserPwd
	email := params.Email
	userAvatar := params.UserAvatar
	//1.check params
	if userName == "" || userPwd == "" || email == "" {
		c.Error(errcode.ErrParam)
		return
	}

	//check if existed
	err, isExisted := model.IsUserExisted(userName)

	if err != nil {
		logger.Error("check if user existed failed ", zap.Error(err))
		c.Error(errcode.ErrDbQuery)
		return

	}
	if isExisted {
		logger.Error("user existed ", zap.String("user", userName))
		c.Error(errcode.ErrDbQuery)
		return
	}
	err = model.CreateUser(userName, userAvatar, userPwd, email)

	if err != nil {
		logger.Error("Create User Failed", zap.Error(err))
		c.Error(errcode.ErrDbQuery)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{},
	})
}