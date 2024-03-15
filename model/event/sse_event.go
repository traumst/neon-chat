package event

type SSEvent string

const (
	Unknown               SSEvent = "unknown"
	PingEventName         SSEvent = "ping"
	ChatAddEventName      SSEvent = "chat-add"
	ChatDropEventName     SSEvent = "chat-drop"
	ChatCloseEventName    SSEvent = "chat-close"
	ChatUserDropEventName SSEvent = "chat-user-drop" // TODO in progress
	MessageAddEventName   SSEvent = "msg-add"
	MessageDropEventName  SSEvent = "msg-drop"
)
