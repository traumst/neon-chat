package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

const minUserNameLen = 3
const maxUserNameLen = 64

type User struct {
	Id     uint   `db:"id"`
	Name   string `db:"name"`
	Email  string `db:"email"`
	Type   string `db:"type"`
	Status string `db:"status"`
	Salt   string `db:"salt"`
}

const UserSchema = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		name TEXT UNIQUE, 
		email TEXT UNIQUE,
		type TEXT,
		status TEXT,
		salt INTEGER
	);`
const UserIndex = `CREATE INDEX IF NOT EXISTS idx_user_status ON users(status);`

func (dbConn *DBConn) UserTableExists() bool {
	return dbConn.TableExists("users")
}

func AddUser(dbConn sqlx.Ext, user *User) (*User, error) {
	if user.Id != 0 {
		return nil, fmt.Errorf("user already has an id[%d]", user.Id)
	} else if len(user.Name) < minUserNameLen || len(user.Name) > maxUserNameLen {
		return nil, fmt.Errorf("user has no name")
	} else if len(user.Email) < minUserNameLen || len(user.Email) > maxUserNameLen {
		return nil, fmt.Errorf("user has no email")
	} else if len(user.Type) < minUserNameLen || len(user.Type) > maxUserNameLen {
		return nil, fmt.Errorf("user has no type")
	} else if len(user.Status) < minUserNameLen || len(user.Status) > maxUserNameLen {
		return nil, fmt.Errorf("user has no status")
	} else if len(user.Salt) == 0 {
		return nil, fmt.Errorf("user has no salt")
	}

	n := user.Name
	e := user.Email
	t := user.Type
	status := user.Status
	salt := user.Salt
	log.Println("...adding user:", n, "email:", e, "type:", t, "status:", status, "salt:", salt)
	result, err := dbConn.Exec(`INSERT INTO users (name, email, type, status, salt) VALUES (?, ?, ?, ?, ?)`,
		n, e, t, status, salt)
	if err != nil {
		return nil, fmt.Errorf("error adding user: %s", err)
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %s", err)
	}
	user.Id = uint(lastId)
	return user, nil
}

func DropUser(dbConn sqlx.Ext, userId uint) error {
	if userId <= 0 {
		return fmt.Errorf("user does not have an id[%d]", userId)
	}

	_, err := dbConn.Exec(`DELETE FROM users where id = ?`, userId)
	if err != nil {
		return fmt.Errorf("error deleting user: %s", err.Error())
	}
	return nil
}

func SearchUser(dbConn sqlx.Ext, login string) (*User, error) {
	if len(login) < minUserNameLen || len(login) > maxUserNameLen {
		return nil, fmt.Errorf("login name/email was not provided")
	}

	var user User
	err := sqlx.Get(dbConn, &user,
		`SELECT * FROM users WHERE name = ? or email = ?`,
		login, login)
	if err != nil {
		return nil, fmt.Errorf("login[%s] not found: %s", login, err)
	}
	return &user, err
}

func SearchUsers(dbConn sqlx.Ext, terms []string) ([]*User, error) {
	if len(terms) <= 0 {
		return nil, fmt.Errorf("name was not provided")
	}

	query, args, err := sqlx.In(`SELECT * FROM users WHERE name IN (?) or email IN (?)`, terms, terms)
	if err != nil {
		return nil, fmt.Errorf("error preparing search query by user names %v, %s", terms, err)
	}
	query = dbConn.Rebind(query)

	var users []*User
	err = sqlx.Select(dbConn, &users, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error getting avatars for %v: %s", terms, err)
	}
	return users, nil
}

func GetUser(dbConn sqlx.Ext, id uint) (*User, error) {
	var user User
	err := sqlx.Get(dbConn, &user, `SELECT * FROM users WHERE id = ?`, id)
	if err != nil {
		return nil, fmt.Errorf("userId[%d] not found: %s", id, err)
	}
	return &user, nil
}

func GetUsers(dbConn sqlx.Ext, userIds []uint) ([]User, error) {
	if len(userIds) <= 0 {
		return []User{}, nil
	}

	query, args, err := sqlx.In(`SELECT * FROM users WHERE id IN (?)`, userIds)
	if err != nil {
		return nil, fmt.Errorf("error preparing query with userIds[%v]: %s", userIds, err)
	}
	query = dbConn.Rebind(query)

	var userChats []User
	err = sqlx.Select(dbConn, &userChats, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error getting users[%v]: %s", userIds, err)
	}
	return userChats, nil
}

func UpdateUserName(dbConn sqlx.Ext, userId uint, userName string) error {
	if userId <= 0 {
		return fmt.Errorf("invalid user id[%d]", userId)
	}

	result, err := dbConn.Exec(`UPDATE users SET name = ? WHERE id = ?;`, userName, userId)
	if err != nil {
		return fmt.Errorf("error updating user name: %s", err)
	}
	count, err := result.RowsAffected()
	if err != nil || count != 1 {
		return fmt.Errorf("error estimating affected rows: %s", err)
	}
	return nil
}

func UpdateUserStatus(dbConn sqlx.Ext, userId uint, userStatus string) error {
	if userId <= 0 {
		return fmt.Errorf("invalid user id[%d]", userId)
	}

	result, err := dbConn.Exec(`UPDATE users SET status = ? WHERE id = ?;`, userStatus, userId)
	if err != nil {
		return fmt.Errorf("error updating user status: %s", err)
	}
	count, err := result.RowsAffected()
	if err != nil || count != 1 {
		return fmt.Errorf("error estimating affected rows: %s", err)
	}
	return nil
}
