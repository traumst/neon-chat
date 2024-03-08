package model

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

type LiveUpdate struct {
	Event  UpdateType
	ChatID int
	MsgID  int
	Author string
	Data   string
	Error  error
}

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
