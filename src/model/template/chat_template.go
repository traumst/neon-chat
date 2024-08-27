package template

import (
	"bytes"
	"fmt"
	"html/template"

	"prplchat/src/model/event"
	ti "prplchat/src/model/template/interface"
)

type ChatTemplate struct {
	ChatId   uint
	ChatName string
	User     ti.Renderable
	Viewer   ti.Renderable
	Owner    ti.Renderable
	Users    []ti.Renderable
	Messages []ti.Renderable
}

func (c ChatTemplate) UserChangeEvent() string {
	return event.UserChange.FormatEventName(0, c.User.(UserTemplate).UserId, 0)
}

func (c ChatTemplate) ChatDropEvent() string {
	return event.ChatDrop.FormatEventName(c.ChatId, c.User.(UserTemplate).UserId, 0)
}

func (c ChatTemplate) ChatCloseEvent() string {
	return event.ChatClose.FormatEventName(c.ChatId, c.User.(UserTemplate).UserId, 0)
}

func (c ChatTemplate) ChatExpelEvent() string {
	return event.ChatExpel.FormatEventName(c.ChatId, c.User.(UserTemplate).UserId, 0)
}

func (c ChatTemplate) ChatLeaveEvent() string {
	return event.ChatLeave.FormatEventName(c.ChatId, c.User.(UserTemplate).UserId, 0)
}

func (c ChatTemplate) MessageAddEvent() string {
	return event.MessageAdd.FormatEventName(c.ChatId, c.User.(UserTemplate).UserId, 0)
}

func (c ChatTemplate) HTML() (string, error) {
	if err := c.validate(); err != nil {
		return "", fmt.Errorf("cannot template, %s", err.Error())
	}
	var buf bytes.Buffer
	chatTmpl := template.Must(template.ParseFiles(
		"static/html/chat/chat_div.html",
		"static/html/user_div.html",
		"static/html/chat/message_li.html",
		"static/html/avatar_div.html"))
	if err := chatTmpl.Execute(&buf, c); err != nil {
		return "", fmt.Errorf("failed to template, %s", err.Error())
	}
	return buf.String(), nil
}

func (c ChatTemplate) ShortHTML() (string, error) {
	if err := c.validateShort(); err != nil {
		return "", fmt.Errorf("cannot template, %s", err.Error())
	}
	var buf bytes.Buffer
	shortTmpl := template.Must(template.ParseFiles("static/html/nav/chat_li.html"))
	if err := shortTmpl.Execute(&buf, c); err != nil {
		return "", fmt.Errorf("failed to template, %s", err.Error())
	}
	return buf.String(), nil
}

func (c ChatTemplate) validate() error {
	if c.User.(UserTemplate).UserId < 1 {
		return fmt.Errorf("ChatTemplate requires User.UserId but is [%d], %v", c.User.(UserTemplate).UserId, c)
	}
	if c.User.(UserTemplate).UserName == "" {
		return fmt.Errorf("ChatTemplate requires User.UserName but is [%s], %v", c.User.(UserTemplate).UserName, c)
	}
	// if c.Viewer.UserId < 1 {
	// 	return fmt.Errorf("ChatTemplate requires Viewer.UserId but is [%d], %v", c.Viewer.UserId, c)
	// }
	// if c.Viewer.UserName == "" {
	// 	return fmt.Errorf("ChatTemplate requires Viewer.UserName but is [%s], %v", c.Viewer.UserName, c)
	// }
	// if c.Owner.UserId < 1 {
	// 	return fmt.Errorf("ChatTemplate requires Owner.UserId but is [%d], %v", c.Owner.UserId, c)
	// }
	// if len(c.Users) < 1 {
	// 	return fmt.Errorf("ChatTemplate requires Users but is empty, %v", c)
	// }
	return nil
}

func (c ChatTemplate) validateShort() error {
	if c.ChatId < 1 {
		return fmt.Errorf("ChatTemplate requires ChatId, but is [%d]", c.ChatId)
	}
	if c.ChatName == "" {
		return fmt.Errorf("ChatTemplate requires ChatName, but is [%s]", c.ChatName)
	}
	if c.Viewer.(UserTemplate).UserId < 1 {
		return fmt.Errorf("ChatTemplate requires Viewer.UserId, but is [%d]", c.Viewer.(UserTemplate).UserId)
	}
	if c.Viewer.(UserTemplate).UserName == "" {
		return fmt.Errorf("ChatTemplate requires Viewer.UserName, but is [%s]", c.Viewer.(UserTemplate).UserName)
	}
	return nil
}
