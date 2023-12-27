package smtp

import (
	"crypto/tls"

	gomail "gopkg.in/mail.v2"
)

type SenderConfig struct {
	From     string
	Password string
	SMTPHost string `default:"smtp.gmail.com"`
	SMTPPort int    `default:"587"`
}

type Sender interface {
	Send(to []string, subj, contentType, body string) error
}

type sender struct {
	from     string
	password string

	smtpHost string
	smtpPort int
}

func NewSender(cfg *SenderConfig) Sender {
	return &sender{
		from:     cfg.From,
		password: cfg.Password,
		smtpHost: cfg.SMTPHost,
		smtpPort: cfg.SMTPPort,
	}
}

func (s *sender) Send(to []string, subj, contentType, body string) error {
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", s.from)

	// Set E-Mail receivers
	m.SetHeader("To", to...)

	// Set E-Mail subject
	m.SetHeader("Subject", subj)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody(contentType, body)

	// Settings for SMTP server
	d := gomail.NewDialer(s.smtpHost, s.smtpPort, s.from, s.password)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d.DialAndSend(m)
}
