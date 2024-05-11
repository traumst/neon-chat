package db

import (
	"fmt"

	"go.chat/src/utils"
)

type UserAvatar struct {
	Id     int    `db:"id"`
	UserId uint   `db:"user_id"`
	Title  string `db:"title"`
	Size   int    `db:"size"`
	Image  []byte `db:"image"`
	Mime   string `db:"mime"`
}

const SchemaUserAvatar string = `
	CREATE TABLE IF NOT EXISTS avatars (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		user_id INTEGER,
		title TEXT,
		size INTEGER,
		image BLOB,
		mime TEXT,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_users_name ON users(name);`

func (db *DBConn) AddAvatar(userId uint, title string, image []byte, mime string) (*UserAvatar, error) {
	if userId <= 0 {
		return nil, fmt.Errorf("avatar must have user id")
	}
	title = utils.TrimSpaces(title)
	title = utils.TrimSpecial(title)
	size := len(image)
	if size <= 0 {
		return nil, fmt.Errorf("avatar requires an image")
	}
	limit := size / (10 * utils.KB)
	if size > limit {
		return nil, fmt.Errorf("avatar image size[%d] is over limit[%d]", size, limit)
	}
	if !db.IsActive() {
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
	avatar := UserAvatar{
		Id:     int(lastId),
		UserId: userId,
		Title:  title,
		Size:   size,
		Image:  image,
		Mime:   mime,
	}
	return &avatar, nil
}

func (db *DBConn) GetAvatar(id int) (*UserAvatar, error) {
	if !db.isConn {
		return nil, fmt.Errorf("db is not connected")
	}
	if id <= 0 {
		return nil, fmt.Errorf("invalid avatarId[%d]", id)
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	var avatar UserAvatar
	err := db.conn.Get(&avatar, `SELECT * FROM avatars WHERE id = ?`, id)
	if err != nil {
		return nil, fmt.Errorf("error getting avatar[%d]: %s", id, err)
	}
	return &avatar, nil
}

func (db *DBConn) GetAvatars(userId uint) ([]*UserAvatar, error) {
	if !db.isConn {
		return nil, fmt.Errorf("db is not connected")
	}
	if userId <= 0 {
		return nil, fmt.Errorf("invalid userId[%d]", userId)
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	avatars := make([]*UserAvatar, 0)
	err := db.conn.Select(&avatars, `SELECT * FROM avatars WHERE user_id = ?`, userId)
	if err != nil {
		return nil, fmt.Errorf("error getting avatars for user[%d]: %s", userId, err)
	}
	return avatars, nil
}

func (db *DBConn) DropAvatar(id int) error {
	if !db.IsActive() {
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
