package model

import "fmt"

type SSEvent string

const (
	Unknown              SSEvent = "unknown"
	PingEventName        SSEvent = "ping"
	ChatAddEventName     SSEvent = "chat-add"
	ChatDropEventName    SSEvent = "chat-drop"
	ChatCloseEventName   SSEvent = "chat-close"
	MessageAddEventName  SSEvent = "msg-add"
	MessageDropEventName SSEvent = "msg-drop"
)

type UpdateType int

const (
	UnknownUpdate  UpdateType = iota
	ChatCreated    UpdateType = iota
	ChatDeleted    UpdateType = iota
	ChatInvite     UpdateType = iota
	MessageAdded   UpdateType = iota
	MessageDeleted UpdateType = iota
)

func (u *UpdateType) String() string {
	switch *u {
	case ChatCreated, ChatInvite:
		return string(ChatAddEventName)
	case MessageAdded:
		return string(MessageAddEventName)
	case MessageDeleted:
		return string(MessageDropEventName)
	default:
		return string(Unknown)
	}
}

type LiveUpdate struct {
	Event  UpdateType
	ChatID int
	MsgID  int
	Author string
	Data   string
	Error  error
}

func (u *LiveUpdate) String() string {
	return fmt.Sprintf("LiveUpdate{event:%s,chat:%d,msg:%d,author:%s,data:%s,error:%v}",
		u.Event.String(), u.ChatID, u.MsgID, u.Author, u.Data, u.Error)
}
