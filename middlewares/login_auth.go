package middlewares

import (
	"fmt"
	"net/http"
	"onbio/logger"
	"onbio/redis"
	"onbio/services"
	"onbio/utils/errcode"

	"github.com/gin-gonic/gin"
	redigo "github.com/gomodule/redigo/redis"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

// OALoginAuth middleware is check ipcc voice agent already logged in
func OnbioLoginAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			jsons = jsoniter.ConfigCompatibleWithStandardLibrary
		)
		userSession, _ := c.Cookie("onbio_user")

		// 未登录
		if userSession == "" {
			c.AbortWithError(http.StatusOK, c.Error(errcode.ErrReqForbidden))
			return
		}

		// 根据cookie获取用户信息
		conn := redis.GetConn("onbio")
		defer conn.Close()

		key := fmt.Sprintf(services.USER_SESSION_REDIS_PRE, userSession)

		sessionContentStr, err := redigo.Bytes(conn.Do("GET", key))
		if err != nil && err != redigo.ErrNil {
			logger.Error("err get redis ", zap.String("key", key), zap.Error(err))
			c.AbortWithError(http.StatusOK, c.Error(errcode.ErrReqForbidden))
			return
		}
		if err == redigo.ErrNil {
			logger.Error("invalid cookie", zap.String("userSession", userSession))
			c.AbortWithError(http.StatusOK, c.Error(errcode.ErrReqForbidden))
			return
		}
		var session services.SessionContent
		err = jsons.Unmarshal(sessionContentStr, &session)
		if err != nil {
			logger.Error("err Unmarshal valid info ", zap.ByteString("key", sessionContentStr), zap.Error(err))
			return
		}
		userId := (int)(session.UserID)
		c.Set("user_name", session.UserName)
		c.Set("user_id", userId)
		c.Set("is_confirmed", session.IsConfirmed)
		c.Set("cookie_key", userSession)
		c.Next()
	}
}
