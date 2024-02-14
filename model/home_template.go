package model

import (
	"bytes"
	"html/template"
	"log"
)

type HomeTemplate struct {
	Chats        []*Chat
	OpenTemplate *ChatTemplate
	ActiveUser   string
}

func (h *HomeTemplate) GetHTML(reqId string) (string, error) {
	log.Printf("----%s---> Home.GetHTML TRACE %+v\n", reqId, h)
	var buf bytes.Buffer
	homeTmpl := template.Must(template.ParseFiles(
		"html/home.html",
		"html/welcome.html",
		"html/chat.html",
		"html/chat_li.html"))
	err := homeTmpl.Execute(&buf, h)
	if err != nil {
		log.Printf("<---%s---- Home.GetHTML ERROR template, %s, [%+v]\n", reqId, err, h)
		return "", err
	}

	log.Printf("<--%s--- Home.GetHTML TRACE serve buf\n", reqId)
	return buf.String(), nil
}
