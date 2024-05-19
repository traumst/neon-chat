package handler

import (
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
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
) (*template.EmailSentTemplate, error) {
	token := utils.RandStringBytes(16)
	expire := time.Now().Add(1 * time.Hour)
	reserve, err := Reserve(db, user, token, expire)
	if err != nil {
		log.Printf("IssueReservationToken ERROR reserving[%s], %s\n", user.Email, err.Error())
		return nil, fmt.Errorf("")
	}
	emailConfig := app.SmtpConfig()
	err = sendSignupCompletionEmail(reserve, emailConfig, user.Email)
	if err != nil {
		log.Printf("IssueReservationToken ERROR sending to email[%s], %s\n", user.Email, err.Error())
		return nil, fmt.Errorf("failed to send email to[%s]", user.Email)
	}
	return &template.EmailSentTemplate{
		SourceEmail: emailConfig.User,
		UserEmail:   user.Email,
		UserName:    user.Name,
		Expire:      reserve.Expire.Format(time.RFC3339),
	}, nil
}

func sendSignupCompletionEmail(reserve *db.Reservation, sender SmtpConfig, receiver string) error {
	subject := "User sing up confirmation"
	body := fmt.Sprintf(
		"<p>Hellow, %s!"+
			"<p>Your verification link: %s"+
			"<p>Reservation expires at %s",
		receiver,
		reserve.Token,
		reserve.Expire.String())
	err := sendEmail(sender.User, sender.Pass, receiver, subject, body)
	return err
}

func sendEmail(sender, password, receiver, subject, body string) error {
	auth := smtp.PlainAuth("", sender, password, "smtp.gmail.com")

	headers := make(map[string]string)
	headers["From"] = sender
	headers["To"] = receiver
	headers["Subject"] = subject
	headers["Date"] = time.Now().Format(time.RFC1123Z)
	//headers["CC"] = "cc@example.com"
	//headers["BCC"] = "bcc@example.com"
	headers["Reply-To"] = "do-not-reply@please.com"

	msg := ""
	for k, v := range headers {
		msg += k + ": " + v + "\r\n"
	}
	msg += "\r\n" + body

	err := smtp.SendMail("smtp.gmail.com:587", auth, sender, []string{receiver}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email about [%s] to [%s], %s", subject, receiver, err.Error())
	}

	return nil
}
