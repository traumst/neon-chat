package event

type SSEvent string

const (
	Unknown              SSEvent = "unknown"
	PingEventName        SSEvent = "ping"
	ChatAddEventName     SSEvent = "chat-add"
	ChatDropEventName    SSEvent = "chat-drop"
	ChatCloseEventName   SSEvent = "chat-close"
	ChatExpelEventName   SSEvent = "chat-expel"
	MessageAddEventName  SSEvent = "msg-add"
	MessageDropEventName SSEvent = "msg-drop"
)
