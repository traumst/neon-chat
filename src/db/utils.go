package db

import (
	"database/sql"
	"log"
)

func (dbConn *DBConn) TableExists(tableName string) bool {
	if tableName == "" {
		panic("check if table exist - empty name")
	}
	table := ""
	err := dbConn.Conn.Get(
		&table,
		"SELECT name FROM sqlite_master WHERE type='table' and name=?;",
		tableName)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			log.Printf("ERROR DBConn.TableExists checked[%s] found[%s], %s", tableName, tableName, err.Error())
			return false
		}
	}
	return table == tableName
}
