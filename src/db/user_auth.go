package db

import (
	"fmt"
)

type UserAuth struct {
	Id     uint   `db:"id"`
	UserId uint   `db:"user_id"`
	Type   string `db:"type"`
	Hash   string `db:"hash"`
}

const SchemaAuth string = `
	CREATE TABLE IF NOT EXISTS auth (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		user_id INTEGER,
		type TEXT,
		hash TEXT,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_auth_user_id_type ON auth(user_id, type, hash);`

func (db *DBConn) AddAuth(auth UserAuth) (*UserAuth, error) {
	if !db.IsActive() {
		return nil, fmt.Errorf("db is not connected")
	}

	db.mu.Lock()
	defer db.mu.Unlock()
	if auth.Id != 0 {
		return nil, fmt.Errorf("auth already has an id[%d]", auth.Id)
	} else if auth.UserId <= 0 {
		return nil, fmt.Errorf("auth has no user id")
	} else if auth.Type == "" {
		return nil, fmt.Errorf("auth type is unknown")
	} else if auth.Hash == "" {
		return nil, fmt.Errorf("auth has no hash")
	}

	result, err := db.conn.Exec(`INSERT INTO auth (user_id, type, hash) VALUES (?, ?, ?)`,
		auth.UserId, auth.Type, auth.Hash)
	if err != nil {
		return nil, fmt.Errorf("error adding auth: %s", err)
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %s", err)
	}

	auth.Id = uint(lastId)
	return &auth, nil
}

func (db *DBConn) GetAuth(userid uint, auth string, hash string) (*UserAuth, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if !db.isConn {
		return nil, fmt.Errorf("db is not connected")
	}

	var dbAuth UserAuth
	err := db.conn.Get(&dbAuth, `SELECT * FROM auth WHERE user_id = ? AND type = ? AND hash = ?`, userid, auth, hash)
	if err != nil {
		return nil, fmt.Errorf("error getting auth: %s", err)
	}

	return &dbAuth, nil
}
