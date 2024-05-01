package event

import (
	"fmt"
)

type SSEvent string

const (
	Unknown              SSEvent = "unknown"
	PingEventName        SSEvent = "ping"
	ChatAddEventName     SSEvent = "chat-add"
	ChatDropEventName    SSEvent = "chat-drop"
	ChatCloseEventName   SSEvent = "chat-close"
	ChatExpelEventName   SSEvent = "chat-expel"
	ChatLeaveEventName   SSEvent = "chat-leave"
	MessageAddEventName  SSEvent = "msg-add"
	MessageDropEventName SSEvent = "msg-drop"
)

func (e SSEvent) Format(
	chatId int,
	userId uint,
	msgId int,
) string {
	switch e {
	case MessageAddEventName:
		return fmt.Sprintf("%s-chat-%d", MessageAddEventName, chatId)
	case MessageDropEventName:
		return fmt.Sprintf("%s-chat-%d-msg-%d", MessageDropEventName, chatId, msgId)
	case ChatAddEventName:
		return string(ChatAddEventName)
	case ChatExpelEventName:
		return fmt.Sprintf("%s-%d-user-%d", ChatExpelEventName, chatId, userId)
	case ChatLeaveEventName:
		return fmt.Sprintf("%s-%d-user-%d", ChatLeaveEventName, chatId, userId)
	case ChatDropEventName:
		return fmt.Sprintf("%s-%d", ChatDropEventName, chatId)
	case ChatCloseEventName:
		return fmt.Sprintf("%s-%d", ChatCloseEventName, chatId)
	default:
		panic(fmt.Sprintf("unknown event type[%v]", e))
	}
}
