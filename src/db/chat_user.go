package db

import (
	"fmt"
	"log"

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

func (db *DBConn) UsersCanChat(chatId uint, userIds ...uint) (bool, error) {
	if !db.ConnIsActive() {
		return false, fmt.Errorf("db is not connected")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	query, args, err := sqlx.In(`SELECT * FROM chat_users WHERE chat_id = ? AND user_id IN (?)`, chatId, userIds)
	if err != nil {
		return false, fmt.Errorf("error preparing query, %s", err)
	}
	query = db.conn.Rebind(query)

	var chatUser []ChatUser
	err = db.conn.Select(&chatUser, query, args...)
	if err != nil {
		return false, fmt.Errorf("error getting chats for userIds[%v]: %s", userIds, err)
	}
	misses := make([]uint, 0)
	for _, uid := range userIds {
		for _, cu := range chatUser {
			if cu.UserId == uid {
				break
			}

			misses = append(misses, uid)
		}
	}
	if len(misses) > 0 {
		return false, fmt.Errorf("some users are not in chat, userIds[%v], misses[%v]", userIds, misses)
	}
	return true, nil
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

	var chatIds []uint
	err := db.conn.Select(&chatIds, `SELECT chat_id FROM chat_users WHERE user_id = ?`, userId)
	if err != nil {
		return nil, fmt.Errorf("error adding user: %s", err.Error())
	}
	return chatIds, nil
}

func (db *DBConn) GetUserChats(userId uint) ([]Chat, error) {
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
	if err != nil {
		return nil, fmt.Errorf("error getting chats for user[%d]: %s", userId, err)
	}
	return userChats, nil
}

func (db *DBConn) GetChatUserIds(chatId uint) ([]uint, error) {
	if chatId == 0 {
		return nil, fmt.Errorf("bad input: chatId[%d]", chatId)
	}
	if !db.ConnIsActive() {
		return nil, fmt.Errorf("db is not connected")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	var userIds []uint
	err := db.conn.Select(&userIds, `SELECT user_id FROM chat_users WHERE chat_id = ?`, chatId)
	if err != nil {
		return nil, fmt.Errorf("error getting chat users: %s", err.Error())
	}
	return userIds, nil
}

func (db *DBConn) GetChatUsers(chatId uint) ([]User, error) {
	userIds, err := db.GetChatUserIds(chatId)
	if err != nil {
		return nil, fmt.Errorf("error getting user ids in chat[%d]: %s", chatId, err)
	}
	if len(userIds) == 0 {
		return []User{}, nil
	}

	return db.GetUsers(userIds)
}

func (db *DBConn) RemoveChatUser(chatId uint, userId uint) error {
	log.Printf("TRACE removing user[%d] from chat[%d]\n", userId, chatId)
	if chatId == 0 || userId == 0 {
		return fmt.Errorf("bad input: chatId[%d], userId[%d]", chatId, userId)
	}
	if !db.ConnIsActive() {
		return fmt.Errorf("db is not connected")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	res, err := db.conn.Exec(`DELETE FROM chat_users WHERE chat_id = ? AND user_id = ?`, chatId, userId)
	if err != nil {
		return fmt.Errorf("error removing user: %s", err.Error())
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error rows affected, %s", err.Error())
	}
	if count != 1 {
		return fmt.Errorf("error rows affected, expected to affect 1 row but affected %d", count)
	}
	return nil
}
