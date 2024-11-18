package db

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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
	// Is db locked for maintenance
	underMaintenance bool
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
	conn.DB.SetMaxOpenConns(1)
	conn.DB.SetMaxIdleConns(1)
	conn.DB.SetConnMaxLifetime(0)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %s", err)
	}
	dbOptions := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA locking_mode = NORMAL;",
		"PRAGMA auto_vacuum = INCREMENTAL;",
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
	log.Printf("db connection established with limit[%d]", conn.Stats().MaxOpenConnections)
	dbConn := DBConn{Conn: conn}
	err = dbConn.init()
	return &dbConn, err
}

func (dbConn *DBConn) ConnClose() {
	// TODO count active tx
	dbConn.Conn.Close()
}

func (dbConn *DBConn) IsAvailable(timeout time.Duration) bool {
	for dbConn.underMaintenance {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		select {
		case <-ticker.C:
			return false
		case <-time.After(timeout):
			return false
		default:
			return true
		}
	}

	return true
}

func (dbConn *DBConn) ScheduleMaintenance() {
	dbConn.underMaintenance = true
	defer func() { dbConn.underMaintenance = false }()
	// TODO wait for all active tx to finish
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for {
		// vacuum requires no other queries running
		// we allow only 1 concurrent connection,
		if _, err := dbConn.Conn.Exec("VACUUM"); err != nil {
			log.Printf("Error running VACUUM: %v", err)
		}
		if _, err := dbConn.Conn.Exec("ANALYZE"); err != nil {
			log.Printf("Error running ANALYZE: %v", err)
		}
		// block until next tick
		<-ticker.C
	}
}

func (dbConn *DBConn) init() error {
	shouldMigrate, err := dbConn.createTables()
	if err != nil {
		log.Printf("ERROR DBConn.init failed to create schema, %s", err)
		return fmt.Errorf("failed to create schema")
	}
	if shouldMigrate {
		err = dbConn.ApplyMigrations()
		if err != nil {
			log.Printf("ERROR DBConn.init failed to apply migrations, %s", err)
			return fmt.Errorf("failed to apply migrations")
		}
	}
	err = dbConn.createIndex()
	if err != nil {
		log.Printf("ERROR DBConn.init failed to create schema, %s", err)
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
		log.Println("TRACE createIndex no indexes to create")
		return nil
	}

	_, err = dbConn.Conn.Exec(strings.TrimRight(indecies, "\n"))
	if err != nil {
		log.Printf("ERROR createSchema failed to create indexes, %s", err.Error())
		return fmt.Errorf("failed to create indexes")
	}
	return nil
}

func (dbConn *DBConn) createTables() (shouldMigrate bool, err error) {
	schema, shouldMigrate := dbConn.concatSchema()
	if schema == "" {
		log.Println("TRACE createTables no tables to create")
		return true, nil
	}

	_, err = dbConn.Conn.Exec(strings.TrimRight(schema, "\n"))
	if err != nil {
		log.Printf("ERROR createTables failed to create schema, %s", err.Error())
		err = fmt.Errorf("failed to create schema")
	}
	return shouldMigrate, err
}

// Note: update when adding new tables
func (dbConn *DBConn) concatSchema() (schema string, shouldMigrate bool) {
	if dbConn.MigrationsTableExists() {
		log.Println("TRACE concatSchema migration table exists")
		// TODO meta-migrate, ie migrations migration
	} else {
		log.Println("TRACE concatSchema migration table will be created")
		schema += MigrationSchema + "\n"
	}

	if dbConn.UserTableExists() {
		log.Println("TRACE concatSchema user table exists")
		shouldMigrate = true
	} else {
		log.Println("TRACE concatSchema user table will be created")
		schema += UserSchema + "\n"
	}

	if dbConn.AuthTableExists() {
		log.Println("TRACE concatSchema auth table exists")
		shouldMigrate = true
	} else {
		log.Println("TRACE concatSchema auth table will be created")
		schema += AuthSchema + "\n"
	}

	if dbConn.AvatarTableExists() {
		log.Println("TRACE concatSchema avatar table exists")
		shouldMigrate = true
	} else {
		log.Println("TRACE concatSchema avatar table will be created")
		schema += AvatarSchema + "\n"
	}

	if dbConn.ReservationTableExists() {
		log.Println("TRACE concatSchema reservation table exists")
		shouldMigrate = true
	} else {
		log.Println("TRACE concatSchema reservation table will be created")
		schema += ReservationSchema + "\n"
	}

	if dbConn.ChatTableExists() {
		log.Println("TRACE concatSchema chat table exists")
		shouldMigrate = true
	} else {
		log.Println("TRACE concatSchema chat table will be created")
		schema += ChatSchema + "\n"
	}

	if dbConn.ChatUserTableExists() {
		log.Println("TRACE concatSchema chat_user table exists")
		shouldMigrate = true
	} else {
		log.Println("TRACE concatSchema chat_user table will be created")
		schema += ChatUserSchema + "\n"
	}

	if dbConn.MessageTableExists() {
		log.Println("TRACE concatSchema messages table exists")
		shouldMigrate = true
	} else {
		log.Println("TRACE concatSchema messages table will be created")
		schema += MessageSchema + "\n"
	}

	if dbConn.QuoteTableExists() {
		log.Println("TRACE concatSchema quotes table exists")
		shouldMigrate = true
	} else {
		log.Println("TRACE concatSchema quotes table will be created")
		schema += QuoteSchema + "\n"
	}

	return schema, shouldMigrate
}
