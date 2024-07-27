package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

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

func (db *DBConn) GetUserChatIds(userId uint) ([]uint, error) {
	if userId == 0 {
		return nil, fmt.Errorf("bad input: userId[%d]", userId)
	}
	if !db.ConnIsActive() {
		return nil, fmt.Errorf("db is not connected")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	var userChatIds []uint
	err := db.conn.Select(&userChatIds, `SELECT chat_id FROM chat_users WHERE user_id = ?`, userId)
	if err != nil {
		return nil, fmt.Errorf("error adding user: %s", err.Error())
	}
	return userChatIds, nil
}

func (db *DBConn) GetUserChats(userId uint) ([]Chat, error) {
	if userId == 0 {
		return nil, fmt.Errorf("invalid userId[%d]", userId)
	}
	if !db.ConnIsActive() {
		return nil, fmt.Errorf("db is not connected")
	}

	chatIds, err := db.GetUserChatIds(userId)
	if err != nil {
		return nil, fmt.Errorf("error getting chat ids for user[%d]: %s", userId, err)
	}
	if len(chatIds) == 0 {
		return []Chat{}, nil
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	query, args, err := sqlx.In(`SELECT * FROM chats WHERE id IN (?)`, chatIds)
	if err != nil {
		return nil, fmt.Errorf("error preparing query for user[%d]: %s", userId, err)
	}
	query = db.conn.Rebind(query)

	var userChats []Chat
	err = db.conn.Select(&userChats, query, args...)
	//err = db.conn.Select(&userChats, `SELECT * FROM chats WHERE id IN (?)`, chatIds)
	if err != nil {
		return nil, fmt.Errorf("error getting chats for user[%d]: %s", userId, err)
	}
	return userChats, nil
}

func (db *DBConn) RemoveChatUser(chatId uint, userId uint) error {
	if chatId == 0 || userId == 0 {
		return fmt.Errorf("bad input: chatId[%d], userId[%d]", chatId, userId)
	}
	if !db.ConnIsActive() {
		return fmt.Errorf("db is not connected")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	_, err := db.conn.Exec(`DELETE FROM chat_users WHERE chat_id = ? AND user_id = ?`, chatId, userId)
	if err != nil {
		return fmt.Errorf("error removing user: %s", err.Error())
	}
	return nil
}
