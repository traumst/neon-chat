package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Auth struct {
	Id     uint   `db:"id"`
	UserId uint   `db:"user_id"`
	Type   string `db:"type"`
	Hash   string `db:"hash"`
}

const AuthSchema string = `
	CREATE TABLE IF NOT EXISTS auth (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		user_id INTEGER,
		type TEXT,
		hash TEXT,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	);`
const AuthIndex string = `CREATE UNIQUE INDEX IF NOT EXISTS idx_auth_type_hash ON auth(type, hash);`

func (dbConn *DBConn) AuthTableExists() bool {
	return dbConn.TableExists("auth")
}

func AddAuth(dbConn sqlx.Ext, auth Auth) (*Auth, error) {
	if auth.Id != 0 {
		return nil, fmt.Errorf("auth already has an id[%d]", auth.Id)
	} else if auth.UserId <= 0 {
		return nil, fmt.Errorf("auth has no user id")
	} else if auth.Type == "" {
		return nil, fmt.Errorf("auth type is unknown")
	} else if auth.Hash == "" {
		return nil, fmt.Errorf("auth has no hash")
	}
	result, err := dbConn.Exec(`INSERT INTO auth (user_id, type, hash) VALUES (?, ?, ?)`,
		auth.UserId,
		auth.Type,
		auth.Hash)
	if err != nil {
		return nil, fmt.Errorf("error adding auth: %s", err.Error())
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %s", err.Error())
	}
	auth.Id = uint(lastId)
	return &auth, nil
}

func GetUserAuth(dbConn sqlx.Ext, userId uint) (*Auth, error) {
	var dbAuth Auth
	err := sqlx.Get(dbConn, &dbAuth, `SELECT * FROM auth WHERE user_id = ?`, userId)
	if err != nil {
		return nil, fmt.Errorf("error getting auth for user[%d]: %s", userId, err.Error())
	}
	return &dbAuth, nil
}

func GetAuth(dbConn sqlx.Ext, auth string, hash string) (*Auth, error) {
	var dbAuth Auth
	err := sqlx.Get(dbConn, &dbAuth, `SELECT * FROM auth WHERE type = ? AND hash = ?`, auth, hash)
	if err != nil {
		return nil, fmt.Errorf("error getting auth by type_hash[%s_%s]: %s", auth, hash, err.Error())
	}
	return &dbAuth, nil
}
