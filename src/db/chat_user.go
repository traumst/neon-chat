package db

import (
	"fmt"
	"log"

	"neon-chat/src/consts"

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
		PRIMARY KEY (chat_id, user_id),
		FOREIGN KEY(chat_id) REFERENCES chats(id) ON DELETE CASCADE,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	);`
const ChatUserIndex = ``

func (dbConn *DBConn) ChatUserTableExists() bool {
	return dbConn.TableExists("chat_users")
}

func UsersCanChat(dbConn sqlx.Ext, chatId uint, userIds ...uint) (bool, error) {
	if chatId == 0 || len(userIds) == 0 {
		return false, fmt.Errorf("bad input: chatId[%d], userIds[%v]", chatId, userIds)
	}
	queryAcc := `SELECT * FROM chat_users WHERE chat_id = ?`
	uids := []interface{}{chatId}
	for _, userId := range userIds {
		queryAcc += ` AND user_id = ?`
		uids = append(uids, userId)
	}
	query, args, err := sqlx.In(queryAcc, uids...)
	if err != nil {
		return false, fmt.Errorf("error preparing query, %s", err)
	} else if len(args) != len(userIds)+1 {
		return false, fmt.Errorf("missing args")
	} else if query == "" {
		return false, fmt.Errorf("empty query")
	}
	if dbConn == nil {
		return false, fmt.Errorf("db conn is nil")
	}
	query = dbConn.Rebind(query)
	var chatUser []ChatUser
	err = sqlx.Select(dbConn, &chatUser, query, args...)
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

func AddChatUser(dbConn sqlx.Ext, chatId uint, userId uint) error {
	if chatId == 0 || userId == 0 {
		return fmt.Errorf("bad input: chatId[%d], userId[%d]", chatId, userId)
	}

	_, err := dbConn.Exec(`INSERT INTO chat_users (chat_id, user_id) VALUES (?, ?)`, chatId, userId)
	if err != nil {
		return fmt.Errorf("error adding user: %s", err.Error())
	}
	return nil
}

func GetUserChatIds(dbConn sqlx.Ext, userId uint) ([]uint, error) {
	if userId == 0 {
		return nil, fmt.Errorf("bad input: userId[%d]", userId)
	}

	var chatIds []uint
	err := sqlx.Select(dbConn, &chatIds, `SELECT chat_id FROM chat_users WHERE user_id = ?`, userId)
	if err != nil {
		return nil, fmt.Errorf("error adding user: %s", err.Error())
	}
	return chatIds, nil
}

func GetUserChats(dbConn sqlx.Ext, userId uint) ([]Chat, error) {
	chatIds, err := GetUserChatIds(dbConn, userId)
	if err != nil {
		return nil, fmt.Errorf("error getting chat ids for user[%d]: %s", userId, err)
	}
	if len(chatIds) == 0 {
		return []Chat{}, nil
	}

	query, args, err := sqlx.In(`SELECT * FROM chats WHERE id IN (?)`, chatIds)
	if err != nil {
		return nil, fmt.Errorf("error preparing query for user[%d]: %s", userId, err)
	}
	query = dbConn.Rebind(query)

	var userChats []Chat
	err = sqlx.Select(dbConn, &userChats, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error getting chats for user[%d]: %s", userId, err)
	}
	return userChats, nil
}

// INNER-JOIN
func GetSharedChatIds(dbConn sqlx.Ext, userIds []uint) ([]uint, error) {
	if len(userIds) != 2 {
		return nil, fmt.Errorf("expected exactly 2 userIds, but got userIds[%v]", userIds)
	}
	withLimit := fmt.Sprintf(`
        SELECT L.chat_id
        FROM chat_users L
        JOIN chat_users R
            ON L.chat_id = R.chat_id
        WHERE L.user_id = ? 
            AND R.user_id = ?
        ORDER BY L.chat_id
        LIMIT %d;`, consts.MaxSharedChats)
	log.Printf("TRACE shared chat ids query for users[%v]: %s\n", userIds, withLimit)
	query, args, err := sqlx.In(withLimit, userIds[0], userIds[1])
	if err != nil {
		return nil, fmt.Errorf("error preparing shared chats query: %s", err)
	}
	query = dbConn.Rebind(query)
	var chatIds []uint
	err = sqlx.Select(dbConn, &chatIds, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error getting shared chat ids: %s", err.Error())
	}
	log.Printf("TRACE shared chat ids between users[%v]: %v\n", userIds, chatIds)
	return chatIds, nil
}

// INNER-JOIN
func GetSharedChats(dbConn sqlx.Ext, userIds []uint) ([]Chat, error) {
	chatIds, err := GetSharedChatIds(dbConn, userIds)
	if err != nil {
		return nil, fmt.Errorf("error getting shared chat ids for users[%v]: %s", userIds, err)
	}
	if len(chatIds) == 0 {
		log.Printf("INFO no shared chats between users[%v]\n", userIds)
		return []Chat{}, nil
	}
	log.Printf("TRACE shared chats between users[%v] ids: %d\n", userIds, len(chatIds))
	query, args, err := sqlx.In(`SELECT * FROM chats WHERE id IN (?)`, chatIds)
	if err != nil {
		return nil, fmt.Errorf("error preparing chatIds query for [%v]: %s", chatIds, err)
	}
	query = dbConn.Rebind(query)
	var sharedChats []Chat
	err = sqlx.Select(dbConn, &sharedChats, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error getting shared chats[%v]: %s", chatIds, err)
	}
	log.Printf("TRACE shared chats between users[%v] found: %d\n", userIds, len(sharedChats))
	return sharedChats, nil
}

func GetChatUserIds(dbConn sqlx.Ext, chatId uint) ([]uint, error) {
	if chatId == 0 {
		return nil, fmt.Errorf("bad input: chatId[%d]", chatId)
	}
	var userIds []uint
	err := sqlx.Select(dbConn, &userIds, `SELECT user_id FROM chat_users WHERE chat_id = ?`, chatId)
	if err != nil {
		return nil, fmt.Errorf("error getting chat users: %s", err.Error())
	}
	return userIds, nil
}

func RemoveChatUser(dbConn sqlx.Ext, chatId uint, userId uint) error {
	log.Printf("TRACE removing user[%d] from chat[%d]\n", userId, chatId)
	if chatId == 0 || userId == 0 {
		return fmt.Errorf("bad input: chatId[%d], userId[%d]", chatId, userId)
	}

	res, err := dbConn.Exec(`DELETE FROM chat_users WHERE chat_id = ? AND user_id = ?`, chatId, userId)
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
