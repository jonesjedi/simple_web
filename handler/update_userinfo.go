package handler

import (
	"net/http"
	logger "onbio/logger"
	"onbio/model"
	"onbio/utils/errcode"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//更新用户信息，暂时只能更新用户头像，这个接口需要验证登录态
//验证登录态，统一放middlewares

type UpdateUserInfoParam struct {
	UserAvatar string `json:"user_avatar" binding:"required"`
}

func HandleUpdateUserInfoRequest(c *gin.Context) {

	var params UpdateUserInfoParam
	err := c.Bind(&params)
	if err != nil {
		logger.Error("params err ,", zap.Error(err))
		c.Error(errcode.ErrParam)
		return
	}

	//新头像链接是否需要经过检验是否涉黄之类
	userAvatar := params.UserAvatar

	userID := (uint64)(c.GetInt("user_id"))

	//更新对应用户的头像
	var user model.User
	user.UserAvatar = userAvatar
	err = model.UpdateUserInfoByID(userID, user)

	if err != nil {
		logger.Error("UpdateUserInfoByID failed ", zap.Error(err))
		c.Error(errcode.ErrInternal)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{},
	})

	return

}
