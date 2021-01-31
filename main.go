package main

import (
	"fmt"
	"net/http"
	"onbio/conf"
	"onbio/logger"
	"onbio/middlewares"
	"onbio/mysql"
	"onbio/redis"
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
	redis.Init("onbio")

	router := gin.New()
	router.Use(middlewares.ResponseHandler())
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
			"msg":  "Api route not found",
		})
	})

	//暂时放在这里导入敏感词，后面应该做到管理端
	/*
		err := import_sensitive.ImportSensitiveWordFromExcel("./admin/import_sensitive/OnBio_sensitive.xlsx")

		if err != nil {
			logger.Error("import sensitive word failed", zap.Error(err))
		}*/

	routers.InitApiRoute(router)

	gin.SetMode(gin.ReleaseMode) //always release

	logger.Info("svr listen on 8000")
	router.Run(":8000")

}
