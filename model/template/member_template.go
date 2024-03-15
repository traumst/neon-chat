package template

import (
	"bytes"
	"html/template"
)

type MemberTemplate struct {
	ChatID int
	Name   string
	Owner  string
	Viewer string
	User   string
}

func (c *MemberTemplate) ShortHTML() (string, error) {
	var buf bytes.Buffer
	shortTmpl := template.Must(template.ParseFiles("html/bits/member_div.html"))
	err := shortTmpl.Execute(&buf, c)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
