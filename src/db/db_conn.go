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

//const migraitonsFolder string = "./migrations"

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

	db.mu.Lock()
	defer db.mu.Unlock()

	log.Println("  initiating db schema")
	schema := fmt.Sprintf("%s\n%s\n%s", SchemaUser, SchemaAuth, SchemaAvatar)
	_, err := db.conn.Exec(schema)
	if err != nil {
		return err
	}

	//err = applyMigrations(db)
	// if err == nil {
	// 	db.isInit = true
	// }
	db.isInit = true
	return err
}

// func applyMigrations(db *DBConn) error {
// 	log.Printf("applyMigrations TRACE IN")
// 	// TODO load "latest" subset
// 	files, err := utils.GetFilenamesIn(migraitonsFolder)
// 	if err != nil {
// 		return fmt.Errorf("failed to list migrations[%s], %s", migraitonsFolder, err.Error())
// 	}
// 	// TODO execute in batches
// 	for _, filename := range files {
// 		log.Printf("	applyMigrations TRACE now on [%s]", filename)
// 		path := strings.Split(filename, ".")
// 		if len(path) != 2 {
// 			return fmt.Errorf("migration title[%s] is not *.sql", filename)
// 		}
// 		title := path[0]
// 		if title == "" {
// 			log.Printf("	applyMigrations WARN blank title [%s]", filename)
// 			continue
// 		}
// 		if ext := path[1]; ext != "sql" {
// 			log.Printf("	applyMigrations TRACE skip non-sql [%s]", filename)
// 			continue
// 		}
// 		bytes, err := os.ReadFile(migraitonsFolder + "/" + title)
// 		if err != nil {
// 			return fmt.Errorf("failed to read migration file content[%s]", title)
// 		}
// 		if _, err = db.conn.Exec(string(bytes[:])); err != nil {
// 			return fmt.Errorf("failed to apply migration[%s]", title)
// 		}
// 		migration, err := db.AddMigration(&Migration{Title: title})
// 		if err != nil || migration.Id < 1 {
// 			return fmt.Errorf("failed to apply migraiton[%s], %s", title, err)
// 		}
// 	}
// 	log.Printf("applyMigrations TRACE OUT")
// 	return nil
// }
