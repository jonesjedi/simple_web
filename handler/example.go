package handler

import (
	"net/http"
	logger "onbio/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//测试请求
func HandleTestRequest(c *gin.Context) {

	//打印请求
	logger.Info("receive request", zap.Any("req", c.Request.Body))
	c.String(http.StatusOK, "test")
}
