package app

type AuthType string

const (
	AuthTypeUnknown AuthType = "unknown"
	AuthTypeLocal   AuthType = "local"
	AuthTypePass    AuthType = "pass"
	AuthTypeSSO     AuthType = "sso"
)

type Auth struct {
	Id     uint
	UserId uint
	Type   AuthType
	Hash   uint
}
