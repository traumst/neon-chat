package models

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
)

type HomeTemplate struct {
	Chats        []*Chat
	OpenTemplate *ChatTemplate
	ActiveUser   string
}

var homeTmpl = template.Must(template.ParseFiles(
	"html/home.html",
	"html/chat.html",
	"html/chat_li.html"))

func (h *HomeTemplate) GetHTML() (string, error) {
	log.Printf("------ Home.GetHTML TRACE %s\n", h.Log())
	var buf bytes.Buffer
	err := homeTmpl.Execute(&buf, h)
	if err != nil {
		log.Printf("------ Home.GetHTML ERROR template, %s, %s\n", err, h.Log())
		return "", err
	}

	return buf.String(), nil
}

func (h *HomeTemplate) Log() string {
	openChatName := "nil"
	if h.OpenTemplate != nil {
		openChatName = h.OpenTemplate.Name
	}
	return fmt.Sprintf("Home{user:[%s],open:[%s],count:%d}", h.ActiveUser, openChatName, len(h.Chats))
}
