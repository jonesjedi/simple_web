package handler

import (
	"fmt"
	"net/http"
	logger "onbio/logger"
	"onbio/model"
	"onbio/redis"
	"onbio/utils/email_html"
	"onbio/utils/errcode"
	"onbio/utils/mailsender"
	"strings"

	"github.com/gin-gonic/gin"
	redigo "github.com/gomodule/redigo/redis"
	jsoniter "github.com/json-iterator/go"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

const (
	USER_EMAIL_RESET_PWD_URL   = "http://onb.io/user/setPassWord?code=%s"
	USER_RESET_PWD_CONTENT_PRE = "user_reset_pwd:%s"
)

type ResetPwdEmailParam struct {
	UserEmailOrUserName string `json:"user_email" binding:"required"`
}

func HandleSendResetPwdEmailRequest(c *gin.Context) {

	var params ResetPwdEmailParam
	err := c.Bind(&params)
	if err != nil {
		logger.Error("params err ")
		c.Error(errcode.ErrParam)
		return
	}
	var user model.User
	//如果包含@，就是邮箱，因为用户名现在不能包含@
	logger.Info("params", zap.Any("params", params), zap.String("user", params.UserEmailOrUserName))
	if strings.Contains(params.UserEmailOrUserName, "@") {
		err, user = model.GetUserInfo(params.UserEmailOrUserName, "", 0)
	} else {
		err, user = model.GetUserInfo("", params.UserEmailOrUserName, 0)
	}
	//判断是哪个用户的email
	if err != nil {
		logger.Info("get user info by email failed ", zap.String("user email", params.UserEmailOrUserName))
		c.Error(errcode.ErrEmail)
		return
	}

	//到这里，就可以发邮件了
	err, code := GenResetPwdCode(user.ID, params.UserEmailOrUserName)
	if err != nil {
		logger.Error("gen valid code failed ,", zap.Error(err))
		c.Error(errcode.ErrInternal)
		return
	}

	sendUrl := fmt.Sprintf(USER_EMAIL_RESET_PWD_URL, code)

	//没有接口，先打个log
	logger.Info("reset  url ", zap.String("url", sendUrl))

	emailBody, err := email_html.GenerateHtml(user.UserName, sendUrl, 2)

	if err != nil {
		logger.Error("generate email body failed ", zap.String("user", user.UserName))
		c.Error(errcode.ErrInternal)
		return
	}
	var ms mailsender.MailSender = &mailsender.Mail{
		Sender:    "Onbio<welcome@onb.io>", // 可以自定义
		Recipient: user.Email,              // 如果处于Sandbox只能发送已验证过的邮箱
		Subject:   "onbio reset pwd email",
		HTMLBody:  emailBody,
		TextBody:  emailBody, // 不支持HTML的话会返回这个
		CharSet:   "UTF-8",   // 固定字符码
	}
	sendRet := ms.SendMail()
	if !sendRet {
		logger.Error("send mail failed ")
		c.Error(errcode.ErrInternal)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"email": user.Email,
		},
	})

}

type ResetPwdContent struct {
	UserID    uint64 `json:"user_id"`
	UserEmail string `json:"user_email"`
}

func GenResetPwdCode(userId uint64, userEmail string) (err error, emailValidCode string) {

	var (
		jsons = jsoniter.ConfigCompatibleWithStandardLibrary
	)

	//先生成一个uuid，随机code
	emailValidCode = fmt.Sprintf("%s", uuid.NewV4())

	var resetPwdContent ResetPwdContent
	resetPwdContent.UserEmail = userEmail
	resetPwdContent.UserID = userId

	//序列化
	resetPwdContentStr, err := jsons.Marshal(resetPwdContent)
	if err != nil {
		logger.Error("err json Marshal", zap.Any("reset pwd content", resetPwdContent))
		return
	}
	conn := redis.GetConn("onbio")
	defer conn.Close()

	key := fmt.Sprintf(USER_RESET_PWD_CONTENT_PRE, emailValidCode)
	_, err = conn.Do("SET", key, string(resetPwdContentStr), "EX", 6*3600, "NX")
	if err != nil && err != redigo.ErrNil {
		logger.Error("err set redis ", zap.String("key", key), zap.Error(err))
		return
	}
	return
}
