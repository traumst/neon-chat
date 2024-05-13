package template

import (
	"bytes"
	"encoding/base64"
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

func (avatar *AvatarTemplate) Base64() string {
	return base64.StdEncoding.EncodeToString(avatar.Image)
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
