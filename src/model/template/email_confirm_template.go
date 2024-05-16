package template

import (
	"bytes"
	"html/template"
)

type EmailConfirmTemplate struct {
	SourceEmail string
	UserEmail   string
	UserName    string
	Expire      string
}

func (w *EmailConfirmTemplate) HTML() (string, error) {
	var buf bytes.Buffer
	welcomeTmpl := template.Must(template.ParseFiles("static/html/chat/email_confirm_div.html"))
	err := welcomeTmpl.Execute(&buf, w)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
