package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"neon-chat/src/utils"
)

type Migration struct {
	Id    uint      `db:"id"`
	Title string    `db:"title"`
	Stamp time.Time `db:"stamp"`
}

const MigrationSchema = `
	CREATE TABLE IF NOT EXISTS _migrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		title TEXT,
		stamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
const MigrationIndex = `CREATE UNIQUE INDEX IF NOT EXISTS idx_applied_migration ON _migrations(title);`

func (db *DBConn) MigrationsTableExists() bool {
	return db.TableExists("_migrations")
}

func (db *DBConn) ApplyMigrations() error {
	log.Printf("applyMigrations TRACE IN")
	utils.LS()
	// TODO load "latest" subset
	files, err := utils.GetFilenamesIn(migraitonsFolder)
	if err != nil {
		return fmt.Errorf("DBConn.ApplyMigrations failed to list migrations[%s], %s", migraitonsFolder, err.Error())
	}

	for _, filename := range files {
		err := applyMigration(db, filename)
		if err != nil {
			log.Printf("DBConn.ApplyMigrations fail while applying migrations, %s", err)
			return fmt.Errorf("fail to apply migration[%s]", filename)
		}
	}
	log.Printf("applyMigrations TRACE OUT")
	return nil
}

func applyMigration(db *DBConn, filename string) error {
	log.Printf("applyMigration TRACE now on [%s]", filename)
	path := strings.Split(filename, ".")
	if len(path) != 2 {
		return fmt.Errorf("migration title[%s] is not *.sql", filename)
	}
	title := path[0]
	if title == "" {
		log.Printf("applyMigration WARN blank title [%s]", filename)
		return nil
	}
	if ext := path[1]; ext != "sql" {
		log.Printf("applyMigration TRACE skip non-sql [%s]", title)
		return nil
	}
	log.Printf("applyMigration TRACE check if already applied [%s]", title)
	isApplied, err := isMigrationApplied(db, title)
	if err == nil && isApplied {
		return nil
	}
	log.Printf("applyMigration TRACE reading [%s]", title)
	bytes, err := os.ReadFile(migraitonsFolder + "/" + filename)
	if err != nil {
		return fmt.Errorf("failed to read migration file content[%s]", title)
	}
	log.Printf("applyMigration TRACE storing migration [%s]", title)
	migration, err := addMigration(db, string(bytes[:]), title)
	if err != nil || migration.Id < 1 {
		return fmt.Errorf("failed to apply migration[%s], %s", title, err)
	}
	log.Printf("applyMigration TRACE applied[%d][%s]", migration.Id, migration.Title)
	return nil
}

func addMigration(db *DBConn, migrate string, title string) (*Migration, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	log.Printf("applyMigration TRACE executing [%s]", title)
	if _, err := db.conn.Exec(migrate); err != nil {
		return nil, fmt.Errorf("failed to add migration[%s], %s", title, err.Error())
	}
	result, err := db.conn.Exec(`INSERT INTO _migrations (title) VALUES (?)`, title)
	if err != nil {
		return nil, fmt.Errorf("error adding migration: %s", err)
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %s", err)
	}
	migration := Migration{
		Id:    uint(lastId),
		Title: title,
	}
	return &migration, nil
}

func isMigrationApplied(db *DBConn, title string) (bool, error) {
	if title == "" {
		return false, fmt.Errorf("migration title is empty")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	var migration Migration
	err := db.conn.Get(&migration, `SELECT * FROM _migrations WHERE title=?`, title)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, fmt.Errorf("error getting migrations: %s", err)
		}
	}
	return migration.Id > 0, nil
}
