package handler

import (
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
	"time"

	t "prplchat/src/model/template"
)

func IsEmailValid(address string) bool {
	_, err := mail.ParseAddress(address)
	return err == nil
}

func sendSignupCompletionEmail(
	tmpl t.VerifyEmailTemplate,
	source string,
	pass string,
) error {
	subject := "User sing up confirmation"
	body, err := tmpl.Email()
	if err != nil {
		return err
	}
	err = sendEmail(source, pass, Email{
		sender:    tmpl.SourceEmail,
		receivers: []string{tmpl.UserEmail},
		subject:   subject,
		body:      body,
	})
	return err
}

type Email struct {
	sender    string
	receivers []string
	subject   string
	body      string
}

func sendEmail(sender string, password string, email Email) error {
	auth := smtp.PlainAuth("", sender, password, "smtp.gmail.com")

	headers := make(map[string]string)
	headers["From"] = sender
	headers["To"] = strings.Join(email.receivers, ",")
	headers["Subject"] = email.subject
	headers["Date"] = time.Now().Format(time.RFC1123Z)
	//headers["CC"] = "cc@example.com"
	//headers["BCC"] = "bcc@example.com"
	headers["Reply-To"] = "do-not-reply@please.com"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	msg := ""
	for k, v := range headers {
		msg += k + ": " + v + "\r\n"
	}
	msg += "\r\n" + email.body

	err := smtp.SendMail("smtp.gmail.com:587", auth, sender, email.receivers, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email about [%s] to [%v], %s", email.subject, email.receivers, err.Error())
	}

	return nil
}
