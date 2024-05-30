package db

import (
	"database/sql"
	"fmt"

	"prplchat/src/utils"
)

type Avatar struct {
	Id     int    `db:"id"`
	UserId uint   `db:"user_id"`
	Title  string `db:"title"`
	Size   int    `db:"size"`
	Image  []byte `db:"image"`
	Mime   string `db:"mime"`
}

const AvatarSchema = `
	CREATE TABLE IF NOT EXISTS avatars (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		user_id INTEGER,
		title TEXT,
		size INTEGER,
		image BLOB,
		mime TEXT,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`
const AvatarIndex = `CREATE INDEX IF NOT EXISTS idx_avatar_user_id ON avatars(user_id);`

func (db *DBConn) AvatarTableExists() bool {
	return db.TableExists("avatars")
}

func (db *DBConn) AddAvatar(userId uint, title string, image []byte, mime string) (*Avatar, error) {
	if userId <= 0 {
		return nil, fmt.Errorf("avatar must have user id")
	}
	title = utils.TrimSpaces(title)
	title = utils.TrimSpecial(title)
	size := len(image)
	if size <= 0 {
		return nil, fmt.Errorf("avatar requires an image")
	}
	limit := 10 * utils.KB
	if size > limit {
		return nil, fmt.Errorf("avatar image size[%d] is over limit[%d]", size, limit)
	}
	if !db.ConnIsActive() {
		return nil, fmt.Errorf("db is not connected")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	result, err := db.conn.Exec(`INSERT INTO avatars (user_id, title, size, image, mime) VALUES (?, ?, ?, ?, ?)`,
		userId, title, size, image, mime)
	if err != nil {
		return nil, fmt.Errorf("error adding avatar: %s", err)
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert avatarId: %s", err)
	}
	avatar := Avatar{
		Id:     int(lastId),
		UserId: userId,
		Title:  title,
		Size:   size,
		Image:  image,
		Mime:   mime,
	}
	return &avatar, nil
}

func (db *DBConn) GetAvatar(userId uint) (*Avatar, error) {
	if !db.isConn {
		return nil, fmt.Errorf("db is not connected")
	}
	if userId <= 0 {
		return nil, fmt.Errorf("invalid userId[%d]", userId)
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	var avatar Avatar
	err := db.conn.Get(&avatar, `SELECT * FROM avatars WHERE user_id = ?`, userId)
	if err != nil {
		return nil, fmt.Errorf("error getting avatar for user[%d]: %s", userId, err)
	}
	return &avatar, nil
}

func (db *DBConn) GetAvatars(userId uint) ([]*Avatar, error) {
	if !db.isConn {
		return nil, fmt.Errorf("db is not connected")
	}
	if userId <= 0 {
		return nil, fmt.Errorf("invalid userId[%d]", userId)
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	avatars := make([]*Avatar, 0)
	err := db.conn.Select(&avatars, `SELECT * FROM avatars WHERE user_id = ?`, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return avatars, nil
		} else {
			return nil, fmt.Errorf("error getting avatar for user[%d]: %s", userId, err)
		}
	}
	return avatars, nil
}

func (db *DBConn) DropAvatar(id int) error {
	if !db.ConnIsActive() {
		return fmt.Errorf("db is not connected")
	}
	if id <= 0 {
		return fmt.Errorf("invalid avatar id[%d]", id)
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	_, err := db.conn.Exec(`DELETE FROM avatars where id = ?`, id)
	if err != nil {
		return fmt.Errorf("error deleting avatar: %s", err.Error())
	}
	return nil
}
