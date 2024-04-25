package template

import (
	"bytes"
	"html/template"
)

type AuthForm struct {
	Id    string
	Label string
	Title string
}

type LoginTemplate struct {
	Login     AuthForm
	Signup    AuthForm
	LoadLocal bool
}

func (lt *LoginTemplate) HTML() (string, error) {
	template := template.Must(template.ParseFiles("static/html/login_div.html"))
	var buf bytes.Buffer
	err := template.Execute(&buf, lt)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
