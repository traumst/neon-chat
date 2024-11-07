package controller

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"neon-chat/src/app"
	"neon-chat/src/app/enum"
	"neon-chat/src/consts"
	"neon-chat/src/convert"
	"neon-chat/src/db"
	"neon-chat/src/handler/email"
	"neon-chat/src/handler/pub"
	"neon-chat/src/state"
	"neon-chat/src/template"
	"neon-chat/src/utils"
	h "neon-chat/src/utils/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	log.Printf("TRACE [%s] IN\n", reqId)
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("action not allowed"))
		return
	}
	loginUser := utils.SanitizeInput(r.FormValue("login-user"))
	loginPass := utils.SanitizeInput(r.FormValue("login-pass"))
	if len(loginUser) < 4 || len(loginPass) < 4 {
		log.Printf("TRACE [%s] empty user[%s]", reqId, loginUser)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad login credentials"))
		return
	}
	log.Printf("TRACE [%s] authentication check for user[%s] auth[%s]\n",
		reqId, loginUser, enum.AuthTypeEmail)
	dbConn := r.Context().Value(consts.DBConn).(*db.DBConn)
	user, auth, err := pub.AuthenticateUser(dbConn, loginUser, loginPass, enum.AuthTypeEmail)
	if err != nil {
		log.Printf("ERROR [%s] unauth user[%s], %s\n", reqId, loginUser, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("authentication failed"))
		return
	}
	if user == nil {
		log.Printf("ERROR [%s] unknown user[%s]\n", reqId, loginUser)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unknown user"))
		return
	} else if user.Status != enum.UserStatusActive {
		log.Printf("ERROR [%s] inactive user[%d] status[%s]\n", reqId, user.Id, user.Status)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("inactive user"))
		return
	}
	if auth == nil {
		log.Printf("ERROR [%s] user password mismatched [%s]\n", reqId, loginUser)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user password mismatched"))
		return
	}
	cookie := h.SetSessionCookie(w, user, auth)
	log.Printf("TRACE [%s] user[%d] authenticated until [%s]\n",
		reqId, user.Id, cookie.Expire.Format(time.RFC1123Z))
	http.Header.Add(w.Header(), "HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	log.Printf("TRACE [%s] \n", reqId)
	user := r.Context().Value(consts.ActiveUser).(*app.User)
	if user == nil {
		log.Printf("INFO [%s] user is unauthorized\n", reqId)
		http.Header.Add(w.Header(), "HX-Refresh", "true")
		return
	}
	h.ClearSessionCookie(w, user.Id)
	http.Header.Add(w.Header(), "HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	log.Printf("TRACE [%s] IN\n", reqId)
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("This shouldn't happen"))
		return
	}
	signupUser := utils.SanitizeInput(r.FormValue("signup-user"))
	signupEmail := utils.SanitizeInput(r.FormValue("signup-email"))
	signupPass := utils.SanitizeInput(r.FormValue("signup-pass"))
	log.Printf("TRACE [%s] authentication check for user[%s] auth[%s]\n",
		reqId, signupUser, enum.AuthTypeEmail)
	if len(signupUser) < 4 || len(signupEmail) < 4 || len(signupPass) < 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad signup credentials"))
		return
	}
	dbConn := r.Context().Value(consts.DBConn).(*db.DBConn)
	user, auth, _ := pub.AuthenticateUser(dbConn, signupUser, signupPass, enum.AuthTypeEmail)
	if user != nil && user.Status == enum.UserStatusActive && auth != nil {
		log.Printf("TRACE [%s] signedIn instead of signUp user[%s]\n", reqId, signupUser)
		h.SetSessionCookie(w, user, auth)
		http.Header.Add(w.Header(), "HX-Refresh", "true")
		w.WriteHeader(http.StatusOK)
		return
	}
	if user != nil {
		log.Printf("ERROR [%s] there is already name[%s] taken by user[%d] in status[%s]\n",
			reqId, signupUser, user.Id, user.Status)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("username is already taken"))
		return
	}
	if !email.IsEmailValid(signupEmail) {
		log.Printf("ERROR [%s] invalid email[%s]\n", reqId, signupEmail)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid email"))
		return
	}
	log.Printf("TRACE [%s] register new user[%s]\n", reqId, signupUser)
	salt := utils.GenerateSalt(signupUser, string(enum.AuthTypeEmail))
	user = &app.User{
		Id:     0,
		Name:   signupUser,
		Email:  signupEmail,
		Type:   enum.UserTypeBasic,
		Status: enum.UserStatusPending,
		Salt:   salt,
	}
	user, auth, err := pub.RegisterUser(dbConn, user, signupPass, enum.AuthTypeEmail)
	if err != nil {
		log.Printf("ERROR [%s] on register user[%s][%s], %s\n", reqId, signupUser, signupEmail, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to register"))
		return
	} else if user == nil || auth == nil {
		log.Printf("ERROR [%s] to register user[%v] - but no err\n", reqId, user)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to registref"))
		return
	}
	log.Printf("TRACE [%s] issuing reservation to [%s]\n", reqId, user.Email)
	emailConfig, err := r.Context().Value(consts.AppState).(*state.State).SmtpConfig()
	if err != nil {
		panic(fmt.Errorf("ERROR  getting smtp config, %s", err.Error()))
	}
	reservation, err := pub.ReserveUserName(dbConn, emailConfig, user)
	if err != nil {
		log.Printf("ERROR [%s] failed to reserve username[%s] for [%s], %s\n",
			reqId, user.Name, user.Email, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to issue reservation token"))
		return
	}
	emailTmpl := template.VerifyEmailTemplate{
		SourceEmail: emailConfig.User,
		UserEmail:   user.Email,
		UserName:    user.Name,
		Token:       reservation.Token,
		//TokenExpire: reservation.Expire.Format(time.RFC3339),
		TokenExpire: reservation.Expire.Format(time.Stamp),
	}
	err = email.SendSignupCompletionEmail(emailTmpl, emailConfig.User, emailConfig.Pass)
	if err != nil {
		log.Printf("ERROR [%s] sending email from [%s] to [%s], %s\n",
			reqId, emailConfig.User, user.Email, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to send email"))
		return
	}
	changesMade := r.Context().Value(consts.TxChangesKey).(*bool)
	*changesMade = true
	html, err := emailTmpl.HTML()
	if err != nil {
		log.Printf("ERROR [%s] templating result html[%v], %s\n",
			reqId, emailTmpl, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to template response"))
		return
	}
	log.Printf("TRACE [%s] OUT\n", reqId)
	w.(*h.StatefulWriter).IndicateChanges()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	reqId := r.Context().Value(consts.ReqIdKey).(string)
	log.Printf("TRACE [%s] IN\n", reqId)
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("This shouldn't happen"))
		return
	}
	signupToken := r.URL.Query().Get("token")
	if signupToken == "" {
		log.Printf("ERROR [%s] missing token\n", reqId)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing token"))
		return
	}
	dbConn := r.Context().Value(consts.DBConn).(*db.DBConn)
	reserve, err := db.GetReservation(dbConn.Tx, signupToken)
	if err != nil {
		log.Printf("ERROR [%s] error reading reservation, %s\n", reqId, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("missing token"))
		return
	} else if reserve == nil {
		log.Printf("WARN [%s] reservation not found\n", reqId)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	} else if reserve.Expire.Before(time.Now()) {
		log.Printf("WARN [%s] reservation[%d] expired\n", reqId, reserve.Id)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("token expired"))
		return
	} else if reserve.UserId <= 0 {
		log.Printf("WARN [%s] reservation[%d] corrupted, userId[%d]\n", reqId, reserve.Id, reserve.UserId)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("corrupted token"))
		return
	}
	dbUser, err := db.GetUser(dbConn.Tx, reserve.UserId)
	if err != nil {
		log.Printf("ERROR [%s] retrieving user[%d], %s\n", reqId, reserve.UserId, err.Error())
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("corrupted token"))
		return
	}
	user := convert.UserDBToApp(dbUser, nil)
	if user.Status != enum.UserStatusPending {
		log.Printf("ERROR [%s] user[%d] status[%s] is not pending\n", reqId, user.Id, user.Status)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid user status"))
		return
	}
	err = db.UpdateUserStatus(dbConn.Tx, user.Id, string(enum.UserStatusActive))
	if err != nil {
		log.Printf("ERROR [%s] failed to update user[%d] status\n", reqId, user.Id)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to update user status"))
		return
	}
	w.(*h.StatefulWriter).IndicateChanges()
	http.Header.Add(w.Header(), "HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}
