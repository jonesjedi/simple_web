package routers

import (
	"onbio/handler"
	"onbio/middlewares"

	"github.com/gin-gonic/gin"
)

// InitApiRoute collect task routers
func InitApiRoute(router *gin.Engine) *gin.Engine {
	appRouter := router.Group("/api")
	appRouter.POST("/user/login", handler.HandleLoginRequest)
	appRouter.GET("/index/logout", handler.HandleTestRequest)
	appRouter.POST("/user/register", handler.HandleRegisteRequest)

	appRouter.POST("/user/send_validate_email", handler.HandleSendValidateEmailRequest)
	appRouter.GET("/user/validate_email", handler.HandleValidateEmailRequest)

	appRouter.POST("/user/send_reset_pwd", handler.HandleSendResetPwdEmailRequest)
	appRouter.POST("/user/reset_pwd", handler.HandleResetPwdRequest)

	// 后台服务路由
	bkRouter := appRouter.Group("/").Use()
	{
		// test
		bkRouter.POST("/hello", handler.HandleTestRequest)

	}
	//需要验证登录态的接口
	linkRouter := appRouter.Group("link/", middlewares.OnbioLoginAuth())
	{
		linkRouter.POST("/updateinfo", handler.HandleUpdateUserInfoRequest)
		linkRouter.GET("/userlink", handler.HandleGetUserLinkRequest)
		linkRouter.GET("/userinfo", handler.HandleGetUserInfoRequest)
		linkRouter.POST("/updatelink", handler.HandleUpdateUserLinkRequest)
		linkRouter.POST("/createlink", handler.HandleCreateUserLinkRequest)
	}

	return router
}
