package db

import (
	"fmt"
	"log"
	"os"
	"strings"
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

const migraitonsFolder string = "./src/db/migrations"

func ConnectDB(dbPath string) (*DBConn, error) {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			return nil, fmt.Errorf("error creating db file: %s", err)
		}
		file.Close()
		log.Printf("db file created [%s]", dbPath)
	} else {
		log.Printf("db file exists [%s]", dbPath)
	}

	log.Printf("db connects to [%s]", dbPath)
	conn, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %s", err)
	}

	log.Printf("db connection established with connections[%d]", conn.Stats().MaxOpenConnections)
	db := DBConn{conn: conn, isConn: true, isInit: false}
	err = db.init()
	return &db, err
}

func (db *DBConn) ConnClose() {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.conn.Close()
}

func (db *DBConn) ConnIsActive() bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	return db.isConn && db.isInit
}

func (db *DBConn) init() error {
	if !db.isConn {
		return fmt.Errorf("DBConn is not connected")
	}
	if db.isInit {
		return fmt.Errorf("DBConn is already initialized")
	}

	shouldMigrate, err := db.createTables()
	if err != nil {
		log.Printf("DBConn.init ERROR failed to create schema, %s", err)
		return fmt.Errorf("failed to create schema")
	}
	if shouldMigrate {
		err = db.ApplyMigrations()
		if err != nil {
			log.Printf("DBConn.init ERROR failed to apply migrations, %s", err)
			return fmt.Errorf("failed to apply migrations")
		}
	}
	err = db.createIndex()
	if err != nil {
		log.Printf("DBConn.init ERROR failed to create schema, %s", err)
		return fmt.Errorf("failed to create schema")
	}

	db.isInit = true
	return nil
}

func (db *DBConn) createIndex() (err error) {
	indecies := MigrationIndex + UserIndex + AuthIndex + AvatarSchema + ReservationIndex
	if indecies == "" {
		log.Println("createIndex TRACE no indexes to create")
		return nil
	}

	db.mu.Lock()
	defer db.mu.Unlock()
	_, err = db.conn.Exec(strings.TrimRight(indecies, "\n"))
	if err != nil {
		log.Printf("createSchema ERROR failed to create indexes, %s", err.Error())
		return fmt.Errorf("failed to create indexes")
	}
	return nil
}

func (db *DBConn) createTables() (shouldMigrate bool, err error) {
	schema, shouldMigrate := db.concatSchema()
	if schema == "" {
		log.Println("createTables TRACE no tables to create")
		return true, nil
	}

	db.mu.Lock()
	defer db.mu.Unlock()
	_, err = db.conn.Exec(strings.TrimRight(schema, "\n"))
	if err != nil {
		log.Printf("createTables ERROR failed to create schema, %s", err.Error())
		err = fmt.Errorf("failed to create schema")
	}
	return shouldMigrate, err
}

func (db *DBConn) concatSchema() (schema string, shouldMigrate bool) {
	if db.MigrationsTableExists() {
		log.Println("concatSchema TRACE migration table exists")
		// TODO meta-migrate, ie migrations migration
	} else {
		log.Println("concatSchema TRACE migration table will be created")
		schema += MigrationSchema + "\n"
	}

	if db.UserTableExists() {
		log.Println("concatSchema TRACE user table exists")
		shouldMigrate = true
	} else {
		log.Println("concatSchema TRACE user table will be created")
		schema += UserSchema + "\n"
	}

	if db.AuthTableExists() {
		log.Println("concatSchema TRACE auth table exists")
		shouldMigrate = true
	} else {
		log.Println("concatSchema TRACE auth table will be created")
		schema += AuthSchema + "\n"
	}

	if db.AvatarTableExists() {
		log.Println("concatSchema TRACE avatar table exists")
		shouldMigrate = true
	} else {
		log.Println("concatSchema TRACE avatar table will be created")
		schema += AvatarSchema + "\n"
	}

	if db.ReservationTableExists() {
		log.Println("concatSchema TRACE reservation table exists")
		shouldMigrate = true
	} else {
		log.Println("concatSchema TRACE reservation table will be created")
		schema += ReservationSchema + "\n"
	}

	if db.ChatTableExists() {
		log.Println("concatSchema TRACE chat table exists")
		shouldMigrate = true
	} else {
		log.Println("concatSchema TRACE chat table will be created")
		schema += ChatSchema + "\n"
	}

	if db.ChatUserTableExists() {
		log.Println("concatSchema TRACE chat_user table exists")
		shouldMigrate = true
	} else {
		log.Println("concatSchema TRACE chat_user table will be created")
		schema += ChatUserSchema + "\n"
	}

	return schema, shouldMigrate
}
