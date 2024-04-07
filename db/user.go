package db

import (
	"fmt"

	"go.chat/model/app"
)

const SchemaUser string = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		name TEXT, 
		type TEXT,
		salt INTEGER
	);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_users_name ON users(name);`

func (db *DBConn) AddUser(user *app.User) (*app.User, error) {
	if !db.IsActive() {
		return nil, fmt.Errorf("db is not connected")
	}

	db.mu.Lock()
	defer db.mu.Unlock()
	if user.Id != 0 {
		return nil, fmt.Errorf("user already has an id[%d]", user.Id)
	} else if user.Name == "" {
		return nil, fmt.Errorf("user has no name")
	} else if user.Type == "" {
		return nil, fmt.Errorf("user has no type")
	} else if len(user.Salt) == 0 {
		return nil, fmt.Errorf("user has no salt")
	}

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

func (db *DBConn) GetUser(name string) (*app.User, error) {
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
	var user app.User
	err := db.conn.Get(&user, `SELECT * FROM users WHERE name = ?`, name)
	if err != nil {
		return nil, fmt.Errorf("user[%s] not found: %s", name, err)
	}
	return &user, err
}
