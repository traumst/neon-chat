package db

import (
	"fmt"
	"log"
)

func (db *DBConn) OpenTx(uid string) error {
	log.Printf("TRACE [%s] OpenTx", uid)
	tx, err := db.conn.Beginx()
	if err != nil {
		return fmt.Errorf("failed to open transaction, %s", err)
	}
	db.tx = tx
	db.txId = uid
	return nil
}

func (db *DBConn) CloseTx(err error) error {
	log.Printf("TRACE [%s] CloseTx, err[%d]", db.txId, err)
	if err != nil {
		return db.rollbackTx()
	} else {
		return db.commitTx()
	}
}

func (db *DBConn) commitTx() error {
	log.Printf("TRACE [%s] commitTx", db.txId)
	if err := db.tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction, %s", err)
	}
	return nil
}

func (db *DBConn) rollbackTx() error {
	log.Printf("TRACE [%s] rollbackTx", db.txId)
	if err := db.tx.Rollback(); err != nil {
		return fmt.Errorf("failed to rollback transaction, %s", err)
	}
	return nil
}
