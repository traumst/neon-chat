package models

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
)

type MessageTemplate struct {
	ID         int
	Author     string
	Text       string
	ActiveUser string
}

var msgTmpl = template.Must(template.ParseFiles("views/message.html"))

func (m *MessageTemplate) GetHTML() (string, error) {
	log.Printf("------ MessageTemplate.GetHTML TRACE\n")
	var buf bytes.Buffer
	err := msgTmpl.Execute(&buf, m)
	if err != nil {
		log.Printf("------ GetHTML ERROR template, %s\n", m.Log())
		return "", err
	}

	return buf.String(), nil
}

func (m *MessageTemplate) Log() string {
	return fmt.Sprintf("MessageTemplate{id:%d,author:[%s],text:[%s]}", m.ID, m.Author, m.Text)
}
