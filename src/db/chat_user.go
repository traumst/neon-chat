package db

import "fmt"

type ChatUser struct {
	ChatId uint `db:"chat_id"`
	UserId uint `db:"user_id"`
}

const ChatUserSchema = `
	CREATE TABLE IF NOT EXISTS chat_users (
		chat_id INTEGER,
		user_id INTEGER,
		FOREIGN KEY(chat_id) REFERENCES chats(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`
const ChatUserIndex = `CREATE INDEX IF NOT EXISTS idx_chat_users_chat_id ON chat_users(chat_id);
CREATE INDEX IF NOT EXISTS idx_chat_users_user_id ON chat_users(user_id);`

func (db *DBConn) ChatUserTableExists() bool {
	return db.TableExists("chat_users")
}

func (db *DBConn) AddChatUser(chatId uint, userId uint) error {
	if chatId == 0 || userId == 0 {
		return fmt.Errorf("bad input: chatId[%d], userId[%d]", chatId, userId)
	}
	if !db.ConnIsActive() {
		return fmt.Errorf("db is not connected")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	_, err := db.conn.Exec(`INSERT INTO chat_users (chat_id, user_id) VALUES (?, ?)`, chatId, userId)
	if err != nil {
		return fmt.Errorf("error adding user: %s", err.Error())
	}
	return nil
}
