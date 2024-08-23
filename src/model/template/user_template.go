package template

import (
	"bytes"
	"fmt"
	"html/template"

	"prplchat/src/model/event"
)

type UserTemplate struct {
	ChatId      uint
	ChatOwnerId uint
	UserId      uint
	UserName    string
	UserEmail   string
	//UserStatus  string
	ViewerId uint
}

func (c *UserTemplate) UserChangeEvent() string {
	return event.UserChange.FormatEventName(0, c.UserId, 0)
}

func (c *UserTemplate) ChatExpelEvent() string {
	return event.ChatExpel.FormatEventName(c.ChatId, c.UserId, 0)
}

func (c *UserTemplate) ChatLeaveEvent() string {
	return event.ChatLeave.FormatEventName(c.ChatId, c.UserId, 0)
}

func (c *UserTemplate) HTML() (string, error) {
	if err := c.validate(); err != nil {
		return "", fmt.Errorf("cannot template, %s", err.Error())
	}
	var buf bytes.Buffer
	shortTmpl := template.Must(template.ParseFiles("static/html/user_div.html"))
	if err := shortTmpl.Execute(&buf, c); err != nil {
		return "", fmt.Errorf("failed to template, %s", err.Error())
	}
	return buf.String(), nil
}

func (c *UserTemplate) ShortHTML() (string, error) {
	if err := c.validate(); err != nil {
		return "", fmt.Errorf("cannot template, %s", err.Error())
	}
	var buf bytes.Buffer
	shortTmpl := template.Must(template.ParseFiles("static/html/search/user_option.html"))
	if err := shortTmpl.Execute(&buf, c); err != nil {
		return "", fmt.Errorf("failed to template, %s", err.Error())
	}
	return buf.String(), nil
}

func (c *UserTemplate) validate() error {
	if c.ChatId < 1 {
		return fmt.Errorf("UserTemplate chat id cannot be 0")
	}
	if c.ChatOwnerId < 1 {
		return fmt.Errorf("UserTemplate chat owner id cannot be 0")
	}
	if c.UserId < 1 {
		return fmt.Errorf("UserTemplate user id cannot be 0")
	}
	if c.UserName == "" {
		return fmt.Errorf("UserTemplate user name cannot be empty")
	}
	if c.UserEmail == "" {
		return fmt.Errorf("UserTemplate user email cannot be empty")
	}
	if c.ViewerId < 1 {
		return fmt.Errorf("UserTemplate viewer id cannot be 0")
	}
	return nil
}
