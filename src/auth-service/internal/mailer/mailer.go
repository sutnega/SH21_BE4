package mailer

import "log"

type Mailer interface {
	SendPassword(email, password string)
}

type consoleMailer struct{}

func New() Mailer {
	return &consoleMailer{}
}

func (m *consoleMailer) SendPassword(email, password string) {
	log.Printf("[MAILER STUB] send to %s: your new password = %s", email, password)
}


