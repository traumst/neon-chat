package event

import "fmt"

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
