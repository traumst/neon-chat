package app

type AuthType string

const (
	AuthTypeUnknown AuthType = "unknown"
	AuthTypeLocal   AuthType = "local"
	AuthTypePass    AuthType = "pass"
	AuthTypeSSO     AuthType = "sso"
)

type Auth struct {
	Id     uint     `db:"id"`
	UserId uint     `db:"user_id"`
	Type   AuthType `db:"type"`
	Hash   string   `db:"hash"`
}
