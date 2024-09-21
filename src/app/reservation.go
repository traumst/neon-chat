package app

import "time"

type Reservation struct {
	Id     uint
	UserId uint
	Token  string
	Expire time.Time
}
