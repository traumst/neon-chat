package template

import (
	"bytes"
	"html/template"
)

type WelcomeTemplate struct {
	User UserTemplate // default user will be served generic message
}

func (wt WelcomeTemplate) HTML() (string, error) {
	var buf bytes.Buffer
	tmpl := template.Must(template.ParseFiles("static/html/chat/welcome_div.html"))
	err := tmpl.Execute(&buf, wt)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
