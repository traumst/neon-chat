package handler

import (
	"fmt"
	"log"
	d "prplchat/src/db"
	"prplchat/src/handler/state"
	a "prplchat/src/model/app"
	t "prplchat/src/model/template"
	"prplchat/src/utils"
	"time"
)

func IssueReservationToken(state *state.State, db *d.DBConn, user *a.User) (*t.VerifyEmailTemplate, error) {
	token := utils.RandStringBytes(16)
	expire := time.Now().Add(1 * time.Hour)
	reserve, err := reserve(db, user, token, expire)
	if err != nil {
		log.Printf("IssueReservationToken ERROR reserving[%s], %s\n", user.Email, err.Error())
		return nil, fmt.Errorf("")
	}
	emailConfig := state.SmtpConfig()
	tmpl := t.VerifyEmailTemplate{
		SourceEmail: emailConfig.User,
		UserEmail:   user.Email,
		UserName:    user.Name,
		Token:       reserve.Token,
		//TokenExpire: reserve.Expire.Format(time.RFC3339),
		TokenExpire: reserve.Expire.Format(time.Stamp),
	}
	err = sendSignupCompletionEmail(tmpl, emailConfig.User, emailConfig.Pass)
	if err != nil {
		log.Printf("IssueReservationToken ERROR sending email from [%s] to [%s], %s\n",
			emailConfig.User, user.Email, err.Error())
		return nil, fmt.Errorf("failed to send email to[%s]", user.Email)
	}
	return &tmpl, nil
}

func reserve(db *d.DBConn, user *a.User, token string, expire time.Time) (*d.Reservation, error) {
	reserve := &d.Reservation{
		Id:     0,
		UserId: user.Id,
		Token:  token,
		Expire: expire,
	}
	reserve, err := db.AddReservation(*reserve)
	if err != nil {
		return nil, fmt.Errorf("reserve[%s] for user[%d], %s", token, user.Id, err)
	} else if reserve == nil {
		return nil, fmt.Errorf("reserve[%s] for user[%d] got reserve NIL", token, user.Id)
	} else if reserve.Id <= 0 {
		return nil, fmt.Errorf("reserve[%s] for user[%d] got reserve id 0", token, user.Id)
	}
	return reserve, nil
}
