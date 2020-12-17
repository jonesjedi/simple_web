package routers

import (
	"onbio/handler"

	"github.com/gin-gonic/gin"
)

// InitApiRoute collect task routers
func InitApiRoute(router *gin.Engine) *gin.Engine {
	appRouter := router.Group("/api")
	appRouter.GET("index/login", handler.HandleTestRequest)
	appRouter.GET("index/logout", handler.HandleTestRequest)
	appRouter.POST("register", handler.HandleRegisteRequest)
	// 后台服务路由
	bkRouter := appRouter.Group("/").Use()
	{
		// test
		bkRouter.POST("/hello", handler.HandleTestRequest)

	}
	return router
}
