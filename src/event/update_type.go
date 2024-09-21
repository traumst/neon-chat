package event

import "fmt"

type EventType string

const (
	Ping         EventType = "ping"
	UserChange   EventType = "user-changed"
	AvatarChange EventType = "avatar-changed"

	ChatAdd   EventType = "chat-add"
	ChatDrop  EventType = "chat-drop"
	ChatClose EventType = "chat-close"

	ChatInvite EventType = "chat-invite"
	ChatExpel  EventType = "chat-expel"
	ChatLeave  EventType = "chat-leave"

	MessageAdd  EventType = "msg-add"
	MessageDrop EventType = "msg-drop"
)

func (e EventType) FormatEventName(
	chatId uint,
	userId uint,
	msgId uint,
) string {
	switch e {
	case UserChange:
		return fmt.Sprintf("%s-%d", UserChange, userId)
	case AvatarChange:
		return fmt.Sprintf("%s-%d", AvatarChange, userId)

	case ChatAdd:
		return string(ChatAdd)
	case ChatDrop:
		return fmt.Sprintf("%s-%d", ChatDrop, chatId)
	case ChatClose:
		return fmt.Sprintf("%s-%d", ChatClose, chatId)

	case ChatInvite:
		return string(ChatInvite)
	case ChatExpel:
		return fmt.Sprintf("%s-%d-user-%d", ChatExpel, chatId, userId)
	case ChatLeave:
		return fmt.Sprintf("%s-%d-user-%d", ChatLeave, chatId, userId)

	case MessageAdd:
		return fmt.Sprintf("%s-chat-%d", MessageAdd, chatId)
	case MessageDrop:
		return fmt.Sprintf("%s-chat-%d-msg-%d", MessageDrop, chatId, msgId)

	default:
		panic(fmt.Sprintf("unknown event type[%v]", e))
	}
}
