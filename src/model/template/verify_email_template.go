package template

import (
	"bytes"
	"fmt"
	"html/template"
)

type VerifyEmailTemplate struct {
	SourceEmail string
	UserEmail   string
	UserName    string
	Token       string
	TokenExpire string
}

func (vet VerifyEmailTemplate) HTML() (string, error) {
	if err := vet.validateHtml(); err != nil {
		return "", fmt.Errorf("cannot template, %s", err.Error())
	}
	var buf bytes.Buffer
	tmpl := template.Must(template.ParseFiles("static/html/verification_sent.html"))
	if err := tmpl.Execute(&buf, vet); err != nil {
		return "", fmt.Errorf("failed to template, %s", err.Error())
	}
	return buf.String(), nil
}

func (vet VerifyEmailTemplate) Email() (string, error) {
	if err := vet.validateEmail(); err != nil {
		return "", fmt.Errorf("cannot template, %s", err.Error())
	}
	var buf bytes.Buffer
	tmpl := template.Must(template.ParseFiles("static/html/email/verification_email.html"))
	if err := tmpl.Execute(&buf, vet); err != nil {
		return "", fmt.Errorf("failed to template, %s", err.Error())
	}
	return buf.String(), nil
}

func (vet VerifyEmailTemplate) validateHtml() error {
	if vet.SourceEmail == "" {
		return fmt.Errorf("VerifyEmailTemplate source email cannot be empty")
	}
	if vet.UserEmail == "" {
		return fmt.Errorf("VerifyEmailTemplate user email cannot be empty")
	}
	if vet.UserName == "" {
		return fmt.Errorf("VerifyEmailTemplate user name cannot be empty")
	}
	if vet.TokenExpire == "" {
		return fmt.Errorf("VerifyEmailTemplate token expire cannot be empty")
	}
	return nil
}

func (vet VerifyEmailTemplate) validateEmail() error {
	if vet.UserName == "" {
		return fmt.Errorf("VerifyEmailTemplate user name cannot be empty")
	}
	if vet.Token == "" {
		return fmt.Errorf("VerifyEmailTemplate token cannot be empty")
	}
	if vet.TokenExpire == "" {
		return fmt.Errorf("VerifyEmailTemplate token expire cannot be empty")
	}
	return nil
}
