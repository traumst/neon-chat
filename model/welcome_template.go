package model

import (
	"bytes"
	"html/template"
)

type WelcomeTemplate struct {
	ActiveUser string
}

func (w *WelcomeTemplate) GetHTML() (string, error) {
	var buf bytes.Buffer
	welcomeTmpl := template.Must(template.ParseFiles("html/bits/welcome_div.html"))
	err := welcomeTmpl.Execute(&buf, w)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
