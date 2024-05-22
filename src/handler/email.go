package handler

import (
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"strings"
	"time"

	"go.chat/src/db"
	a "go.chat/src/model/app"
	"go.chat/src/model/template"
	"go.chat/src/utils"
)

func IsEmailValid(address string) bool {
	_, err := mail.ParseAddress(address)
	return err == nil
}

func IssueReservationToken(
	app *AppState,
	db *db.DBConn,
	user *a.User,
) (*template.VerifyEmailTemplate, error) {
	token := utils.RandStringBytes(16)
	expire := time.Now().Add(1 * time.Hour)
	reserve, err := Reserve(db, user, token, expire)
	if err != nil {
		log.Printf("IssueReservationToken ERROR reserving[%s], %s\n", user.Email, err.Error())
		return nil, fmt.Errorf("")
	}
	emailConfig := app.SmtpConfig()
	tmpl := template.VerifyEmailTemplate{
		SourceEmail: emailConfig.User,
		UserEmail:   user.Email,
		UserName:    user.Name,
		Token:       reserve.Token,
		TokenExpire: reserve.Expire.Format(time.RFC3339),
	}
	err = sendSignupCompletionEmail(tmpl)
	if err != nil {
		log.Printf("IssueReservationToken ERROR sending to email[%s], %s\n", user.Email, err.Error())
		return nil, fmt.Errorf("failed to send email to[%s]", user.Email)
	}
	return &tmpl, nil
}

func sendSignupCompletionEmail(tmpl template.VerifyEmailTemplate) error {
	subject := "User sing up confirmation"
	body, err := tmpl.Email()
	if err != nil {
		return err
	}
	err = sendEmail(tmpl.SourceEmail, tmpl.UserEmail, Email{
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
