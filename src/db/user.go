package db

import (
	"fmt"
)

const minUserNameLen = 3
const maxUserNameLen = 64

type User struct {
	Id     uint   `db:"id"`
	Name   string `db:"name"`
	Email  string `db:"email"`
	Type   string `db:"type"`
	Status string `db:"status"`
	Salt   string `db:"salt"`
}

const UserSchema = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		name TEXT UNIQUE, 
		email TEXT UNIQUE,
		type TEXT,
		status TEXT,
		salt INTEGER
	);`
const UserIndex = `CREATE INDEX IF NOT EXISTS idx_user_status ON users(status);`

func (db *DBConn) UserTableExists() bool {
	return db.TableExists("users")
}

func (db *DBConn) AddUser(user *User) (*User, error) {
	if user.Id != 0 {
		return nil, fmt.Errorf("user already has an id[%d]", user.Id)
	} else if len(user.Name) < minUserNameLen || len(user.Name) > maxUserNameLen {
		return nil, fmt.Errorf("user has no name")
	} else if len(user.Email) < minUserNameLen || len(user.Email) > maxUserNameLen {
		return nil, fmt.Errorf("user has no email")
	} else if len(user.Type) < minUserNameLen || len(user.Type) > maxUserNameLen {
		return nil, fmt.Errorf("user has no type")
	} else if len(user.Status) < minUserNameLen || len(user.Status) > maxUserNameLen {
		return nil, fmt.Errorf("user has no status")
	} else if len(user.Salt) == 0 {
		return nil, fmt.Errorf("user has no salt")
	}
	if !db.ConnIsActive() {
		return nil, fmt.Errorf("db is not connected")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	result, err := db.conn.Exec(`INSERT INTO users (name, email, type, status, salt) VALUES (?, ?, ?, ?, ?)`,
		user.Name, user.Email, user.Type, user.Status, user.Salt[:])
	if err != nil {
		return nil, fmt.Errorf("error adding user: %s", err)
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %s", err)
	}
	user.Id = uint(lastId)
	return user, nil
}

func (db *DBConn) DropUser(userId uint) error {
	if !db.ConnIsActive() {
		return fmt.Errorf("db is not connected")
	}
	if userId <= 0 {
		return fmt.Errorf("user does not have an id[%d]", userId)
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	_, err := db.conn.Exec(`DELETE FROM users where id = ?`, userId)
	if err != nil {
		return fmt.Errorf("error deleting user: %s", err.Error())
	}
	return nil
}

func (db *DBConn) SearchUser(login string) (*User, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}
	if !db.isConn || !db.isInit {
		return nil, fmt.Errorf("db is not connected")
	}
	if len(login) < minUserNameLen || len(login) > maxUserNameLen {
		return nil, fmt.Errorf("login name/email was not provided")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	var user User
	err := db.conn.Get(&user, `SELECT * FROM users WHERE name = ? or email = ?`, login, login)
	if err != nil {
		return nil, fmt.Errorf("login[%s] not found: %s", login, err)
	}
	return &user, err
}

func (db *DBConn) SearchUsers(name string) ([]*User, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}
	if !db.isConn || !db.isInit {
		return nil, fmt.Errorf("db is not connected")
	}
	if len(name) < minUserNameLen || len(name) > maxUserNameLen {
		return nil, fmt.Errorf("name was not provided")
	}
	users := make([]*User, 0)
	approxName := fmt.Sprintf("%%%s%%", name)

	db.mu.Lock()
	defer db.mu.Unlock()

	err := db.conn.Select(&users, `SELECT * FROM users WHERE name like ? or email like ?`,
		approxName, approxName)
	if err != nil {
		return nil, fmt.Errorf("user[%s] not found: %s", name, err)
	}
	return users, err
}

func (db *DBConn) GetUser(id uint) (*User, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}
	if !db.isConn || !db.isInit {
		return nil, fmt.Errorf("db is not connected")
	}
	if id <= 0 {
		return nil, fmt.Errorf("id was 0")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	var user User
	err := db.conn.Get(&user, `SELECT * FROM users WHERE id = ?`, id)
	if err != nil {
		return nil, fmt.Errorf("userId[%d] not found: %s", id, err)
	}
	return &user, err
}

func (db *DBConn) UpdateUserName(userId uint, userName string) error {
	if userId <= 0 {
		return fmt.Errorf("invalid user id[%d]", userId)
	}
	if !db.ConnIsActive() {
		return fmt.Errorf("db is not connected")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	result, err := db.conn.Exec(`UPDATE users SET name = ? WHERE id = ?;`, userName, userId)
	if err != nil {
		return fmt.Errorf("error updating user name: %s", err)
	}
	count, err := result.RowsAffected()
	if err != nil || count != 1 {
		return fmt.Errorf("error estimating affected rows: %s", err)
	}
	return nil
}

func (db *DBConn) UpdateUserStatus(userId uint, userStatus string) error {
	if userId <= 0 {
		return fmt.Errorf("invalid user id[%d]", userId)
	}
	if !db.ConnIsActive() {
		return fmt.Errorf("db is not connected")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	result, err := db.conn.Exec(`UPDATE users SET status = ? WHERE id = ?;`, userStatus, userId)
	if err != nil {
		return fmt.Errorf("error updating user status: %s", err)
	}
	count, err := result.RowsAffected()
	if err != nil || count != 1 {
		return fmt.Errorf("error estimating affected rows: %s", err)
	}
	return nil
}
