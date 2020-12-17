package main

import (
	"fmt"
	"net/http"
	"onbio/conf"
	"onbio/logger"
	"onbio/middlewares"
	"onbio/mysql"
	"onbio/routers"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	//读取配置

	if len(os.Args) < 2 {
		fmt.Println("usage : ./program  configpath")
		return
	}
	filePath := os.Args[1]
	error := conf.LoadConf(filePath)
	if error != nil {
		fmt.Println("load conf failed")
		return
	}
	logger.Init()
	//初始化日志
	logger.Info("test")

	//初始化db
	mysql.Init("onbio")

	router := gin.New()
	router.Use(middlewares.ResponseHandler())
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
			"msg":  "Api route not found",
		})
	})

	routers.InitApiRoute(router)

	gin.SetMode(gin.ReleaseMode) //always release

	logger.Info("svr listen on 8000")
	router.Run(":8000")

}
