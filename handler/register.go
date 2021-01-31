package handler

import (
	"net/http"
	"onbio/logger"
	"onbio/model"
	"onbio/redis"
	"onbio/utils/errcode"

	"github.com/gin-gonic/gin"
	redigo "github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
)

const (
	USER_SENSITIVE_WORD_KEY = "onbio_user_sensitive_word"
)

type RegisterParam struct {
	UserName string `json:"user_name" binding:"required"`
	UserPwd  string `json:"user_pwd" binding:"required"`
	//UserAvatar string `json:"user_avatar" binding:"required"`
	Email string `json:"email" binding:"required"`
}

func checkIfUserNameValid(userName string) (isValid bool) {

	len := len(userName)

	if len > 30 || len < 6 {
		logger.Info("use name len is invalid", zap.String("user", userName))
		return false
	}

	for _, r := range userName {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r != '.') && (r != '_') {
			logger.Info("use name content is invalid", zap.String("user", userName))
			return false
		}
	}

	return true
}

func checkIfUserNameSensitve(userName string) (isSensi bool) {
	//检查是否敏感词
	// 根据cookie获取用户信息
	conn := redis.GetConn("onbio")
	defer conn.Close()

	key := USER_SENSITIVE_WORD_KEY

	isSensitive, err := redigo.Int(conn.Do("sismember", key, userName))

	if err != nil {
		logger.Error("check if sensitive err .", zap.Error(err))
		return true
	}
	logger.Info("is sensitive ", zap.Int("is sensitive", isSensitive))
	if isSensitive == 1 {
		logger.Info("check if sensitive return true", zap.String("user", userName))
		return true
	}
	return false
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
	userAvatar := ""
	//1.check params
	if userName == "" || userPwd == "" || email == "" {
		c.Error(errcode.ErrParam)
		return
	}

	isUserNameValid := checkIfUserNameValid(userName)

	if !isUserNameValid {
		logger.Info("userName is invalid ", zap.String("user", userName))
		c.Error(errcode.ErrUserNameInvalid)
		return
	}
	logger.Info("check if sensitive ")
	isSensitive := checkIfUserNameSensitve(userName)

	if isSensitive {
		logger.Info("userName is sensitive ", zap.String("user", userName))
		c.Error(errcode.ErrUserIsSensitive)
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
		c.Error(errcode.ErrUserExisted)
		return
	}

	err, isExisted = model.IsEmailExisted(email)

	if err != nil {
		logger.Error("check if email existed failed ", zap.Error(err))
		c.Error(errcode.ErrDbQuery)
		return

	}
	if isExisted {
		logger.Error("email existed ", zap.String("email", email))
		c.Error(errcode.ErrEmailExisted)
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
