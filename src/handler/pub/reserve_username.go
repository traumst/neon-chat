package pub

import (
	"fmt"
	"time"

	"neon-chat/src/convert"
	"neon-chat/src/db"
	"neon-chat/src/model/app"
	"neon-chat/src/utils"
)

func ReserveUserName(
	dbConn *db.DBConn,
	emailConfig *utils.SmtpConfig,
	user *app.User,
) (*app.Reservation, error) {
	token := utils.RandStringBytes(16)
	expire := time.Now().Add(1 * time.Hour)
	reserve := &db.Reservation{
		Id:     0,
		UserId: user.Id,
		Token:  token,
		Expire: expire,
	}
	reserve, err := dbConn.AddReservation(*reserve)
	if err != nil {
		return nil, fmt.Errorf("failed to reserve[%s] for user[%d] got err, %s", token, user.Id, err.Error())
	} else if reserve == nil {
		return nil, fmt.Errorf("failed to reserve[%s] for user[%d] got reserve NIL", token, user.Id)
	} else if reserve.Id <= 0 {
		return nil, fmt.Errorf("failed to reserve[%s] for user[%d] got reserve id 0", token, user.Id)
	}
	return convert.ReservationDBToApp(reserve), nil
}
