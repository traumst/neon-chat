package db

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DBConn struct {
	mu     sync.Mutex
	conn   *sqlx.DB
	isConn bool
	isInit bool
}

func ConnectDB(dbPath string) (*DBConn, error) {
	if fi, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			return nil, fmt.Errorf("error creating db file: %s", err)
		}
		file.Close()
	} else {
		log.Printf("  opening db file [%s] [%d]", fi.Name(), fi.Size())
	}

	conn, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %s", err)
	}

	db := DBConn{conn: conn, isConn: true, isInit: false}
	err = db.init()
	return &db, err
}

func (db *DBConn) Close() {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.conn.Close()
}

func (db *DBConn) IsActive() bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	return db.isConn && db.isInit
}

func (db *DBConn) init() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if !db.isConn {
		return fmt.Errorf("DBConn is not connected")
	}

	if db.isInit {
		return fmt.Errorf("DBConn is already initialized")
	}

	schema := fmt.Sprintf("%s\n%s", SchemaUser, SchemaAuth)
	_, err := db.conn.Exec(schema)
	if err == nil {
		db.isInit = true
	}
	return err
}
