package template

import (
	"bytes"
	"html/template"
)

type WelcomeTemplate struct {
	User UserTemplate
}

func (w *WelcomeTemplate) HTML() (string, error) {
	var buf bytes.Buffer
	welcomeTmpl := template.Must(template.ParseFiles("static/html/welcome_div.html"))
	err := welcomeTmpl.Execute(&buf, w)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
