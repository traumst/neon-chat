package template

import (
	"bytes"
	"fmt"
	"html/template"
)

type AuthTemplate struct{}

func (at AuthTemplate) GetId() uint {
	return 0 // ie no user id here
}

func (at AuthTemplate) HTML() (string, error) {
	template := template.Must(template.ParseFiles("static/html/nav/auth_div.html"))
	var buf bytes.Buffer
	err := template.Execute(&buf, at)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (at AuthTemplate) ShortHTML() (string, error) {
	return "", fmt.Errorf("not implemented")
}
