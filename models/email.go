package models

import (
	"fmt"

	"github.com/go-mail/mail/v2"
)

const (
	DefaultSender = "support@lenslocked.com"
)

type Email struct {
	From      string
	To        string
	Subject   string
	Plaintext string
	HTML      string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type EmailService struct {
	DefaultSender string // Use when no other sender is specified

	// unexported fields
	dialer *mail.Dialer
}

func NewEmailService(config SMTPConfig) *EmailService {
	es := EmailService{
		dialer: mail.NewDialer(config.Host, config.Port, config.Username, config.Password),
	}

	return &es
}

func (e *EmailService) Send(email Email) error {
	msg := mail.NewMessage()
	msg.SetHeader("To", email.To)
	e.setFrom(msg, email)
	msg.SetHeader("Subject", email.Subject)

	switch {
	case email.Plaintext != "" && email.HTML != "":
		msg.SetBody("text/plain", email.Plaintext)
		msg.AddAlternative("text/html", email.HTML)
	case email.Plaintext != "":
		msg.SetBody("text/plain", email.Plaintext)
	case email.HTML != "":
		msg.SetBody("text/html", email.HTML)
	}

	err := e.dialer.DialAndSend(msg)
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}

	return nil
}

func (e *EmailService) setFrom(msg *mail.Message, email Email) {
	var from string

	switch {
	case email.From != "":
		from = email.From
	case e.DefaultSender != "":
		from = e.DefaultSender
	default:
		from = DefaultSender
	}

	msg.SetHeader("From", from)
}

func (e *EmailService) ForgotPassword(to, resetUrl string) error {
	email := Email{
		Subject:   "Reset your password",
		To:        to,
		Plaintext: "To reset your password, please visit the following link: " + resetUrl,
		HTML:      `<p>To reset your password, please visit the following link:</p><a href="` + resetUrl + `">` + resetUrl + `</a>`,
	}

	err := e.Send(email)
	if err != nil {
		return fmt.Errorf("forgot password email: %w", err)
	}

	return nil
}
