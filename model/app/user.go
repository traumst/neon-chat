package app

type UserType string

// TODO add flags/permissions mapping

const (
	UserTypeAdmin UserType = "user-type-admin"
	UserTypeFree  UserType = "user-type-free"
)

type User struct {
	Id   uint     `db:"id"`
	Name string   `db:"name"`
	Type UserType `db:"type"`
	Salt string   `db:"salt"`
}
