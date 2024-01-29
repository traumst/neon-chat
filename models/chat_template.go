package models

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strings"
)

type ChatTemplate struct {
	ID         int
	Name       string
	ActiveUser string
	Users      []string
	Messages   []MessageTemplate
}

func (c *ChatTemplate) Log() string {
	return fmt.Sprintf("ChatTemplate{id:%d,name:[%s],active_user:[%s],users:[%s]}", c.ID, c.Name, c.ActiveUser, strings.Join(c.Users, ","))
}

var chatTmpl = template.Must(template.ParseFiles("views/chat.html", "views/message.html"))
var shortTmpl = template.Must(template.ParseFiles("views/chat_li.html"))

func (ct *ChatTemplate) GetHTML() (string, error) {
	log.Printf("------ ChatTemplate.GetHTML TRACE enter user[%s] into chat[%s]\n", ct.ActiveUser, ct.Log())
	var buf bytes.Buffer
	err := chatTmpl.Execute(&buf, ct)
	if err != nil {
		log.Printf("------ ChatTemplate.GetHTML ERROR template, %s\n", ct.Log())
		return "", err
	}
	log.Printf("------ ChatTemplate.GetHTML TRACE serve buf\n")
	return buf.String(), nil
}

func (ct *ChatTemplate) GetShortHTML() (string, error) {
	log.Printf("------ ChatTemplate.GetShortHTML TRACE\n")
	var buf bytes.Buffer
	err := shortTmpl.Execute(&buf, ct)
	if err != nil {
		log.Printf("------ ChatTemplate.GetShortHTML ERROR template, %s\n", shortTmpl.Name())
		return "", err
	}
	log.Printf("------ ChatTemplate.GetShortHTML TRACE serve buf\n")
	return buf.String(), nil
}
