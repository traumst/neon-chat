package db

import (
	"database/sql"
	"fmt"

	"prplchat/src/utils"

	"github.com/jmoiron/sqlx"
)

type Avatar struct {
	Id     uint   `db:"id"`
	UserId uint   `db:"user_id"`
	Title  string `db:"title"`
	Size   int    `db:"size"`
	Image  []byte `db:"image"`
	Mime   string `db:"mime"`
}

const AvatarMaxUploadBytesSize int64 = 50 * utils.KB

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
	title = utils.ReplaceWithSingleSpace(title)
	title = utils.RemoveSpecialChars(title)
	size := len(image)
	if size <= 0 {
		return nil, fmt.Errorf("avatar requires an image")
	}
	if int64(size) > AvatarMaxUploadBytesSize {
		return nil, fmt.Errorf("avatar image size[%d] is over limit[%d]", size, AvatarMaxUploadBytesSize)
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
		Id:     uint(lastId),
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

func (db *DBConn) GetAvatars(userIds []uint) ([]*Avatar, error) {
	if !db.isConn {
		return nil, fmt.Errorf("db is not connected")
	}
	if len(userIds) <= 0 {
		return nil, fmt.Errorf("empty input userIds")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	query, args, err := sqlx.In(`SELECT * FROM avatars WHERE user_id IN (?)`, userIds)
	if err != nil {
		return nil, fmt.Errorf("error preparing select avatars query for userIds %v, %s", userIds, err)
	}
	query = db.conn.Rebind(query)

	var avatars []Avatar
	err = db.conn.Select(&avatars, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error getting avatars for userIds %v: %s", userIds, err)
	}
	var avatarsPtrs []*Avatar
	for i := range avatars {
		avatarsPtrs = append(avatarsPtrs, &avatars[i])
	}
	return avatarsPtrs, nil
}

func (db *DBConn) GetUserAvatars(userId uint) ([]*Avatar, error) {
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

func (db *DBConn) DropAvatar(id uint) error {
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
