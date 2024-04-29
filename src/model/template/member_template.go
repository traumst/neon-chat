package template

import (
	"bytes"
	"html/template"
)

type MemberTemplate struct {
	ChatId         int
	ChatName       string
	User           UserTemplate
	Viewer         UserTemplate
	Owner          UserTemplate
	ChatExpelEvent string
}

func (c *MemberTemplate) ShortHTML() (string, error) {
	var buf bytes.Buffer
	shortTmpl := template.Must(template.ParseFiles("static/html/chat/member_div.html"))
	err := shortTmpl.Execute(&buf, c)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
