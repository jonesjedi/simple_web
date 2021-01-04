package mailsender

import (
	"fmt"

	// go get -u github.com/aws/aws-sdk-go
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	// AKID 全称 AWS_ACCESS_KEY_ID
	// 创建aws sdk session使用，其他方式对于部署和开发较复杂
	// 现使用Hard-Coded模式硬编码
	// 开发使用dev用户生成的访问权限，线上可更换为生产环境权限
	// 初始为空请获得到访问权限填充
	AKID = ""
	// SECREY_KEY 全称 AWS_SECRET_ACCESS_KEY
	SECREY_KEY = ""
)

// MailSender 邮件发送者接口
type MailSender interface {
	SendMail() bool
}

// Mail 发送电子邮件的结构体
type Mail struct {
	// 发送者
	Sender string
	// 收件人
	Recipient string
	// 主题
	Subject string
	// HTML邮件主体
	HTMLBody string
	// 文字邮件主题（不支持HTML的邮件客户端使用）
	TextBody string
	// 字符集
	CharSet string
}

// SendMail 发送邮件的实现
func (m *Mail) SendMail() bool {
	// 创建新Session访问aws服务，需使用相对的Region
	// 当前只有新加坡区可以使用
	// 当前TOKEN参数可为空
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(AKID, SECREY_KEY, "")},
	)
	svc := ses.New(sess)

	// 组装邮件
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(m.Recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(m.CharSet),
					Data:    aws.String(m.HTMLBody),
				},
				Text: &ses.Content{
					Charset: aws.String(m.CharSet),
					Data:    aws.String(m.TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(m.CharSet),
				Data:    aws.String(m.Subject),
			},
		},
		Source: aws.String(m.Sender),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// 发送邮件
	result, err := svc.SendEmail(input)

	// 如果发生错误显示错误信息
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}

		return false
	}
	// 如果得到MessageId证明操作成功
	if res := result.MessageId; *res != "" {
		return true
	}
	return false
}
