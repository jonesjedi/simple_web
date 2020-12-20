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

//form 表单提交
func HandleLoginRequest(c *gin.Context) {

	var (
		jsons = jsoniter.ConfigCompatibleWithStandardLibrary
	)

	userName := c.PostForm("user_name")
	userPwd := c.PostForm("user_pwd")

	if userName == "" || userPwd == "" {
		c.Error(errcode.ErrParam)
		return
	}

	//原本这里要校验登录次数 。先打个桩，暂时不实现

	//检查登录密码
	err, user := model.CheckUserPwd(userName, userPwd)

	if err != nil {
		logger.Error("check user pwd failed ", zap.Error(err))
		c.Error(errcode.ErrUserPwd)
		return
	}

	//校验正常的情况下，需要设置cookie，并跳转个人首页
	//先生成一个cookie
	sessionKey := fmt.Sprintf("%s", uuid.NewV4())

	var sessionContent services.SessionContent
	sessionContent.Email = user.Email
	sessionContent.IsConfirmed = user.IsConfirmed
	sessionContent.UserAvatar = user.UserAvatar
	sessionContent.UserID = user.ID
	sessionContent.UserLink = user.UserLink
	sessionContent.LoginTime = uint64(time.Now().Unix())

	//序列化
	sessionStr, err := jsons.Marshal(sessionContent)
	if err != nil {
		logger.Error("err json Marshal", zap.Any("session", sessionContent))
		return
	}
	conn := redis.GetConn("onbio")
	defer conn.Close()

	key := fmt.Sprintf(services.USER_SESSION_REDIS_PRE, sessionKey)
	_, err = conn.Do("SET", key, string(sessionStr))
	if err != nil && err != redigo.ErrNil {
		logger.Error("err set redis ", zap.String("key", key), zap.Error(err))
		c.Error(errcode.ErrRedisOper)
		return
	}

	//设置cookie
	c.SetCookie("onbio_user", sessionKey, 86400, "/", ".onb.io", false, true)

	//跳转
	c.Redirect(http.StatusFound, services.USER_REDIRECT_URL)

}
