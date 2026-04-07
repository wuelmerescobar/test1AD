package mailer

import (
	"fmt"
	"net/smtp"
)

type Mailer struct {
	Host string
	Port string
	User string
	Pass string
	From string
}

func New(host, port, user, pass, from string) *Mailer {
	return &Mailer{
		Host: host,
		Port: port,
		User: user,
		Pass: pass,
		From: from,
	}
}

func (m *Mailer) Enabled() bool {
	return m.Host != "" && m.Port != "" && m.User != "" && m.Pass != "" && m.From != ""
}

func (m *Mailer) Send(to, subject, body string) error {
	if !m.Enabled() {
		return nil
	}

	addr := fmt.Sprintf("%s:%s", m.Host, m.Port)
	auth := smtp.PlainAuth("", m.User, m.Pass, m.Host)

	msg := []byte(
		"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-version: 1.0;\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\";\r\n\r\n" +
			body + "\r\n",
	)

	return smtp.SendMail(addr, auth, m.From, []string{to}, msg)
}

func (m *Mailer) SendWelcomeStaffEmail(to, name string) error {
	body := fmt.Sprintf("Hello %s,\n\nYour staff account for the Library Management System has been created successfully.\n\nRegards,\nLibrary Admin", name)
	return m.Send(to, "Welcome to the Library Management System", body)
}
