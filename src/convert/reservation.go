package convert

import (
	"neon-chat/src/app"
	"neon-chat/src/db"
)

func ReservationAppToDB(res *app.Reservation) *db.Reservation {
	return &db.Reservation{
		Id:     res.Id,
		UserId: res.UserId,
		Token:  res.Token,
		Expire: res.Expire,
	}
}

func ReservationDBToApp(res *db.Reservation) *app.Reservation {
	return &app.Reservation{
		Id:     res.Id,
		UserId: res.UserId,
		Token:  res.Token,
		Expire: res.Expire,
	}
}
