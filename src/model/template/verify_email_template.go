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

func (w *VerifyEmailTemplate) HTML() (string, error) {
	if err := w.validateHtml(); err != nil {
		return "", fmt.Errorf("cannot template, %s", err.Error())
	}
	var buf bytes.Buffer
	welcomeTmpl := template.Must(template.ParseFiles("static/html/verification_sent.html"))
	if err := welcomeTmpl.Execute(&buf, w); err != nil {
		return "", fmt.Errorf("failed to template, %s", err.Error())
	}
	return buf.String(), nil
}

func (w *VerifyEmailTemplate) Email() (string, error) {
	if err := w.validateEmail(); err != nil {
		return "", fmt.Errorf("cannot template, %s", err.Error())
	}
	var buf bytes.Buffer
	welcomeTmpl := template.Must(template.ParseFiles("static/html/email/verification_email.html"))
	if err := welcomeTmpl.Execute(&buf, w); err != nil {
		return "", fmt.Errorf("failed to template, %s", err.Error())
	}
	return buf.String(), nil
}

func (w *VerifyEmailTemplate) validateHtml() error {
	if w.SourceEmail == "" {
		return fmt.Errorf("VerifyEmailTemplate source email cannot be empty")
	}
	if w.UserEmail == "" {
		return fmt.Errorf("VerifyEmailTemplate user email cannot be empty")
	}
	if w.UserName == "" {
		return fmt.Errorf("VerifyEmailTemplate user name cannot be empty")
	}
	if w.TokenExpire == "" {
		return fmt.Errorf("VerifyEmailTemplate token expire cannot be empty")
	}
	return nil
}

func (w *VerifyEmailTemplate) validateEmail() error {
	if w.UserName == "" {
		return fmt.Errorf("VerifyEmailTemplate user name cannot be empty")
	}
	if w.Token == "" {
		return fmt.Errorf("VerifyEmailTemplate token cannot be empty")
	}
	if w.TokenExpire == "" {
		return fmt.Errorf("VerifyEmailTemplate token expire cannot be empty")
	}
	return nil
}
