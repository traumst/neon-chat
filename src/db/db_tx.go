package db

import (
	"fmt"
	"log"
)

func (db *DBConn) OpenTx(sessionId string) error {
	log.Printf("TRACE [%s] OpenTx", sessionId)
	tx, err := db.conn.Beginx()
	if err != nil {
		return fmt.Errorf("failed to open transaction, %s", err)
	}
	db.tx = tx
	db.txId = sessionId
	return nil
}

func (db *DBConn) CloseTx(err error, hasChanges bool) error {
	log.Printf("TRACE [%s] CloseTx hasChanges[%t], err[%d]", db.txId, hasChanges, err)
	if err != nil || !hasChanges {
		// only discards im-memory changes, no db io
		return db.rollbackTx()
	} else {
		// must write to database's journal even if tx is empty
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
