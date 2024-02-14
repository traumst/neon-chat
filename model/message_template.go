package model

import (
	"bytes"
	"html/template"
	"log"
)

type MessageTemplate struct {
	ID         int
	Author     string
	Text       string
	ActiveUser string
}

func (m *MessageTemplate) GetHTML(reqId string) (string, error) {
	log.Printf("----%s---> MessageTemplate.GetHTML TRACE [%+v]\n", reqId, m)
	var buf bytes.Buffer
	msgTmpl := template.Must(template.ParseFiles("html/message.html"))
	err := msgTmpl.Execute(&buf, m)
	if err != nil {
		log.Printf("<---%s---- MessageTemplate.GetHTML ERROR template, %s, [%+v]\n", reqId, err, m)
		return "", err
	}

	log.Printf("<--%s--- MessageTemplate.GetHTML TRACE serve buf\n", reqId)
	return buf.String(), nil
}
