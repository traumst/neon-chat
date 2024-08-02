package db

import "fmt"

type Message struct {
	Id       uint   `db:"id"`
	ChatId   uint   `db:"chat_id"`
	OwnerId  uint   `db:"owner_id"`
	AuthorId uint   `db:"author_id"`
	Text     string `db:"text"`
}

const MessageSchema = `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		chat_id INTEGER, 
		owner_id INTEGER,
		author_id INTEGER,
		text TEXT,
		FOREIGN KEY(chat_id) REFERENCES chats(id)
		FOREIGN KEY(owner_id) REFERENCES users(id)
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
	} else if msg.OwnerId == 0 {
		return nil, fmt.Errorf("message has no owner")
	} else if msg.AuthorId == 0 {
		return nil, fmt.Errorf("message has no author")
	}
	if !db.ConnIsActive() {
		return nil, fmt.Errorf("db is not connected")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	result, err := db.conn.Exec(`INSERT INTO messages (chat_id, owner_id, author_id, text) VALUES (?, ?, ?, ?)`,
		msg.ChatId, msg.OwnerId, msg.AuthorId, msg.Text)
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
