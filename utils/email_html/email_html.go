package email_html

import (
	"onbio/logger"

	"github.com/matcornic/hermes/v2"
	"go.uber.org/zap"
)

type EmailType uint32

const (
	Register_Account EmailType = 1
	Reset_Pwd        EmailType = 2
)

func (p EmailType) GetContent() string {
	switch p {
	case Register_Account:
		return "Thanks for signing up to Linktree, great to have you!"
	case Reset_Pwd:
		return "You are trying to reset your password"
	default:
		return "UNKNOWN"
	}
}
func (p EmailType) GetIntroContent() string {
	switch p {
	case Register_Account:
		return "Please verify your email address by clicking the link below."
	case Reset_Pwd:
		return "To reset your password, please click here"
	default:
		return "UNKNOWN"
	}
}

func GenerateHtml(userName, url string, emailType EmailType) (emailBody string, err error) {

	h := hermes.Hermes{
		// Optional Theme
		// Theme: new(Default)
		Product: hermes.Product{
			// Appears in header & footer of e-mails
			Name: "Onb.io",
			Link: "https://onb.io/",
			// Optional product logo
			Logo: "http://onb.io/_nuxt/assets/images/logo.png",
		},
	}
	intros := []string{emailType.GetContent()}
	email := hermes.Email{
		Body: hermes.Body{
			Name:   userName,
			Intros: intros,
			Actions: []hermes.Action{
				{
					Instructions: emailType.GetIntroContent(),
					Button: hermes.Button{
						Color: "#22BC66", // Optional action button color
						Text:  "Go",
						Link:  url,
					},
				},
			},
			Outros: []string{
				"",
			},
		},
	}

	// Generate an HTML email with the provided contents (for modern clients)
	emailBody, err = h.GenerateHTML(email)
	if err != nil {
		logger.Error("GenerateHTML failed ", zap.Error(err))
		return
	}
	return
}
