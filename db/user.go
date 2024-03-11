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

func (db *DBConn) AddUser(user app.User) (*app.User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if !db.IsActive() {
		return nil, fmt.Errorf("db is not connected")
	}
	if user.Id != 0 {
		return nil, fmt.Errorf("user already has an id[%d]", user.Id)
	} else if user.Name == "" {
		return nil, fmt.Errorf("user has no name")
	} else if user.Salt == nil || len(user.Salt) == 0 {
		return nil, fmt.Errorf("user has no salt")
	}

	result, err := db.conn.Exec(`INSERT INTO users (name, salt) VALUES (?, ?)`, user.Name, user.Salt)
	if err != nil {
		return nil, fmt.Errorf("error adding user: %s", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %s", err)
	}

	user.Id = uint(lastID)
	return &user, nil
}

func (db *DBConn) GetUser(name string) (*app.User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if !db.isConn || !db.isInit {
		return nil, fmt.Errorf("db is not connected")
	}
	if name == "" {
		return nil, fmt.Errorf("no id or name provided")
	}

	var user app.User
	err := db.conn.Get(&user, `SELECT * FROM users WHERE name = ?`, name)
	if err != nil {
		return nil, fmt.Errorf("error getting user by name[%s]: %s", name, err)
	}

	return &user, err
}
