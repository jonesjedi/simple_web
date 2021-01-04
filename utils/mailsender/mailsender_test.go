package mailsender

import (
	"fmt"
	"testing"
)

// TestSendMail 测试发送邮件功能
func TestSendMail(t *testing.T) {
	var ms MailSender = &Mail{
		Sender:    "sender@onb.io",       // 可以自定义
		Recipient: "jonesjedi@gmail.com", // 如果处于Sandbox只能发送已验证过的邮箱
		Subject:   "Amazon SES Test (AWS SDK for Go)",
		HTMLBody: "<h1>Amazon SES Test Email (AWS SDK for Go)</h1><p>This email was sent with " +
			"<a href='https://aws.amazon.com/ses/'>Amazon SES</a> using the " +
			"<a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.</p>",
		TextBody: "This email was sent with Amazon SES using the AWS SDK for Go.", // 不支持HTML的话会返回这个
		CharSet:  "UTF-8",                                                         // 固定字符码
	}
	fmt.Println(ms.SendMail())
}
