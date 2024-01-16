package models

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
)

type Message struct {
	ID     int
	Author string
	Text   string
}

var msgTmpl = template.Must(template.ParseFiles("views/message.html"))

func (m *Message) GetHTML() (string, error) {
	var buf bytes.Buffer
	err := msgTmpl.Execute(&buf, m)
	if err != nil {
		log.Printf("------ GetHTML ERROR template, %s\n", m.Log())
		return "", err
	}

	return buf.String(), nil
}

func (m *Message) Log() string {
	return fmt.Sprintf("Message{id:%d,author:[%s],text:[%s]}", m.ID, m.Author, m.Text)
}
