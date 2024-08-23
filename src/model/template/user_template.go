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

func (ut UserTemplate) UserChangeEvent() string {
	return event.UserChange.FormatEventName(0, ut.UserId, 0)
}

func (ut UserTemplate) ChatExpelEvent() string {
	return event.ChatExpel.FormatEventName(ut.ChatId, ut.UserId, 0)
}

func (ut UserTemplate) ChatLeaveEvent() string {
	return event.ChatLeave.FormatEventName(ut.ChatId, ut.UserId, 0)
}

func (ut UserTemplate) HTML() (string, error) {
	if err := ut.validate(); err != nil {
		return "", fmt.Errorf("cannot template, %s", err.Error())
	}
	var buf bytes.Buffer
	tmpl := template.Must(template.ParseFiles("static/html/user_div.html"))
	if err := tmpl.Execute(&buf, ut); err != nil {
		return "", fmt.Errorf("failed to template, %s", err.Error())
	}
	return buf.String(), nil
}

func (ut UserTemplate) ShortHTML() (string, error) {
	if err := ut.validate(); err != nil {
		return "", fmt.Errorf("cannot template, %s", err.Error())
	}
	var buf bytes.Buffer
	tmpl := template.Must(template.ParseFiles("static/html/search/user_option.html"))
	if err := tmpl.Execute(&buf, ut); err != nil {
		return "", fmt.Errorf("failed to template, %s", err.Error())
	}
	return buf.String(), nil
}

func (ut UserTemplate) validate() error {
	if ut.ChatId < 1 {
		return fmt.Errorf("UserTemplate chat id cannot be 0")
	}
	if ut.ChatOwnerId < 1 {
		return fmt.Errorf("UserTemplate chat owner id cannot be 0")
	}
	if ut.UserId < 1 {
		return fmt.Errorf("UserTemplate user id cannot be 0")
	}
	if ut.UserName == "" {
		return fmt.Errorf("UserTemplate user name cannot be empty")
	}
	if ut.UserEmail == "" {
		return fmt.Errorf("UserTemplate user email cannot be empty")
	}
	if ut.ViewerId < 1 {
		return fmt.Errorf("UserTemplate viewer id cannot be 0")
	}
	return nil
}
