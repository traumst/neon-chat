package db

import (
	"database/sql"
	"fmt"

	"neon-chat/src/utils"

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
const AvatarIndex = ``

func (db *DBConn) AvatarTableExists() bool {
	return db.TableExists("avatars")
}

func AddAvatar(dbConn sqlx.Ext, userId uint, title string, image []byte, mime string) (*Avatar, error) {
	if userId <= 0 {
		return nil, fmt.Errorf("avatar must have user id")
	}
	title = utils.ReplaceWithSingleSpace(title)
	title = utils.RemoveSpecialChars(title)
	size := len(image)
	if size <= 0 {
		return nil, fmt.Errorf("avatar requires an image")
	}
	if int64(size) > utils.MaxUploadBytesSize {
		return nil, fmt.Errorf("avatar image size[%d] is over limit[%d]", size, utils.MaxUploadBytesSize)
	}

	result, err := dbConn.Exec(`INSERT INTO avatars (user_id, title, size, image, mime) VALUES (?, ?, ?, ?, ?)`,
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

func GetAvatar(dbConn sqlx.Ext, userId uint) (*Avatar, error) {
	if userId <= 0 {
		return nil, fmt.Errorf("invalid userId[%d]", userId)
	}

	var avatar Avatar
	err := sqlx.Get(dbConn, &avatar, `SELECT * FROM avatars WHERE user_id = ?`, userId)
	if err != nil {
		return nil, fmt.Errorf("error getting avatar for user[%d]: %s", userId, err)
	}
	return &avatar, nil
}

func GetAvatars(dbConn sqlx.Ext, userIds []uint) ([]*Avatar, error) {
	if len(userIds) <= 0 {
		return nil, fmt.Errorf("empty input userIds")
	}

	query, args, err := sqlx.In(`SELECT * FROM avatars WHERE user_id IN (?)`, userIds)
	if err != nil {
		return nil, fmt.Errorf("error preparing select avatars query for userIds %v, %s", userIds, err)
	}
	query = dbConn.Rebind(query)

	var avatars []Avatar
	err = sqlx.Select(dbConn, &avatars, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error getting avatars for userIds %v: %s", userIds, err)
	}
	var avatarsPtrs []*Avatar
	for i := range avatars {
		avatarsPtrs = append(avatarsPtrs, &avatars[i])
	}
	return avatarsPtrs, nil
}

func GetUserAvatars(dbConn sqlx.Ext, userId uint) ([]*Avatar, error) {
	if userId <= 0 {
		return nil, fmt.Errorf("invalid userId[%d]", userId)
	}

	avatars := make([]*Avatar, 0)
	err := sqlx.Select(dbConn, &avatars, `SELECT * FROM avatars WHERE user_id = ?`, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return avatars, nil
		} else {
			return nil, fmt.Errorf("error getting avatar for user[%d]: %s", userId, err)
		}
	}
	return avatars, nil
}

func DropAvatar(dbConn sqlx.Ext, id uint) error {
	if id <= 0 {
		return fmt.Errorf("invalid avatar id[%d]", id)
	}

	_, err := dbConn.Exec(`DELETE FROM avatars where id = ?`, id)
	if err != nil {
		return fmt.Errorf("error deleting avatar: %s", err.Error())
	}
	return nil
}
