package main

import "time"

const (
	emailMinLength   = 5
	emailMaxLength   = 255
	subjectMinLength = 0
	subjectMaxLength = 255
	messageMinLength = 3
	messageMaxLength = 5000
)

// Message represents a message
type Message struct {
	ID      int64  `db:"id" json:"id"`
	Email   string `db:"email" json:"email"`
	Subject string `db:"subject" json:"subject"`
	Message string `db:"message" json:"message"`
	Date    string `db:"date" json:"date"`
	IP      string `db:"ip" json:"ip"`
}

// NewMessage returns a *Message with the current date
func NewMessage() *Message {
	return &Message{Date: time.Now().Format(time.RFC3339)}
}

// Valid checks if a message is valid before insertion
func (m Message) Valid() bool {
	return m.validEmail() && m.validSubject() && m.validMessage()
}

// TODO: Regex to check email
func (m Message) validEmail() bool {
	return len(m.Email) >= emailMinLength &&
		len(m.Email) <= emailMaxLength
}

func (m Message) validSubject() bool {
	return len(m.Subject) >= subjectMinLength &&
		len(m.Subject) <= subjectMaxLength
}

func (m Message) validMessage() bool {
	return len(m.Message) >= messageMinLength &&
		len(m.Message) <= messageMaxLength
}
