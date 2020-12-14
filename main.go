package main

import (
	"fmt"
	"net/http"
	"onbio/conf"
	"onbio/log"
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

	//初始化日志
	log.Init()
	log.Info("test")

	router := gin.New()

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
			"msg":  "Api route not found",
		})
	})

	routers.InitApiRoute(router)

	gin.SetMode(gin.ReleaseMode) //always release

	log.Error("test")
	log.Info("svr listen on 8000")
	router.Run(":8000")

}
