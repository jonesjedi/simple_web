package handler

import (
	"fmt"
	"net/http"
	logger "onbio/logger"
	"onbio/redis"
	"onbio/services"
	"onbio/utils/errcode"

	"github.com/gin-gonic/gin"
	redigo "github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
)

func HandleLogoutRequest(c *gin.Context) {

	//先获取对应的cookie key，删掉在redis里的key，下个请求就失效了

	sessionKey := c.GetString("cookie_key")

	if sessionKey == "" {
		logger.Error("logout :empty sessionKey ")
		c.Error(errcode.ErrInternal)
		return
	}
	conn := redis.GetConn("onbio")
	defer conn.Close()

	key := fmt.Sprintf(services.USER_SESSION_REDIS_PRE, sessionKey)

	_, err := conn.Do("DEL", key)
	if err != nil && err != redigo.ErrNil {
		logger.Error("err get redis ", zap.String("key", key), zap.Error(err))
		c.AbortWithError(http.StatusOK, c.Error(errcode.ErrRedisOper))
		return
	}

	//设置cookie过期
	c.SetCookie("onbio_user", sessionKey, 0, "/", ".onb.io", false, true)

	c.Redirect(http.StatusFound, services.USER_REDIRECT_URL)
	return

}
