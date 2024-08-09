package db

import "fmt"

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
		FOREIGN KEY(chat_id) REFERENCES chats(id)
		FOREIGN KEY(author_id) REFERENCES users(id)
	);`

const MessageIndex = `CREATE INDEX IF NOT EXISTS idx_message_chat_id ON messages(chat_id);
CREATE INDEX IF NOT EXISTS idx_message_author_id ON messages(author_id);`

func (db *DBConn) MessageTableExists() bool {
	return db.TableExists("messages")
}

func (db *DBConn) AddMessage(msg *Message) (*Message, error) {
	if msg.Id != 0 {
		return nil, fmt.Errorf("message already has an id[%d]", msg.Id)
	} else if msg.ChatId == 0 {
		return nil, fmt.Errorf("message has no chat")
	} else if msg.AuthorId == 0 {
		return nil, fmt.Errorf("message has no author")
	}
	if !db.ConnIsActive() {
		return nil, fmt.Errorf("db is not connected")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	result, err := db.conn.Exec(`INSERT INTO messages (chat_id, author_id, text) VALUES (?, ?, ?)`,
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

func (db *DBConn) GetMessage(msgId uint) (*Message, error) {
	if msgId == 0 {
		return nil, fmt.Errorf("bad input: msgId[%d]", msgId)
	}
	if !db.ConnIsActive() {
		return nil, fmt.Errorf("db is not connected")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	var message Message
	err := db.conn.Get(&message, `SELECT * FROM messages where id = ?`, msgId)
	if err != nil {
		return nil, fmt.Errorf("error getting message: %s", err)
	}
	return &message, nil
}

func (db *DBConn) DeleteMessage(msgId uint) error {
	if msgId == 0 {
		return fmt.Errorf("cannot delete message with id [%d]", msgId)
	}
	if !db.ConnIsActive() {
		return fmt.Errorf("db is not connected")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	_, err := db.conn.Exec(`DELETE FROM messages WHERE id = ?`, msgId)
	if err != nil {
		return fmt.Errorf("failed to delete message: %s", err.Error())
	}
	return nil
}