package template

import (
	"bytes"
	"html/template"
)

type AvatarTemplate struct {
	Id     int
	Title  string
	UserId uint
	Size   string
	Image  []byte
	Mime   string
}

func (a *AvatarTemplate) HTML() (string, error) {
	var buf bytes.Buffer
	tmpl := template.Must(template.ParseFiles("static/html/avatar_div.html"))
	err := tmpl.Execute(&buf, a)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
