package db

import (
	"fmt"
	"log"
	"neon-chat/src/utils"
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

func (dbConn *DBConn) ConnClose(timeout time.Duration) error {
	activeCount := utils.MaintenanceManager.WaitUsersLeave(timeout)
	if activeCount != 0 {
		log.Printf("WARN ConnClose aborts [%d] active users after timeout %s", activeCount, timeout.String())
	} else {
		log.Printf("INFO ConnClose waited for users to leave")
	}
	return dbConn.Conn.Close()
}

func (dbConn *DBConn) ScheduleMaintenance() {
	err := doMaintenance(dbConn)
	if err != nil {
		log.Fatalf("ERROR startup maintenance failed to run, %s", err)
	}

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		err = doMaintenance(dbConn)
		if err != nil {
			log.Printf("ERROR maintenance failed to run, %s", err)
		}
	}
}

func doMaintenance(dbConn *DBConn) error {
	log.Println("TRACE maintenance about to start")
	if err := utils.MaintenanceManager.RaiseFlag(); err != nil {
		log.Printf("TRACE maintenance will not run, %s", err)
	}
	defer utils.MaintenanceManager.ClearFlag()

	activeCount := utils.MaintenanceManager.WaitUsersLeave(1 * time.Minute)
	if activeCount != 0 {
		return fmt.Errorf("failed while waiting for [%d] users to leave", activeCount)
	}

	log.Println("TRACE running incremental_vacuum")
	if _, err := dbConn.Conn.Exec("PRAGMA incremental_vacuum"); err != nil {
		log.Printf("ERROR running incremental_vacuum: %v", err)
	}
	log.Println("TRACE running analyze")
	if _, err := dbConn.Conn.Exec("ANALYZE"); err != nil {
		log.Printf("ERROR running ANALYZE: %v", err)
	}
	log.Println("TRACE maintenance done")
	return nil
}

func (dbConn *DBConn) init() error {
	tables, index := dbConn.concatSchema()
	var total int64

	rows, err := dbConn.createTables(tables)
	if err != nil {
		log.Printf("ERROR DBConn.init failed to create tables %s", tables)
		return fmt.Errorf("failed to create tables, %s", err)
	} else {
		log.Printf("INFO DBConn.init created %d tables", rows)
		total += rows
	}

	rows, err = dbConn.createIndex(index)
	if err != nil {
		log.Printf("ERROR DBConn.init failed to create index %s", index)
		return fmt.Errorf("failed to create index, %s", err)
	} else {
		log.Printf("INFO DBConn.init created %d indexes", rows)
		total += rows
	}

	rows, err = dbConn.TryApplyMigrations()
	if err != nil {
		log.Printf("ERROR DBConn.init failed to apply migrations, %s", err)
		return fmt.Errorf("failed to apply migrations")
	} else {
		log.Printf("INFO DBConn.init applied %d migrations", rows)
	}

	return nil
}

func (dbConn *DBConn) createTables(schema string) (int64, error) {
	if schema == "" {
		log.Println("TRACE createTables no tables to create")
		return 0, nil
	}
	res, err := dbConn.Conn.Exec(strings.TrimRight(schema, "\n"))
	if err != nil {
		return 0, fmt.Errorf("failed to create tables, %s", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return rowsAffected, fmt.Errorf("failed to get rows affected, %s", err)
	}
	return rowsAffected, nil
}

func (dbConn *DBConn) createIndex(index string) (int64, error) {
	if index == "" {
		log.Println("TRACE createIndex no indexes to create")
		return 0, nil
	}
	res, err := dbConn.Conn.Exec(strings.TrimRight(index, "\n"))
	if err != nil {
		return 0, fmt.Errorf("failed to create indexes, %s", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return rowsAffected, fmt.Errorf("failed to get rows affected, %s", err)
	}
	return rowsAffected, nil
}
