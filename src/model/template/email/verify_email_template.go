package email

import (
	"bytes"
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
	var buf bytes.Buffer
	welcomeTmpl := template.Must(template.ParseFiles("static/html/verification_sent.html"))
	err := welcomeTmpl.Execute(&buf, w)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
