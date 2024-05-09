package event

import "fmt"

type UpdateType string

const (
	Ping       UpdateType = "ping"
	UserChange UpdateType = "user-changed"

	ChatAdd   UpdateType = "chat-add"
	ChatDrop  UpdateType = "chat-drop"
	ChatClose UpdateType = "chat-close"

	ChatInvite UpdateType = "chat-invite"
	ChatExpel  UpdateType = "chat-expel"
	ChatLeave  UpdateType = "chat-leave"

	MessageAdd  UpdateType = "msg-add"
	MessageDrop UpdateType = "msg-drop"
)

func (e UpdateType) FormatEventName(
	chatId int,
	userId uint,
	msgId int,
) string {
	switch e {
	case UserChange:
		return fmt.Sprintf("%s-%d", UserChange, userId)

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
