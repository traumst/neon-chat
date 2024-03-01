package model

type UpdateType int
type LiveUpdate struct {
	Event  UpdateType
	ChatID int
	MsgID  int
	Author string
	Data   string
}

const (
	UnknownUpdate  UpdateType = iota
	ChatCreated    UpdateType = iota
	ChatInvite     UpdateType = iota
	MessageAdded   UpdateType = iota
	MessageDeleted UpdateType = iota
)

func (u *UpdateType) String() string {
	switch *u {
	case ChatCreated:
		return string(ChatAddEventName)
	case ChatInvite:
		return string(ChatInviteEventName)
	case MessageAdded:
		return string(MessageAddEventName)
	case MessageDeleted:
		return string(MessageDropEventName)
	default:
		return string(Unknown)
	}
}
