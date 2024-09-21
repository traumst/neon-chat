package app

import "neon-chat/src/app/enum"

type Auth struct {
	Id     uint          `db:"id"`
	UserId uint          `db:"user_id"`
	Type   enum.AuthType `db:"type"`
	Hash   string        `db:"hash"`
}
