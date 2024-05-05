package event

import "fmt"

type UpdateType string

const (
	Ping           UpdateType = "ping"
	UserChanged    UpdateType = "user-changed"
	ChatCreated    UpdateType = "chat-add"
	ChatDeleted    UpdateType = "chat-drop"
	ChatClose      UpdateType = "chat-close"
	ChatInvite     UpdateType = "chat-invite"
	ChatExpel      UpdateType = "chat-expel"
	ChatLeave      UpdateType = "chat-leave"
	MessageAdded   UpdateType = "msg-add"
	MessageDeleted UpdateType = "msg-drop"
)

func (e UpdateType) FormatEventName(
	chatId int,
	userId uint,
	msgId int,
) string {
	switch e {
	case UserChanged:
		return fmt.Sprintf("%s-%d", UserChanged, userId)
	case ChatCreated:
		return string(ChatCreated)
	case ChatExpel:
		return fmt.Sprintf("%s-%d-user-%d", ChatExpel, chatId, userId)
	case ChatLeave:
		return fmt.Sprintf("%s-%d-user-%d", ChatLeave, chatId, userId)
	case ChatDeleted:
		return fmt.Sprintf("%s-%d", ChatDeleted, chatId)
	case ChatClose:
		return fmt.Sprintf("%s-%d", ChatClose, chatId)
	case MessageAdded:
		return fmt.Sprintf("%s-chat-%d", MessageAdded, chatId)
	case MessageDeleted:
		return fmt.Sprintf("%s-chat-%d-msg-%d", MessageDeleted, chatId, msgId)
	default:
		panic(fmt.Sprintf("unknown event type[%v]", e))
	}
}
