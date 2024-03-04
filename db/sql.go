package db

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jmoiron/sqlx"

	"go.chat/model"
)

type DBConn struct {
	mu     sync.Mutex
	conn   *sqlx.DB
	isConn bool
	isInit bool
}

func ConnectDB(dbPath string) DBConn {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
	}

	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	return DBConn{conn: db, isConn: true, isInit: false}
}

func (db *DBConn) IsActive() bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	return db.isConn && db.isInit
}

func (db *DBConn) InitDB() {
	db.mu.Lock()
	defer db.mu.Unlock()

	if !db.isConn {
		log.Fatal("DBConn is not connected")
	}

	if db.isInit {
		log.Fatal("DBConn is already initialized")
	}

	_, err := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			name TEXT, 
			salt INTEGER
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_users_name ON users(name);

		CREATE TABLE IF NOT EXISTS auth (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			user_id INTEGER,
			type TEXT,
			hash INTEGER
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_auth_user_id_type ON auth(user_id, type, hash);
	`)
	if err != nil {
		log.Fatalf("Error initializing database: %s", err)
	}

	db.isInit = true
}

func (db *DBConn) Close() {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.conn.Close() // TODO can panic I bet
}

func (db *DBConn) AddUser(user model.User) (*model.User, error) {
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

func (db *DBConn) GetUser(name string) (*model.User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if !db.isConn || !db.isInit {
		return nil, fmt.Errorf("db is not connected")
	}
	if name == "" {
		return nil, fmt.Errorf("no id or name provided")
	}

	var user model.User
	err := db.conn.Get(&user, `SELECT * FROM users WHERE name = ?`, name)
	if err != nil {
		return nil, fmt.Errorf("error getting user by name[%s]: %s", name, err)
	}

	return &user, err
}

func (db *DBConn) AddAuth(auth model.Auth) (*model.Auth, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if !db.isConn {
		return nil, fmt.Errorf("db is not connected")
	}
	if auth.Id != 0 {
		return nil, fmt.Errorf("auth already has an id[%d]", auth.Id)
	} else if auth.UserId == 0 {
		return nil, fmt.Errorf("auth has no user id")
	} else if auth.Type == model.AuthTypeUnknown {
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

func (db *DBConn) GetAuth(userid int, auth model.AuthType, hash int) (*model.Auth, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if !db.isConn {
		return nil, fmt.Errorf("db is not connected")
	}
	if auth == model.AuthTypeUnknown {
		return nil, fmt.Errorf("auth type is unknown")
	}

	var dbAuth model.Auth
	err := db.conn.Get(&dbAuth, `SELECT * FROM auth WHERE user_id = ? AND type = ? and hash = ?`, userid, auth, hash)
	if err != nil {
		return nil, fmt.Errorf("error getting auth: %s", err)
	}

	return &dbAuth, nil
}
