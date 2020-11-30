package handler

import (
	"net/http"
	"onbio/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HandleYunApiRequest(c *gin.Context) {

	//打印请求
	log.Info("receive request", zap.Any("req", c.Request.Body))

	c.String(http.StatusOK, "test")
}
