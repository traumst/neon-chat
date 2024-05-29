package template

import (
	"bytes"
	"html/template"
)

type AuthTemplate struct{}

func (lt *AuthTemplate) HTML() (string, error) {
	template := template.Must(template.ParseFiles("static/html/nav/auth_div.html"))
	var buf bytes.Buffer
	err := template.Execute(&buf, lt)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
