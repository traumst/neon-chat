package db

import (
	"fmt"
)

type User struct {
	Id   uint   `db:"id"`
	Name string `db:"name"`
	Type string `db:"type"`
	Salt string `db:"salt"`
}

const SchemaUser string = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		name TEXT, 
		type TEXT,
		salt INTEGER
	);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_users_name ON users(name);`

func (db *DBConn) AddUser(user *User) (*User, error) {
	if user.Id != 0 {
		return nil, fmt.Errorf("user already has an id[%d]", user.Id)
	} else if user.Name == "" {
		return nil, fmt.Errorf("user has no name")
	} else if user.Type == "" {
		return nil, fmt.Errorf("user has no type")
	} else if len(user.Salt) <= 0 {
		return nil, fmt.Errorf("user has no salt")
	}
	if !db.IsActive() {
		return nil, fmt.Errorf("db is not connected")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	result, err := db.conn.Exec(`INSERT INTO users (name, type, salt) VALUES (?, ?, ?)`,
		user.Name, user.Type, user.Salt[:])
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

func (db *DBConn) UpdateName(user User) error {
	if user.Id <= 0 {
		return fmt.Errorf("invalid user id[%d]", user.Id)
	}
	if !db.IsActive() {
		return fmt.Errorf("db is not connected")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	result, err := db.conn.Exec(`UPDATE users SET name = ? WHERE id = ?;`,
		user.Name, user.Id)
	if err != nil {
		return fmt.Errorf("error adding user: %s", err)
	}
	count, err := result.RowsAffected()
	if err != nil || count != 1 {
		return fmt.Errorf("error estimating affected rows: %s", err)
	}
	return nil
}

func (db *DBConn) DropUser(userId uint) error {
	if !db.IsActive() {
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

func (db *DBConn) SearchUser(name string) (*User, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}
	if !db.isConn || !db.isInit {
		return nil, fmt.Errorf("db is not connected")
	}
	if name == "" {
		return nil, fmt.Errorf("name was not provided")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	var user User
	err := db.conn.Get(&user, `SELECT * FROM users WHERE name = ?`, name)
	if err != nil {
		return nil, fmt.Errorf("user[%s] not found: %s", name, err)
	}
	return &user, err
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
