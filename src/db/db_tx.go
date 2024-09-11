package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type DBTx struct {
	uid  string
	conn *DBConn
	tx   *sqlx.Tx
}

func (db *DBConn) OpenTx(uid string) (*DBTx, error) {
	log.Printf("TRACE [%s] OpenTx", uid)
	tx, err := db.conn.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to open transaction, %s", err)
	}
	return &DBTx{conn: db, tx: tx}, nil
}

func (dbtx *DBTx) CloseTx(err error) error {
	log.Printf("TRACE [%s] CloseTx, err[%d]", dbtx.uid, err)
	if err != nil {
		return dbtx.rollbackTx()
	} else {
		return dbtx.commitTx()
	}
}

func (dbtx *DBTx) commitTx() error {
	log.Printf("TRACE [%s] commitTx", dbtx.uid)
	if err := dbtx.tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction, %s", err)
	}
	return nil
}

func (dbtx *DBTx) rollbackTx() error {
	log.Printf("TRACE [%s] rollbackTx", dbtx.uid)
	if err := dbtx.tx.Rollback(); err != nil {
		return fmt.Errorf("failed to rollback transaction, %s", err)
	}
	return nil
}
