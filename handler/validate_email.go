package handler

import (
	"fmt"
	"net/http"
	logger "onbio/logger"
	"onbio/model"
	"onbio/redis"
	"onbio/utils/errcode"

	"github.com/gin-gonic/gin"
	redigo "github.com/gomodule/redigo/redis"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

//验证邮箱接口，主要根据code来验证，是否是之前发出去的邮件
//Get接口，这里只需要一个code做验证
/**
1.根据code验证是否是来自之前发出去的邮件
2.将用户数据更新为已验证

**/

const (
	Validate_Succ_URL = "http://onb.io/user/mailPassed"
	Validate_Fail_URL = "http://onb.io/"
)

func HandleValidateEmailRequest(c *gin.Context) {
	code := c.DefaultQuery("code", "")

	if code == "" {
		logger.Info("err code ")
		c.Error(errcode.ErrParam)
		return
	}

	//验证这个code是否有效
	err, valid := ValidateCode(code)

	if err != nil {
		logger.Error("Validate code failed ", zap.Error(err))
		c.Error(errcode.ErrCodeInValid)
		return
	}

	userId := valid.UserID

	var user model.User
	user.IsConfirmed = 1
	err = model.UpdateUserInfoByID(userId, user)

	if err != nil {
		logger.Error("UpdateUserInfoByID failed ", zap.Error(err))
		//c.Error(errcode.ErrInternal)
		c.Redirect(http.StatusFound, Validate_Fail_URL)
		return
	}
	c.Redirect(http.StatusFound, Validate_Succ_URL)
	/***直接跳转了，不要返回数据
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{},
	})***/
}

func ValidateCode(code string) (err error, valid ValidContent) {

	var (
		jsons = jsoniter.ConfigCompatibleWithStandardLibrary
	)

	conn := redis.GetConn("onbio")
	defer conn.Close()

	key := fmt.Sprintf(USER_VALID_CONTENT_PRE, code)
	validContent, err := redigo.Bytes(conn.Do("GET", key))
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
	err = jsons.Unmarshal(validContent, &valid)
	if err != nil {
		logger.Error("err Unmarshal valid info ", zap.ByteString("key", validContent), zap.Error(err))
		return
	}
	return
}
