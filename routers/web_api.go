package routers

import (
	"onbio/handler"

	"github.com/gin-gonic/gin"
)

// InitApiRoute collect task routers
func InitApiRoute(router *gin.Engine) *gin.Engine {
	appRouter := router.Group("/api")
	appRouter.GET("index/login", handler.HandleHello)
	appRouter.GET("index/logout", handler.HandleHello)

	// 后台服务路由
	bkRouter := appRouter.Group("/").Use()
	{
		// test
		bkRouter.POST("/hello", handler.HandleHello)

	}
	return router
}
