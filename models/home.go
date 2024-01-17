package models

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
)

type Home struct {
	Chats        []*Chat
	OpenTemplate *ChatTemplate
	ActiveUser   string
}

var homeTmpl = template.Must(template.ParseFiles(
	"views/home.html",
	"views/chat.html"))

func (h *Home) GetHTML() (string, error) {
	log.Printf("------ Home.GetHTML TRACE %s\n", h.Log())
	var buf bytes.Buffer
	err := homeTmpl.Execute(&buf, h)
	if err != nil {
		log.Printf("------ Home.GetHTML ERROR template, %s\n", h.Log())
		return "", err
	}

	return buf.String(), nil
}

// TODO return template
func (h *Home) Log() string {
	openChatName := "nil"
	if h.OpenTemplate != nil && h.OpenTemplate.Chat != nil {
		openChatName = h.OpenTemplate.Chat.Name
	}
	return fmt.Sprintf("Home{user:[%s],open:[%s],count:%d}", h.ActiveUser, openChatName, len(h.Chats))
}
