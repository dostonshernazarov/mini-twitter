package etc

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"

	"github.com/dostonshernazarov/mini-twitter/internal/entity"
	"github.com/dostonshernazarov/mini-twitter/internal/pkg/config"
)

func SendMessage(to []string, message entity.SMTPCode, htmlPath string, cfg config.Config) error {
	t, err := template.ParseFiles(htmlPath)
	if err != nil {
		log.Println("Error parsing html file in send mail", err.Error())
		return err
	}

	var k bytes.Buffer
	err = t.Execute(&k, message)
	if err != nil {
		log.Println("failed to executing email body", err.Error())
		return err
	}

	if k.String() == "" {
		log.Println("Error buffer")
	}
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(fmt.Sprintf("Subject: %s", "MiniTwitter - Verification Code\n") + mime + k.String())

	// Authentication.
	auth := smtp.PlainAuth("", cfg.SMTPEmail, cfg.SMTPPassword, cfg.SMTPHost)

	// Sending email.
	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)
	err = smtp.SendMail(addr, auth, cfg.SMTPEmail, to, msg)

	return err
}
