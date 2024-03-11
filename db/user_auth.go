package db

import (
	"fmt"

	"go.chat/model/app"
)

const SchemaAuth string = `
	CREATE TABLE IF NOT EXISTS auth (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		user_id INTEGER,
		type TEXT,
		hash INTEGER
	);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_auth_user_id_type ON auth(user_id, type, hash);`

func (db *DBConn) AddAuth(auth app.Auth) (*app.Auth, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if !db.IsActive() {
		return nil, fmt.Errorf("db is not connected")
	}
	if auth.Id != 0 {
		return nil, fmt.Errorf("auth already has an id[%d]", auth.Id)
	} else if auth.UserId == 0 {
		return nil, fmt.Errorf("auth has no user id")
	} else if auth.Type == app.AuthTypeUnknown {
		return nil, fmt.Errorf("auth type is unknown")
	} else if auth.Hash == 0 {
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

func (db *DBConn) GetAuth(userid int, auth app.AuthType, hash int) (*app.Auth, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if !db.isConn {
		return nil, fmt.Errorf("db is not connected")
	}
	if auth == app.AuthTypeUnknown {
		return nil, fmt.Errorf("auth type is unknown")
	}

	var dbAuth app.Auth
	err := db.conn.Get(&dbAuth, `SELECT * FROM auth WHERE user_id = ? AND type = ? and hash = ?`, userid, auth, hash)
	if err != nil {
		return nil, fmt.Errorf("error getting auth: %s", err)
	}

	return &dbAuth, nil
}
