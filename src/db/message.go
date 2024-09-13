package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Message struct {
	Id       uint   `db:"id"`
	ChatId   uint   `db:"chat_id"`
	AuthorId uint   `db:"author_id"`
	Text     string `db:"text"`
}

const MessageSchema = `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		chat_id INTEGER, 
		author_id INTEGER,
		text TEXT,
		FOREIGN KEY(chat_id) REFERENCES chats(id),
		FOREIGN KEY(author_id) REFERENCES users(id)
	);`

const MessageIndex = `CREATE INDEX IF NOT EXISTS idx_message_text ON messages(text);`

func (db *DBConn) MessageTableExists() bool {
	return db.TableExists("messages")
}

func AddMessage(dbConn sqlx.Ext, msg *Message) (*Message, error) {
	if msg.Id != 0 {
		return nil, fmt.Errorf("message already has an id[%d]", msg.Id)
	} else if msg.ChatId == 0 {
		return nil, fmt.Errorf("message has no chat")
	} else if msg.AuthorId == 0 {
		return nil, fmt.Errorf("message has no author")
	}
	result, err := dbConn.Exec(`INSERT INTO messages (chat_id, author_id, text) VALUES (?, ?, ?)`,
		msg.ChatId, msg.AuthorId, msg.Text)
	if err != nil {
		return nil, fmt.Errorf("error adding message: %s", err)
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %s", err)
	}
	msg.Id = uint(lastId)

	return msg, nil
}

func GetMessage(dbConn sqlx.Ext, msgId uint) (*Message, error) {
	if msgId == 0 {
		return nil, fmt.Errorf("bad input: msgId[%d]", msgId)
	}
	var message Message
	err := sqlx.Get(dbConn, &message, `SELECT * FROM messages where id = ?`, msgId)
	if err != nil {
		return nil, fmt.Errorf("error getting message: %s", err)
	}
	return &message, nil
}

func GetMessages(dbConn sqlx.Ext, chatId uint, offset int) ([]Message, error) {
	if chatId == 0 {
		return nil, fmt.Errorf("bad input: chatId[%d]", chatId)
	}
	var messages []Message
	err := sqlx.Select(dbConn, &messages, `SELECT * FROM messages where chat_id = ?`, chatId)
	if err != nil {
		return nil, fmt.Errorf("error getting chat[%d] messages: %s", chatId, err)
	}
	return messages, nil
}

func DeleteMessage(dbConn sqlx.Ext, msgId uint) error {
	if msgId == 0 {
		return fmt.Errorf("cannot delete message with id [%d]", msgId)
	}
	_, err := dbConn.Exec(`DELETE FROM messages WHERE id = ?`, msgId)
	if err != nil {
		return fmt.Errorf("failed to delete message: %s", err.Error())
	}
	return nil
}
