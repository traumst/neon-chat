package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
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
		FOREIGN KEY(owner_id) REFERENCES users(id) ON DELETE CASCADE
	);`

const ChatIndex = `CREATE INDEX IF NOT EXISTS idx_chat_title ON chats(title);`

func (dbConn *DBConn) ChatTableExists() bool {
	return dbConn.TableExists("chats")
}

func AddChat(dbConn sqlx.Ext, chat *Chat) (*Chat, error) {
	if chat.Id != 0 {
		return nil, fmt.Errorf("chat already has an id[%d]", chat.Id)
	} else if len(chat.Title) < minChatTitleLen || len(chat.Title) > maxChatTitleLen {
		return nil, fmt.Errorf("chat has no title")
	} else if chat.OwnerId == 0 {
		return nil, fmt.Errorf("chat has no owner")
	}

	result, err := dbConn.Exec(`INSERT INTO chats (title, owner_id) VALUES (?, ?)`, chat.Title, chat.OwnerId)
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

func GetChat(dbConn sqlx.Ext, chatId uint) (*Chat, error) {
	if chatId == 0 {
		return nil, fmt.Errorf("bad input: chatId[%d]", chatId)
	}

	var chat Chat
	err := sqlx.Get(dbConn, &chat, `SELECT * FROM chats WHERE id = ?`, chatId)
	if err != nil {
		return nil, fmt.Errorf("error getting chat: %s", err)
	}
	return &chat, nil
}

func GetOwner(dbConn sqlx.Ext, chatId uint) (*User, error) {
	if chatId == 0 {
		return nil, fmt.Errorf("bad input: chatId[%d]", chatId)
	}

	var user User
	err := sqlx.Get(dbConn, &user, `SELECT * FROM users WHERE id in (SELECT owner_id FROM chats WHERE id = ?)`, chatId)
	if err != nil {
		return nil, fmt.Errorf("error getting chat[%d] owner: %s", chatId, err.Error())
	}
	return &user, nil
}

func DeleteChat(dbConn sqlx.Ext, chatId uint) error {
	if chatId == 0 {
		return fmt.Errorf("bad input: chatId[%d]", chatId)
	}
	_, err := dbConn.Exec(`DELETE FROM chats WHERE id = ?`, chatId)
	if err != nil {
		return fmt.Errorf("error deleting chat[%d]: %s", chatId, err)
	}
	return nil
}
