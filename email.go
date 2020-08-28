package main

import (
	"errors"
	"fmt"
	"gregoryalbouy-server-go/clog"
	"net/smtp"
	"os"
)

// Email represents an email
type Email struct {
	host     string
	port     string
	identity string
	password string
	from     string
	to       []string
	Subject  string
	Body     string
}

// NewEmail sets an email using params that must be set as environment variables
func NewEmail() *Email {
	return &Email{
		host:     os.Getenv("EMAIL_SMTP_HOST"),
		port:     os.Getenv("EMAIL_SMTP_PORT"),
		identity: os.Getenv("EMAIL_IDENTITY"),
		password: os.Getenv("EMAIL_PASSWORD"),
		from:     os.Getenv("EMAIL_FROM"),
		to:       []string{os.Getenv("EMAIL_TO")},
	}
}

// NewEmailFromMessage takes a *Message in input and sets an *Email
// ready to be sent
func NewEmailFromMessage(msg *Message) *Email {
	m := NewEmail()
	m.Subject = fmt.Sprintf("Nouveau message de la part de %s", msg.Email)
	m.Body = fmt.Sprintf("%s\n\nDate: %s\nIP: %s\n", msg.Message, msg.Date, msg.IP)
	return m
}

// Send sends an email
func (m *Email) Send() error {
	if err := m.valid(); err != nil {
		return err
	}

	addr := m.host + ":" + m.port
	auth := smtp.PlainAuth(m.identity, m.from, m.password, m.host)
	msg := []byte(fmt.Sprintf("Subject: %s\n\n%s", m.Subject, m.Body))

	return smtp.SendMail(addr, auth, m.from, m.to, msg)
}

func (m *Email) valid() error {
	if m.host == "" ||
		m.port == "" ||
		// m.Identity == "" ||
		m.password == "" ||
		m.from == "" ||
		m.to[0] == "" {
		clog.Errorlb(fmt.Sprintf("%+v", m), "Missing mail environment variable")
		return errors.New("Internal error")
	}

	if m.Body == "" {
		clog.Errorlb(m, "Missing user values")
		return errors.New("Email incomplete")
	}

	return nil
}
