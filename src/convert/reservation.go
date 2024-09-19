package convert

import (
	"neon-chat/src/db"
	"neon-chat/src/model/app"
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
