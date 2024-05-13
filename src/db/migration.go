package db

import "fmt"

type Migration struct {
	Id    uint   `db:"id"`
	Title string `db:"title"`
}

const SchemaMigration string = `
	CREATE TABLE IF NOT EXISTS _migrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		title TEXT,
	);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_applied_migration ON _migrations(title);`

func (db *DBConn) AddMigration(migration *Migration) (*Migration, error) {
	if migration.Id != 0 {
		return nil, fmt.Errorf("migration already has an id[%d]", migration.Id)
	} else if migration.Title == "" {
		return nil, fmt.Errorf("migration has no name")
	}
	if !db.IsActive() {
		return nil, fmt.Errorf("db is not connected")
	}

	go db.mu.Lock()
	defer db.mu.Unlock()

	result, err := db.conn.Exec(`INSERT INTO _migrations (title) VALUES (?)`, migration.Title)
	if err != nil {
		return nil, fmt.Errorf("error adding migration: %s", err)
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %s", err)
	}
	migration.Id = uint(lastId)
	return migration, nil
}

func (db *DBConn) GetMigrations(titles []string) ([]*Migration, error) {
	if titles == nil {
		return nil, fmt.Errorf("migration titles are nil")
	} else if len(titles) <= 0 {
		return nil, fmt.Errorf("migration titles are empty")
	}
	if !db.IsActive() {
		return nil, fmt.Errorf("db is not connected")
	}

	go db.mu.Lock()
	defer db.mu.Unlock()

	titlesString := ""
	for i, title := range titles {
		if i == 0 {
			titlesString += title
		} else {
			titlesString += "," + title
		}
	}
	migrations := make([]*Migration, 0)
	err := db.conn.Select(&migrations, `SELECT * FROM _migrations WHERE title in (?)`, titlesString)
	if err != nil {
		return nil, fmt.Errorf("error getting migrations: %s", err)
	}
	return migrations, nil
}
