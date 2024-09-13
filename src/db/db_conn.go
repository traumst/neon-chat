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
	Conn *sqlx.DB
	TxId string
	Tx   *sqlx.Tx
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
	db := DBConn{Conn: conn}
	err = db.init()
	return &db, err
}

func (db *DBConn) ConnClose() {
	// TODO count active tx
	db.Conn.Close()
}

func (db *DBConn) init() error {
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

	return nil
}

func (db *DBConn) createIndex() (err error) {
	// Note: update when adding new tables
	indecies := MigrationIndex +
		UserIndex + AuthIndex + AvatarSchema + ReservationIndex +
		ChatIndex + ChatUserIndex + MessageIndex + QuoteIndex
	if indecies == "" {
		log.Println("createIndex TRACE no indexes to create")
		return nil
	}

	_, err = db.Conn.Exec(strings.TrimRight(indecies, "\n"))
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

	_, err = db.Conn.Exec(strings.TrimRight(schema, "\n"))
	if err != nil {
		log.Printf("createTables ERROR failed to create schema, %s", err.Error())
		err = fmt.Errorf("failed to create schema")
	}
	return shouldMigrate, err
}

// Note: update when adding new tables
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

	if db.MessageTableExists() {
		log.Println("concatSchema TRACE messages table exists")
		shouldMigrate = true
	} else {
		log.Println("concatSchema TRACE messages table will be created")
		schema += MessageSchema + "\n"
	}

	if db.QuoteTableExists() {
		log.Println("concatSchema TRACE quotes table exists")
		shouldMigrate = true
	} else {
		log.Println("concatSchema TRACE quotes table will be created")
		schema += QuoteSchema + "\n"
	}

	return schema, shouldMigrate
}
