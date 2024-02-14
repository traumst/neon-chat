package model

import (
	"bytes"
	"html/template"
	"log"
)

type ChatTemplate struct {
	ID         int
	Name       string
	ActiveUser string
	Users      []string
	Messages   []MessageTemplate
}

func (ct *ChatTemplate) GetHTML(reqId string) (string, error) {
	log.Printf("---%s--> ChatTemplate.GetHTML TRACE enter user[%s] into chat[%+v]\n", reqId, ct.ActiveUser, ct)
	var buf bytes.Buffer
	chatTmpl := template.Must(template.ParseFiles("html/chat.html", "html/message.html"))
	err := chatTmpl.Execute(&buf, ct)
	if err != nil {
		log.Printf("<--%s--- ChatTemplate.GetHTML ERROR template, [%s]\n", reqId, err)
		return "", err
	}
	log.Printf("<--%s--- ChatTemplate.GetHTML TRACE serve buf\n", reqId)
	return buf.String(), nil
}

func (ct *ChatTemplate) GetShortHTML(reqId string) (string, error) {
	log.Printf("---%s--> ChatTemplate.GetShortHTML TRACE enter user[%s] into chat[%+v]\n", reqId, ct.ActiveUser, ct)
	var buf bytes.Buffer
	shortTmpl := template.Must(template.ParseFiles("html/chat_li.html"))
	err := shortTmpl.Execute(&buf, ct)
	if err != nil {
		log.Printf("<--%s--- ChatTemplate.GetShortHTML ERROR template, [%s]\n", reqId, err)
		return "", err
	}
	log.Printf("<--%s--- ChatTemplate.GetShortHTML TRACE serve buf\n", reqId)
	return buf.String(), nil
}
