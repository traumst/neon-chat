package template

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
)

type AvatarTemplate struct {
	Id     uint
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
	if err := a.validate(); err != nil {
		return "", fmt.Errorf("cannot template, %s", err.Error())
	}
	var buf bytes.Buffer
	tmpl := template.Must(template.ParseFiles("static/html/avatar_div.html"))
	if err := tmpl.Execute(&buf, a); err != nil {
		return "", fmt.Errorf("failed to template, %s", err.Error())
	}
	return buf.String(), nil
}

func (a *AvatarTemplate) validate() error {
	if a.Id < 1 {
		return fmt.Errorf("AvatarTemplate id cannot be 0")
	}
	if a.Title == "" {
		return fmt.Errorf("AvatarTemplate title cannot be empty")
	}
	if a.UserId < 1 {
		return fmt.Errorf("AvatarTemplate user id cannot be 0")
	}
	if a.Size == "" {
		return fmt.Errorf("AvatarTemplate size cannot be empty")
	}
	if len(a.Image) < 1 {
		return fmt.Errorf("AvatarTemplate image cannot be empty")
	}
	if a.Mime == "" {
		return fmt.Errorf("AvatarTemplate mime cannot be empty")
	}
	return nil
}
