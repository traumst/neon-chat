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

func (e SSEvent) Format(chatID int, msgID int) string {
	switch e {
	case ChatAddEventName:
		return string(ChatAddEventName)
	case ChatDropEventName:
		return fmt.Sprintf("%s-%d", ChatDropEventName, chatID)
	case ChatCloseEventName:
		return fmt.Sprintf("%s-%d", ChatCloseEventName, chatID)
	case MessageAddEventName:
		return fmt.Sprintf("%s-chat-%d", MessageAddEventName, chatID)
	case MessageDropEventName:
		return fmt.Sprintf("%s-chat-%d-msg-%d", MessageDropEventName, chatID, msgID)
	default:
		panic(fmt.Sprintf("unknown event type[%v]", e))
	}
}

type UpdateType int

const (
	UnknownUpdate  UpdateType = iota
	ChatCreated    UpdateType = iota
	ChatDeleted    UpdateType = iota
	ChatClose      UpdateType = iota
	ChatInvite     UpdateType = iota
	MessageAdded   UpdateType = iota
	MessageDeleted UpdateType = iota
)

func (u *UpdateType) String() SSEvent {
	switch *u {
	case ChatCreated, ChatInvite:
		return ChatAddEventName
	case MessageAdded:
		return MessageAddEventName
	case MessageDeleted:
		return MessageDropEventName
	case ChatDeleted:
		return ChatDropEventName
	case ChatClose:
		return ChatCloseEventName
	default:
		panic(fmt.Sprintf("unknown update type[%d]", *u))
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
	return fmt.Sprintf("LiveUpdate{event:%s,chat:%d,msg:%d,author:%s,error:%v}",
		u.Event.String(), u.ChatID, u.MsgID, u.Author, u.Error)
}
