package models

import (
	"bytes"
	"log"
	"text/template"
)

type ChatTemplate struct {
	Chat       *Chat
	ActiveUser string
}

var chatTmpl = template.Must(template.ParseFiles(
	"views/chat.html",
	"views/message.html"))

func (ct *ChatTemplate) GetHTML() (string, error) {
	log.Printf("------ ChatTemplate.GetHTML TRACE enter user[%s] into chat[%s]\n", ct.ActiveUser, ct.Chat.Log())
	var buf bytes.Buffer
	err := chatTmpl.Execute(&buf, ct)
	if err != nil {
		log.Printf("------ GetHTML ERROR template, %s\n", ct.Chat.Log())
		return "", err
	}
	log.Printf("------ ChatTemplate.GetHTML TRACE serve buf\n")
	return buf.String(), nil
}
