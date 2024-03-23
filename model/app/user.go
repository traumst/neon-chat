package app

type UserType string

// TODO add flags/permissions mapping

const (
	Admin UserType = "admin"
	Free  UserType = "free"
)

type User struct {
	Id   uint     `db:"id"`
	Name string   `db:"name"`
	Type UserType `db:"type"`
	Salt []byte   `db:"salt"`
}
