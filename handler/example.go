package handler

import (
	"github.com/gin-gonic/gin"
)

func HandleHello(c *gin.Context) {
	c.JSON(200, gin.H{"status": "Hello World"})
}
