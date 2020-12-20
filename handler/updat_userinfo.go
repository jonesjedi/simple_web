package handler

import "github.com/gin-gonic/gin"

//更新用户信息，暂时只能更新用户头像，这个接口需要验证登录态
//验证登录态，统一放middlewares
func HandleUpdateUserInfoRequest(c *gin.Context) {
	return
}
