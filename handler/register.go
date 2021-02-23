package handler

import (
	"fmt"
	"net/http"
	"onbio/logger"
	"onbio/model"
	"onbio/redis"
	"onbio/services"
	"onbio/utils/errcode"
	"time"

	"github.com/gin-gonic/gin"
	redigo "github.com/gomodule/redigo/redis"
	jsoniter "github.com/json-iterator/go"
	uuid "github.com/satori/go.uuid"
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
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && (r != '.') && (r != '_') {
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

	//注册成功之后，自动登录一次
	err, user := model.GetUserInfo("", userName, 0)

	if err != nil {
		logger.Info("get user info by username failed ", zap.String("user name ", userName))
		c.Error(errcode.ErrEmail)
		return
	}

	//先生成一个cookie
	sessionKey := fmt.Sprintf("%s", uuid.NewV4())

	var sessionContent services.SessionContent
	sessionContent.Email = user.Email
	sessionContent.IsConfirmed = user.IsConfirmed
	sessionContent.UserAvatar = user.UserAvatar
	sessionContent.UserID = user.ID
	sessionContent.UserName = user.UserName
	sessionContent.UserLink = user.UserLink
	sessionContent.LoginTime = uint64(time.Now().Unix())

	var (
		jsons = jsoniter.ConfigCompatibleWithStandardLibrary
	)
	//序列化
	sessionStr, err := jsons.Marshal(sessionContent)
	if err != nil {
		logger.Error("err json Marshal", zap.Any("session", sessionContent))
		return
	}
	conn := redis.GetConn("onbio")
	defer conn.Close()

	key := fmt.Sprintf(services.USER_SESSION_REDIS_PRE, sessionKey)
	_, err = conn.Do("SET", key, string(sessionStr), "EX", 86400)
	if err != nil && err != redigo.ErrNil {
		logger.Error("err set redis ", zap.String("key", key), zap.Error(err))
		c.Error(errcode.ErrRedisOper)
		return
	}

	//设置cookie
	c.SetCookie("onbio_user", sessionKey, 86400, "/", ".onb.io", false, true)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{},
	})
}
