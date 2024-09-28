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

func (dbConn *DBConn) MigrationsTableExists() bool {
	return dbConn.TableExists("_migrations")
}

func (dbConn *DBConn) ApplyMigrations() error {
	log.Printf("TRACE applyMigrations IN")
	//utils.LS()
	// TODO load "latest" subset
	files, err := utils.GetFilenamesIn(migraitonsFolder)
	if err != nil {
		return fmt.Errorf("DBConn.ApplyMigrations failed to list migrations[%s], %s", migraitonsFolder, err.Error())
	}

	for _, filename := range files {
		err := applyMigration(dbConn, filename)
		if err != nil {
			log.Printf("DBConn.ApplyMigrations fail while applying migrations, %s", err)
			return fmt.Errorf("fail to apply migration[%s]", filename)
		}
	}
	log.Printf("TRACE applyMigrations OUT")
	return nil
}

func applyMigration(dbConn *DBConn, filename string) error {
	log.Printf("TRACE applyMigration now on [%s]", filename)
	path := strings.Split(filename, ".")
	if len(path) != 2 {
		return fmt.Errorf("migration title[%s] is not *.sql", filename)
	}
	title := path[0]
	if title == "" {
		log.Printf("WARN applyMigration blank title [%s]", filename)
		return nil
	}
	if ext := path[1]; ext != "sql" {
		log.Printf("TRACE applyMigration skip non-sql [%s]", title)
		return nil
	}
	log.Printf("TRACE applyMigration check if already applied [%s]", title)
	isApplied, err := isMigrationApplied(dbConn, title)
	if err == nil && isApplied {
		return nil
	}
	log.Printf("TRACE applyMigration reading [%s]", title)
	bytes, err := os.ReadFile(migraitonsFolder + "/" + filename)
	if err != nil {
		return fmt.Errorf("failed to read migration file content[%s]", title)
	}
	log.Printf("TRACE applyMigration storing migration [%s]", title)
	migration, err := addMigration(dbConn, string(bytes[:]), title)
	if err != nil || migration.Id < 1 {
		return fmt.Errorf("failed to apply migration[%s], %s", title, err)
	}
	log.Printf("TRACE applyMigration applied[%d][%s]", migration.Id, migration.Title)
	return nil
}

func addMigration(dbConn *DBConn, migrate string, title string) (*Migration, error) {
	log.Printf("TRACE applyMigration executing [%s]", title)
	if _, err := dbConn.Conn.Exec(migrate); err != nil {
		return nil, fmt.Errorf("failed to add migration[%s], %s", title, err.Error())
	}
	result, err := dbConn.Conn.Exec(`INSERT INTO _migrations (title) VALUES (?)`, title)
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

func isMigrationApplied(dbConn *DBConn, title string) (bool, error) {
	if title == "" {
		return false, fmt.Errorf("migration title is empty")
	}

	var migration Migration
	err := dbConn.Conn.Get(&migration, `SELECT * FROM _migrations WHERE title=?`, title)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, fmt.Errorf("error getting migrations: %s", err)
		}
	}
	return migration.Id > 0, nil
}
