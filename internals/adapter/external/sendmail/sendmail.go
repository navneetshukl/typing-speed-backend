package sendmail

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type MailSender interface {
	SendMail(from, to, subject, body string) error
}

type GoMail struct {
	Host string
	Port int
	User string
	Pass string
}

func NewGoMail(host string, port int) MailSender {
	return &GoMail{
		Host: host,
		Port: port,
	}
}

func (a *GoMail) SendMail(from, to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(a.Host, a.Port, a.User, a.Pass)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send mail: %w", err)
	}

	return nil
}
