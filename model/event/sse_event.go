package event

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
