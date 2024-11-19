package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

func (dbConn *DBConn) OpenTx(txId string) (*sqlx.Tx, string, error) {
	log.Printf("TRACE [%s] OpenTx", txId)
	tx, err := dbConn.Conn.Beginx()
	if err != nil {
		return nil, "", fmt.Errorf("failed to open transaction, %s", err)
	}
	dbConn.txCount += 1
	dbConn.Tx = tx
	dbConn.TxId = txId
	return tx, txId, nil
}

func (dbConn *DBConn) CloseTx(err error, hasChanges bool) error {
	log.Printf("TRACE [%s] CloseTx hasChanges[%t], err[%d]", dbConn.TxId, hasChanges, err)
	if err != nil || !hasChanges {
		// only discards im-memory changes, no db io
		return dbConn.rollbackTx()
	} else {
		// must write to database's journal even if tx is empty
		return dbConn.commitTx()
	}
}

func (dbConn *DBConn) commitTx() error {
	log.Printf("TRACE [%s] commitTx", dbConn.TxId)
	if err := dbConn.Tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction, %s", err)
	}
	return nil
}

func (dbConn *DBConn) rollbackTx() error {
	log.Printf("TRACE [%s] rollbackTx", dbConn.TxId)
	if err := dbConn.Tx.Rollback(); err != nil {
		return fmt.Errorf("failed to rollback transaction, %s", err)
	}
	return nil
}
