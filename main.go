package main

import (
	"fmt"
	"onbio/conf"
	"onbio/handler"
	"onbio/log"
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

	//初始化日志
	log.Init()
	log.Info("test")

	router := gin.Default()

	gin.SetMode(gin.ReleaseMode)

	router.POST("/request", handler.HandleYunApiRequest)

	router.GET("/request", handler.HandleYunApiRequest)

	log.Info("svr listen on 8000")
	router.Run(":8000")

}
