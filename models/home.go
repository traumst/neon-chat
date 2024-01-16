package models

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
)

type Home struct {
	Chats      []*Chat
	OpenChat   *Chat
	ActiveUser string
}

var homeTmpl = template.Must(template.ParseFiles("views/home.html", "views/chat.html"))

func (h *Home) GetHTML() (string, error) {
	var buf bytes.Buffer
	err := homeTmpl.Execute(&buf, h)
	if err != nil {
		log.Printf("------ GetHTML ERROR template, %s\n", h.Log())
		return "", err
	}

	return buf.String(), nil
}

func (h *Home) Log() string {
	return fmt.Sprintf("Home{user:[%s],open:[%s],count:%d}", h.ActiveUser, h.OpenChat.Name, len(h.Chats))
}
