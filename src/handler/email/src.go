package email

import (
	"net/mail"

	"neon-chat/src/template"
)

func IsEmailValid(address string) bool {
	_, err := mail.ParseAddress(address)
	return err == nil
}

func SendSignupCompletionEmail(
	tmpl template.VerifyEmailTemplate,
	source string,
	pass string,
) error {
	subject := "User sing up confirmation"
	body, err := tmpl.Email()
	if err != nil {
		return err
	}
	email := email{
		sender:    tmpl.SourceEmail,
		receivers: []string{tmpl.UserEmail},
		subject:   subject,
		body:      body,
	}
	err = email.sendAs(source, pass)
	return err
}
