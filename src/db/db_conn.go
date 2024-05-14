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
	if !db.isConn {
		return fmt.Errorf("DBConn is not connected")
	}

	if db.isInit {
		return fmt.Errorf("DBConn is already initialized")
	}

	var shouldMigrate bool
	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		db.mu.Lock()
		defer db.mu.Unlock()
		shouldMigrate, err = db.createSchema()
		if err != nil {
			log.Printf("DBConn.init ERROR failed to create schema, %s", err)
			err = fmt.Errorf("failed to create schema")
		}
	}()
	wg.Wait()
	if err != nil {
		return err
	}

	if shouldMigrate {
		err = db.ApplyMigrations()
		if err != nil {
			log.Printf("DBConn.init ERROR failed to apply migrations, %s", err)
			return fmt.Errorf("failed to apply migrations")
		}
	}

	db.isInit = true
	return nil
}

func (db *DBConn) createSchema() (shouldMigrate bool, err error) {
	log.Println("createSchema TRACE in")
	schema := ""

	if db.UserTableExists() {
		log.Println("createSchema TRACE user table exists")
		shouldMigrate = true
	} else {
		log.Println("createSchema TRACE user table will be created")
		schema += UserSchema + "\n"
	}

	if db.AuthTableExists() {
		log.Println("createSchema TRACE auth table exists")
		shouldMigrate = true
	} else {
		log.Println("createSchema TRACE auth table will be created")
		schema += AuthSchema + "\n"
	}

	if db.AvatarTableExists() {
		log.Println("createSchema TRACE avatar table exists")
		shouldMigrate = true
	} else {
		log.Println("createSchema TRACE avatar table will be created")
		schema += AvatarSchema + "\n"
	}

	if db.MigrationsTableExists() {
		log.Println("createSchema TRACE migration table exists")
		// TODO meta-migrate, ie migrations migration
	} else {
		log.Println("createSchema TRACE migration table will be created")
		schema += MigrationSchema + "\n"
	}

	if schema == "" {
		log.Println("createSchema TRACE no tables to create")
		return true, nil
	}

	_, err = db.conn.Exec(strings.TrimRight(schema, "\n"))
	if err != nil {
		log.Printf("createSchema ERROR failed to create schema, %s", err.Error())
		err = fmt.Errorf("failed to create schema")
	}

	return shouldMigrate, err
}
