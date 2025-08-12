package utils

import (
	"fmt"
	"log"
	"mail-service/internal/config"

	"gopkg.in/gomail.v2"
)

type MailSender struct {
	cfg *config.MailConfig
}

func NewMailSender(cfg *config.MailConfig) *MailSender {
	return &MailSender{cfg: cfg}
}

func (m *MailSender) SendEmail(to, subject, htmlBody string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.cfg.FromEmail)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(m.cfg.Host, m.cfg.Port, m.cfg.FromEmail, m.cfg.Password)

	if err := d.DialAndSend(msg); err != nil {
		log.Println("Send email failed:", err)
		return err
	}
	return nil
}

func (m *MailSender) SendVerificationEmail(userEmail, token string) error {
	link := fmt.Sprintf("%s/api/verify-account?token=%s", m.cfg.AuthServiceURL, token)
	subject := "Verify your account"
	html := fmt.Sprintf(`
		<h2>Welcome!</h2>
		<p>Please verify your email by clicking the link below:</p>
		<a href="%s">Verify Email</a>
	`, link)
	return m.SendEmail(userEmail, subject, html)
}
