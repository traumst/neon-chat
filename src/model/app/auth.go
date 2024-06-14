package app

type AuthType string

const (
	AuthTypeEmail  AuthType = "email"
	AuthTypeGoogle AuthType = "google"
)

type Auth struct {
	Id     uint     `db:"id"`
	UserId uint     `db:"user_id"`
	Type   AuthType `db:"type"`
	Hash   string   `db:"hash"`
}
