package db

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DBConn struct {
	// READONLY static conn
	Conn *sqlx.DB
	// READWRITE per session tx
	Tx *sqlx.Tx
	// ReqId from original request for tracing
	TxId string
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
	conn, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %s", err)
	}
	dbOptions := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA locking_mode = NORMAL;",
		// "PRAGMA auto_vacuum = INCREMENTAL;",
		"PRAGMA foreign_keys = ON;",
		"PRAGMA journal_size_limit = 67108864;",
		"PRAGMA page_size = 4096;",
		"PRAGMA cache_size = 2000;",
		"PRAGMA mmap_size = 134217728;",
	}
	log.Printf("applying db options...")
	for _, pragma := range dbOptions {
		log.Printf("applying db option [%s]", pragma)
		_, err := conn.Exec(pragma)
		if err != nil {
			log.Fatalf("db connection error applying pragma [%s]: %s", pragma, err)
		}
	}
	log.Printf("db connection established with connections[%d]", conn.Stats().MaxOpenConnections)
	dbConn := DBConn{Conn: conn}
	err = dbConn.init()
	return &dbConn, err
}

func (dbConn *DBConn) ConnClose() {
	// TODO count active tx
	dbConn.Conn.Close()
}

func (dbConn *DBConn) init() error {
	shouldMigrate, err := dbConn.createTables()
	if err != nil {
		log.Printf("DBConn.init ERROR failed to create schema, %s", err)
		return fmt.Errorf("failed to create schema")
	}
	if shouldMigrate {
		err = dbConn.ApplyMigrations()
		if err != nil {
			log.Printf("DBConn.init ERROR failed to apply migrations, %s", err)
			return fmt.Errorf("failed to apply migrations")
		}
	}
	err = dbConn.createIndex()
	if err != nil {
		log.Printf("DBConn.init ERROR failed to create schema, %s", err)
		return fmt.Errorf("failed to create schema")
	}

	return nil
}

func (dbConn *DBConn) createIndex() (err error) {
	// Note: update when adding new tables
	indecies := MigrationIndex +
		UserIndex + AuthIndex + AvatarSchema + ReservationIndex +
		ChatIndex + ChatUserIndex + MessageIndex + QuoteIndex
	if indecies == "" {
		log.Println("createIndex TRACE no indexes to create")
		return nil
	}

	_, err = dbConn.Conn.Exec(strings.TrimRight(indecies, "\n"))
	if err != nil {
		log.Printf("createSchema ERROR failed to create indexes, %s", err.Error())
		return fmt.Errorf("failed to create indexes")
	}
	return nil
}

func (dbConn *DBConn) createTables() (shouldMigrate bool, err error) {
	schema, shouldMigrate := dbConn.concatSchema()
	if schema == "" {
		log.Println("createTables TRACE no tables to create")
		return true, nil
	}

	_, err = dbConn.Conn.Exec(strings.TrimRight(schema, "\n"))
	if err != nil {
		log.Printf("createTables ERROR failed to create schema, %s", err.Error())
		err = fmt.Errorf("failed to create schema")
	}
	return shouldMigrate, err
}

// Note: update when adding new tables
func (dbConn *DBConn) concatSchema() (schema string, shouldMigrate bool) {
	if dbConn.MigrationsTableExists() {
		log.Println("concatSchema TRACE migration table exists")
		// TODO meta-migrate, ie migrations migration
	} else {
		log.Println("concatSchema TRACE migration table will be created")
		schema += MigrationSchema + "\n"
	}

	if dbConn.UserTableExists() {
		log.Println("concatSchema TRACE user table exists")
		shouldMigrate = true
	} else {
		log.Println("concatSchema TRACE user table will be created")
		schema += UserSchema + "\n"
	}

	if dbConn.AuthTableExists() {
		log.Println("concatSchema TRACE auth table exists")
		shouldMigrate = true
	} else {
		log.Println("concatSchema TRACE auth table will be created")
		schema += AuthSchema + "\n"
	}

	if dbConn.AvatarTableExists() {
		log.Println("concatSchema TRACE avatar table exists")
		shouldMigrate = true
	} else {
		log.Println("concatSchema TRACE avatar table will be created")
		schema += AvatarSchema + "\n"
	}

	if dbConn.ReservationTableExists() {
		log.Println("concatSchema TRACE reservation table exists")
		shouldMigrate = true
	} else {
		log.Println("concatSchema TRACE reservation table will be created")
		schema += ReservationSchema + "\n"
	}

	if dbConn.ChatTableExists() {
		log.Println("concatSchema TRACE chat table exists")
		shouldMigrate = true
	} else {
		log.Println("concatSchema TRACE chat table will be created")
		schema += ChatSchema + "\n"
	}

	if dbConn.ChatUserTableExists() {
		log.Println("concatSchema TRACE chat_user table exists")
		shouldMigrate = true
	} else {
		log.Println("concatSchema TRACE chat_user table will be created")
		schema += ChatUserSchema + "\n"
	}

	if dbConn.MessageTableExists() {
		log.Println("concatSchema TRACE messages table exists")
		shouldMigrate = true
	} else {
		log.Println("concatSchema TRACE messages table will be created")
		schema += MessageSchema + "\n"
	}

	if dbConn.QuoteTableExists() {
		log.Println("concatSchema TRACE quotes table exists")
		shouldMigrate = true
	} else {
		log.Println("concatSchema TRACE quotes table will be created")
		schema += QuoteSchema + "\n"
	}

	return schema, shouldMigrate
}
