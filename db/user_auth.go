package db

import (
	"fmt"

	a "go.chat/model/app"
)

const SchemaAuth string = `
	CREATE TABLE IF NOT EXISTS auth (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		user_id INTEGER,
		type TEXT,
		hash TEXT
	);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_auth_user_id_type ON auth(user_id, type, hash);`

func (db *DBConn) AddAuth(auth a.UserAuth) (*a.UserAuth, error) {
	if !db.IsActive() {
		return nil, fmt.Errorf("db is not connected")
	}

	db.mu.Lock()
	defer db.mu.Unlock()
	if auth.Id != 0 {
		return nil, fmt.Errorf("auth already has an id[%d]", auth.Id)
	} else if auth.UserId == 0 {
		return nil, fmt.Errorf("auth has no user id")
	} else if auth.Type == a.AuthTypeUnknown {
		return nil, fmt.Errorf("auth type is unknown")
	} else if auth.Hash == "" {
		return nil, fmt.Errorf("auth has no hash")
	}

	result, err := db.conn.Exec(`INSERT INTO auth (user_id, type, hash) VALUES (?, ?, ?)`,
		auth.UserId, auth.Type, auth.Hash)
	if err != nil {
		return nil, fmt.Errorf("error adding auth: %s", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %s", err)
	}

	auth.Id = uint(lastID)
	return &auth, nil
}

func (db *DBConn) GetAuth(userid uint, auth a.AuthType, hash string) (*a.UserAuth, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if !db.isConn {
		return nil, fmt.Errorf("db is not connected")
	}
	if auth == a.AuthTypeUnknown {
		return nil, fmt.Errorf("auth type is unknown")
	}

	var dbAuth a.UserAuth
	err := db.conn.Get(&dbAuth, `SELECT * FROM auth WHERE user_id = ? AND type = ? and hash = ?`, userid, auth, hash)
	if err != nil {
		return nil, fmt.Errorf("error getting auth: %s", err)
	}

	return &dbAuth, nil
}
