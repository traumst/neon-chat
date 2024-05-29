package email

import (
	"bytes"
	"text/template"
)

type EmailTemplate struct {
	Header string
	Body   string
	Footer string
}

func (w *VerifyEmailTemplate) Email() (string, error) {
	var buf bytes.Buffer
	welcomeTmpl := template.Must(template.ParseFiles("static/html/email/verification_email.html"))
	err := welcomeTmpl.Execute(&buf, w)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
