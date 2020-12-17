package middlewares

import (
	"fmt"
	"net/http"

	"onbio/utils/errcode"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// ResponseHandler across domain
func ResponseHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var ee error

				ee = errors.New(fmt.Sprintf("[Recovery] %s panic recovered", err))
				c.JSON(http.StatusOK, gin.H{
					"code": http.StatusInternalServerError,
					"msg":  ee.Error(),
					"data": gin.H{},
				})
				return
			}
		}()
		c.Next()
		err := c.Errors.ByType(gin.ErrorTypeAny).Last()
		if err != nil {
			if err.Meta != nil {
				c.JSON(http.StatusOK, err.Meta)
			} else {
				if e, ok := err.Err.(errcode.StandardError); ok {
					c.JSON(http.StatusOK, gin.H{
						"code": e.Code,
						"msg":  e.Msg,
						"data": gin.H{},
					})
				} else {
					c.JSON(http.StatusOK, gin.H{
						"code": http.StatusInternalServerError,
						"msg":  err.Error(),
						"data": gin.H{},
					})
				}
			}
		}
	}
}
