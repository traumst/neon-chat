package email

import (
	"bytes"
	"log"
	"text/template"
)

type InfoMessage struct {
	Header string
	Body   string
	Footer string
}

func (m *InfoMessage) InformUser() string {
	// TODO generic info message
	tmpl := template.Must(template.ParseFiles(
		"static/html/verification_sent.html"))
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, m)
	if err != nil {
		log.Printf("failed to template an html, %s", err.Error())
		return "ERROR_INFORM_USER"
	}
	return buf.String()
}
