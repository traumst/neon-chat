package db

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Quote struct {
	MsgId   uint `db:"msg_id"`
	QuoteId uint `db:"quote_id"`
}

const QuoteSchema = `
	CREATE TABLE IF NOT EXISTS quotes (
		msg_id INTEGER, 
		quote_id INTEGER,
		FOREIGN KEY(msg_id) REFERENCES messages(id),
		FOREIGN KEY(quote_id) REFERENCES messages(id)
	);`

const QuoteIndex = ``

func (db *DBConn) QuoteTableExists() bool {
	return db.TableExists("quotes")
}

func (db *DBConn) AddQuote(quote *Quote) (*Quote, error) {
	if quote.MsgId == 0 || quote.QuoteId == 0 {
		return nil, fmt.Errorf("bad arg - msg_id[%d] quote_id[%d]", quote.MsgId, quote.QuoteId)
	}
	if !db.ConnIsActive() {
		return nil, fmt.Errorf("db is not connected")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	result, err := db.conn.Exec(`INSERT INTO quotes (msg_id, quote_id) VALUES (?, ?)`, quote.MsgId, quote.QuoteId)
	if err != nil {
		return nil, fmt.Errorf("error adding quote: sqlx %s", err)
	} else if affected, _ := result.RowsAffected(); affected != 1 {
		return nil, fmt.Errorf("error adding quote: no rows affected")
	}

	return quote, nil
}

func (db *DBConn) GetQuote(msgId uint) (*Quote, error) {
	if msgId == 0 {
		return nil, fmt.Errorf("bad input: msgId[%d]", msgId)
	}
	if !db.ConnIsActive() {
		return nil, fmt.Errorf("db is not connected")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	var quote Quote
	err := db.conn.Get(&quote, `SELECT * FROM quotes WHERE msg_id = ?`, msgId)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("error getting quote: %s", err)
	} else if quote.MsgId != msgId {
		return nil, fmt.Errorf("error getting quote: expected[%d] different from actual[%d]", msgId, quote.MsgId)
	} else if quote.QuoteId == 0 {
		return nil, nil
	}

	return &quote, nil
}

func (db *DBConn) GetQuotes(msgIds []uint) ([]*Quote, error) {
	if msgIds == nil {
		return nil, fmt.Errorf("bad input: msgIds[%v]", msgIds)
	} else if len(msgIds) == 0 {
		return nil, nil
	}
	if !db.ConnIsActive() {
		return nil, fmt.Errorf("db is not connected")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	query, args, err := sqlx.In(`SELECT * FROM quotes WHERE msg_id IN (?)`, msgIds)
	if err != nil {
		return nil, fmt.Errorf("error preparing select quotes query for msgIds %v, %s", msgIds, err)
	}
	query = db.conn.Rebind(query)

	var quotes []Quote
	err = db.conn.Select(&quotes, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error getting quotes for msgIds %v: %s", msgIds, err)
	}
	var quotePtrs []*Quote
	for _, quote := range quotes {
		quotePtrs = append(quotePtrs, &quote)
	}
	return quotePtrs, nil
}
