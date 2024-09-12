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

// TODO chat search
const ChatIndex = `CREATE INDEX IF NOT EXISTS idx_chat_title ON chats(title);`

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
	if db.tx == nil {
		return nil, fmt.Errorf("db has no transaction")
	}

	result, err := db.tx.Exec(`INSERT INTO chats (title, owner_id) VALUES (?, ?)`, chat.Title, chat.OwnerId)
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

func (db *DBConn) GetOwner(chatId uint) (*User, error) {
	if chatId == 0 {
		return nil, fmt.Errorf("bad input: chatId[%d]", chatId)
	}
	if !db.ConnIsActive() {
		return nil, fmt.Errorf("db is not connected")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	var user User
	err := db.conn.Get(&user, `
SELECT * FROM users WHERE id in (
	SELECT owner_id FROM chats WHERE id = ?
)`, chatId)
	if err != nil {
		return nil, fmt.Errorf("error getting chat[%d] owner: %s", chatId, err.Error())
	}
	return &user, nil
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
