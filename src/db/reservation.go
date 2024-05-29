package db

import (
	"database/sql"
	"fmt"
	"time"
)

type Reservation struct {
	Id     uint      `db:"id"`
	UserId uint      `db:"user_id"`
	Token  string    `db:"token"`
	Expire time.Time `db:"expire"`
}

const ReservationSchema = `
	CREATE TABLE IF NOT EXISTS reservations (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		user_id INTEGER UNIQUE,
		token TEXT NOT NULL UNIQUE,
		expire DATETIME NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`
const ReservationIndex = `CREATE INDEX IF NOT EXISTS idx_reserve_expire ON reservations(expire);`

func (db *DBConn) ReservationTableExists() bool {
	return db.TableExists("reservations")
}

func (db *DBConn) AddReservation(reserve Reservation) (*Reservation, error) {
	if !db.IsActive() {
		return nil, fmt.Errorf("db is not connected")
	}
	if reserve.Id != 0 {
		return nil, fmt.Errorf("reserve already has an id[%d]", reserve.Id)
	} else if reserve.UserId <= 0 {
		return nil, fmt.Errorf("reserve has no user id")
	} else if reserve.Token == "" {
		return nil, fmt.Errorf("reserve token is empty")
	} else if reserve.Expire.IsZero() {
		return nil, fmt.Errorf("reserve expiration is zero")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	result, err := db.conn.Exec(`INSERT INTO reservations (user_id, token, expire) VALUES (?, ?, ?) 
		ON CONFLICT(user_id) DO UPDATE 
			SET token = excluded.token, 
				expire = excluded.expire;`,
		reserve.UserId, reserve.Token, reserve.Expire)
	if err != nil {
		return nil, fmt.Errorf("error adding reserve: %s", err.Error())
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %s", err.Error())
	}

	reserve.Id = uint(lastId)
	return &reserve, nil
}

func (db *DBConn) GetReservation(token string) (*Reservation, error) {
	if !db.IsActive() {
		return nil, fmt.Errorf("db is not connected")
	}
	if token == "" {
		return nil, fmt.Errorf("invalid token")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	reserve := Reservation{}
	err := db.conn.Get(&reserve, `SELECT * FROM reservations WHERE token=?`, token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf(", %s", err.Error())
	}
	return &reserve, nil
}