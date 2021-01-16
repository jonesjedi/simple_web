package handler

import (
	"github.com/gin-gonic/gin"

	logger "onbio/logger"
	"onbio/model"
	"onbio/utils/errcode"
	"go.uber.org/zap"
	"net/http"
)

type DeleteUserLinkParam struct {
	LinkID  int `json:"id" binding:"required"`
}


func HandleDeleteUserLinkRequest(c *gin.Context) {

	var params DeleteUserLinkParam
	err := c.Bind(&params)
	if err != nil {
		logger.Error("params err ")
		c.Error(errcode.ErrParam)
		return
	}
	if params.LinkID == 0{
		logger.Error("params err ")
		c.Error(errcode.ErrParam)
		return
	}
	linkID := (uint64)(params.LinkID)
	userID := (uint64)(c.GetInt("user_id"))

	err = model.DeleteUserLink(userID,linkID)

	if err != nil {
		logger.Error("delete user link failed ",zap.Uint64("userID",userID),zap.Uint64("linkID",linkID))
		c.Error(errcode.ErrInternal)
		return 
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{},
	})
}
