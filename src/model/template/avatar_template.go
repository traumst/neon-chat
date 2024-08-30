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

func (at AvatarTemplate) GetId() uint {
	return at.Id
}

func (at AvatarTemplate) ShortHTML() (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (at AvatarTemplate) Base64() string {
	return base64.StdEncoding.EncodeToString(at.Image)
}

func (at AvatarTemplate) HTML() (string, error) {
	if err := at.validate(); err != nil {
		return "", fmt.Errorf("cannot template, %s", err.Error())
	}
	var buf bytes.Buffer
	tmpl := template.Must(template.ParseFiles("static/html/avatar_div.html"))
	if err := tmpl.Execute(&buf, at); err != nil {
		return "", fmt.Errorf("failed to template, %s", err.Error())
	}
	return buf.String(), nil
}

func (at AvatarTemplate) validate() error {
	if at.Id < 1 {
		return fmt.Errorf("AvatarTemplate id cannot be 0")
	}
	if at.Title == "" {
		return fmt.Errorf("AvatarTemplate title cannot be empty")
	}
	if at.UserId < 1 {
		return fmt.Errorf("AvatarTemplate user id cannot be 0")
	}
	if at.Size == "" {
		return fmt.Errorf("AvatarTemplate size cannot be empty")
	}
	if len(at.Image) < 1 {
		return fmt.Errorf("AvatarTemplate image cannot be empty")
	}
	if at.Mime == "" {
		return fmt.Errorf("AvatarTemplate mime cannot be empty")
	}
	return nil
}
