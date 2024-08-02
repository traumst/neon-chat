package db

import (
	"fmt"
)

const minChatTitleLen = 5
const maxChatTitleLen = 256

type Chat struct {
	Id      uint   `db:"id"`
	Title   string `db:"title"`
	OwnerId uint   `db:"owner_id"`
}

const ChatSchema = `
	CREATE TABLE IF NOT EXISTS chats (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		title TEXT,
		owner_id INTEGER,
		FOREIGN KEY(owner_id) REFERENCES users(id)
	);`
const ChatIndex = `CREATE INDEX IF NOT EXISTS idx_chat_owner_id ON chats(owner_id);`

func (db *DBConn) ChatTableExists() bool {
	return db.TableExists("chats")
}

func (db *DBConn) AddChat(chat *Chat) (*Chat, error) {
	if chat.Id != 0 {
		return nil, fmt.Errorf("chat already has an id[%d]", chat.Id)
	} else if len(chat.Title) < minChatTitleLen || len(chat.Title) > maxChatTitleLen {
		return nil, fmt.Errorf("chat has no title")
	} else if chat.OwnerId == 0 {
		return nil, fmt.Errorf("chat has no owner")
	}
	if !db.ConnIsActive() {
		return nil, fmt.Errorf("db is not connected")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	result, err := db.conn.Exec(`INSERT INTO chats (title, owner_id) VALUES (?, ?)`, chat.Title, chat.OwnerId)
	if err != nil {
		return nil, fmt.Errorf("error adding user: %s", err)
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %s", err)
	}
	chat.Id = uint(lastId)
	return chat, nil
}

func (db *DBConn) GetChat(chatId uint) (*Chat, error) {
	if chatId == 0 {
		return nil, fmt.Errorf("bad input: chatId[%d]", chatId)
	}
	if !db.ConnIsActive() {
		return nil, fmt.Errorf("db is not connected")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	var chat Chat
	err := db.conn.Get(&chat, `SELECT * FROM chats WHERE id = ?`, chatId)
	if err != nil {
		return nil, fmt.Errorf("error getting chat: %s", err)
	}
	return &chat, nil
}

func (db *DBConn) UserCanChat(chatId uint, userId uint) (bool, error) {
	if !db.ConnIsActive() {
		return false, fmt.Errorf("db is not connected")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	var chatUser *ChatUser
	err := db.conn.Get(&chatUser, `SELECT * FROM chat_users where chat_id = ? and user_id = ?`, chatId, userId)
	if err != nil {
		return false, fmt.Errorf("error getting chats: %s", err)
	}
	return true, nil
}

func (db *DBConn) DeleteChat(chatId uint) error {
	if !db.ConnIsActive() {
		return fmt.Errorf("db is not connected")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	_, err := db.conn.Exec(`DELETE FROM chats WHERE id = ?`, chatId)
	if err != nil {
		return fmt.Errorf("error deleting chat[%d]: %s", chatId, err)
	}
	return nil
}
