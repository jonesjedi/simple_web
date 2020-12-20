package handler

import (
	"fmt"
	"net/http"
	logger "onbio/logger"
	"onbio/model"
	"onbio/redis"
	"onbio/utils"
	"onbio/utils/errcode"

	"github.com/gin-gonic/gin"
	redigo "github.com/gomodule/redigo/redis"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type ResetPwdParam struct {
	Code   string `json:"code" binding:"required"`
	NewPwd string `json:"new_pwd" binding:"required"`
}

///重置密码请求
func HandleResetPwdRequest(c *gin.Context) {

	var params ResetPwdParam
	err := c.Bind(&params)
	if err != nil {
		logger.Error("params err ")
		c.Error(errcode.ErrParam)
		return
	}

	//验证code，看来源是否正确
	err, resetInfo := ValidateResetPwdCode(params.Code)

	if err != nil {
		logger.Error("ValidateResetPwdCode failed ", zap.Error(err))
		c.Error(errcode.ErrCodeInValid)
		return
	}

	//获取对应的用户ID,更新密码
	userID := resetInfo.UserID

	userPwd, _ := utils.Md5Sum(params.NewPwd)

	var user model.User
	user.UserPwd = userPwd
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

}
func ValidateResetPwdCode(code string) (err error, resetPwdContent ResetPwdContent) {

	var (
		jsons = jsoniter.ConfigCompatibleWithStandardLibrary
	)

	conn := redis.GetConn("onbio")
	defer conn.Close()

	key := fmt.Sprintf(USER_RESET_PWD_CONTENT_PRE, code)
	resetPwdContentStr, err := redigo.Bytes(conn.Do("GET", key))
	if err != nil && err != redigo.ErrNil {
		logger.Error("err get redis ", zap.String("key", key), zap.Error(err))
		return
	}
	if err == redigo.ErrNil {
		logger.Error("invalid code", zap.String("code", code))
		return
	}
	//删掉这个code,只能验证一次
	_, err = redigo.Int(conn.Do("DEL", key))
	if err != nil && err != redigo.ErrNil {
		logger.Error("err del redis", zap.String("key", key), zap.Error(err))
		return
	}

	//反序列化
	err = jsons.Unmarshal(resetPwdContentStr, &resetPwdContent)
	if err != nil {
		logger.Error("err Unmarshal valid info ", zap.ByteString("key", resetPwdContentStr), zap.Error(err))
		return
	}
	return
}
