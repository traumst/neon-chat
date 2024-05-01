package event

type UpdateType string

const (
	UnknownUpdate  UpdateType = "unknown"
	ChatCreated    UpdateType = "chat_created"
	ChatDeleted    UpdateType = "chat_deleted"
	ChatClose      UpdateType = "chat_close"
	ChatInvite     UpdateType = "chat_invite"
	ChatExpel      UpdateType = "chat_expel"
	ChatLeave      UpdateType = "chat_leave"
	MessageAdded   UpdateType = "message_added"
	MessageDeleted UpdateType = "message_deleted"
)
