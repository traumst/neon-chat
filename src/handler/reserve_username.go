package handler

import (
	"fmt"
	"log"
	"time"

	"neon-chat/src/db"
	"neon-chat/src/handler/email"
	"neon-chat/src/model/app"
	"neon-chat/src/model/template"
	"neon-chat/src/utils"
)

func ReserveUserName(dbConn *db.DBConn, emailConfig *utils.SmtpConfig, user *app.User) (template.VerifyEmailTemplate, error) {
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
		return template.VerifyEmailTemplate{},
			fmt.Errorf("failed to reserve[%s] for user[%d] got err, %s", token, user.Id, err.Error())
	} else if reserve == nil {
		return template.VerifyEmailTemplate{},
			fmt.Errorf("failed to reserve[%s] for user[%d] got reserve NIL", token, user.Id)
	} else if reserve.Id <= 0 {
		return template.VerifyEmailTemplate{},
			fmt.Errorf("failed to reserve[%s] for user[%d] got reserve id 0", token, user.Id)
	}
	tmpl := template.VerifyEmailTemplate{
		SourceEmail: emailConfig.User,
		UserEmail:   user.Email,
		UserName:    user.Name,
		Token:       reserve.Token,
		//TokenExpire: reserve.Expire.Format(time.RFC3339),
		TokenExpire: reserve.Expire.Format(time.Stamp),
	}
	err = email.SendSignupCompletionEmail(tmpl, emailConfig.User, emailConfig.Pass)
	if err != nil {
		log.Printf("IssueReservationToken ERROR sending email from [%s] to [%s], %s\n",
			emailConfig.User, user.Email, err.Error())
		return template.VerifyEmailTemplate{}, fmt.Errorf("failed to send email to[%s]", user.Email)
	}
	return tmpl, nil
}
