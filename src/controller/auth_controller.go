package controller

import (
	"fmt"
	"log"
	"net/http"
	"time"

	d "go.chat/src/db"
	"go.chat/src/handler"
	a "go.chat/src/model/app"
	"go.chat/src/model/template"
	"go.chat/src/utils"
	h "go.chat/src/utils/http"
)

// TODO support other types
const (
	LocalUserType = a.UserTypeLocal
	LocalAuthType = a.AuthTypeLocal
)

func Login(app *handler.AppState, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] Login TRACE IN\n", h.GetReqId(r))
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("action not allowed"))
		return
	}
	loginUser := utils.Trim(r.FormValue("login-user"))
	loginPass := utils.Trim(r.FormValue("login-pass"))
	if loginUser == "" || loginPass == "" || len(loginUser) < 4 || len(loginPass) < 4 {
		log.Printf("[%s] Login TRACE empty user[%s]", h.GetReqId(r), loginUser)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad login credentials"))
		return
	}
	log.Printf("[%s] Login TRACE authentication check for user[%s] auth[%s]\n",
		h.GetReqId(r), loginUser, LocalAuthType)
	user, auth, err := handler.Authenticate(db, loginUser, loginPass, LocalAuthType)
	if err != nil {
		log.Printf("[%s] Login ERROR unauth user[%s], %s\n", h.GetReqId(r), loginUser, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("authentication failed"))
		return
	}
	if user == nil {
		log.Printf("[%s] Login ERROR unknown user[%s]\n", h.GetReqId(r), loginUser)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unknown user"))
		return
	} else if user.Status != a.UserStatusActive {
		log.Printf("[%s] Login ERROR inactive user[%d] status[%s]\n", h.GetReqId(r), user.Id, user.Status)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("inactive user"))
		return
	}
	if auth == nil {
		log.Printf("[%s] Login ERROR user password mismatched [%s]\n", h.GetReqId(r), loginUser)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user password mismatched"))
		return
	}
	cookie := h.SetSessionCookie(w, user, auth)
	log.Printf("[%s] Login TRACE user[%d] authenticated until [%s]\n",
		h.GetReqId(r), user.Id, cookie.Expire.Format(time.RFC1123Z))
	http.Header.Add(w.Header(), "HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func Logout(app *handler.AppState, w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] Logout TRACE \n", h.GetReqId(r))
	user, err := handler.ReadSession(app, w, r)
	if user == nil {
		log.Printf("[%s] Logout INFO user is not authorized, %s\n", h.GetReqId(r), err.Error())
		RenderLogin(w, r, nil)
		return
	}
	h.ClearSessionCookie(w, user.Id)
	app.UntrackUser(user.Id)
	http.Header.Add(w.Header(), "HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func SignUp(app *handler.AppState, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] SignUp TRACE IN\n", h.GetReqId(r))
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("This shouldn't happen"))
		return
	}
	signupUser := utils.Trim(r.FormValue("signup-user"))
	signupEmail := utils.Trim(r.FormValue("signup-email"))
	signupPass := utils.Trim(r.FormValue("signup-pass"))
	log.Printf("[%s] SignUp TRACE authentication check for user[%s] auth[%s]\n",
		h.GetReqId(r), signupUser, LocalAuthType)
	if signupUser == "" || signupEmail == "" || signupPass == "" ||
		len(signupUser) < 4 || len(signupEmail) < 4 || len(signupPass) < 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad signup credentials"))
		return
	}
	user, auth, _ := handler.Authenticate(db, signupUser, signupPass, LocalAuthType)
	if user != nil && user.Status == a.UserStatusActive && auth != nil {
		log.Printf("[%s] SignUp TRACE signedIn instead of signUp user[%s]\n", h.GetReqId(r), signupUser)
		h.SetSessionCookie(w, user, auth)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if user != nil {
		log.Printf("[%s] SignUp ERROR there is already name[%s] taken by user[%d] in status[%s]\n",
			h.GetReqId(r), signupUser, user.Id, user.Status)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("username is already taken"))
		return
	}
	if !handler.IsEmailValid(signupEmail) {
		log.Printf("[%s] SignUp ERROR invalid email[%s]\n", h.GetReqId(r), signupEmail)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid email"))
		return
	}
	log.Printf("[%s] SignUp TRACE register new user[%s]\n", h.GetReqId(r), signupUser)
	salt := utils.GenerateSalt(signupUser, string(LocalUserType))
	user = &a.User{
		Id:     0,
		Name:   signupUser,
		Email:  signupEmail,
		Type:   LocalUserType,
		Status: a.UserStatusPending,
		Salt:   salt,
	}
	user, auth, err := handler.Register(db, user, signupPass, LocalAuthType)
	if err != nil {
		log.Printf("[%s] SignUp ERROR on register user[%s][%s], %s\n", h.GetReqId(r), signupUser, signupEmail, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Failed to register user [%s:%s]", LocalUserType, signupUser)))
		return
	} else if user == nil || auth == nil {
		log.Printf("[%s] SignUp ERROR to register user[%v]\n", h.GetReqId(r), user)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Failed to register user [%s:%s]", LocalUserType, signupUser)))
		return
	}
	defer func() {
		// TODO delete user and auth
		// if r := recover(); r != nil {
		// 	handler.DeleteUser()
		// }
	}()
	log.Printf("[%s] SignUp TRACE issuing reservation to [%s]\n", h.GetReqId(r), user.Email)
	sentEmail, err := handler.IssueReservationToken(app, db, user)
	if err != nil {
		log.Printf("[%s] SignUp ERROR failed to issue reservation token to email[%s], %s\n",
			h.GetReqId(r), user.Email, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to issue reservation token"))
		return
	}
	defer func() {
		// TODO delete reservation token
		// if r := recover(); r != nil {
		// 	handler.DeleteReservation()
		// }
	}()
	html, err := sentEmail.HTML()
	if err != nil {
		log.Printf("[%s] SignUp ERROR templating result html[%v], %s\n", h.GetReqId(r), sentEmail, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to template response"))
		return
	}
	log.Printf("[%s] SignUp TRACE OUT\n", h.GetReqId(r))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func ConfirmEmail(app *handler.AppState, db *d.DBConn, w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] ConfirmEmail TRACE IN\n", h.GetReqId(r))
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("This shouldn't happen"))
		return
	}
	signupToken := r.URL.Query().Get("token")
	if signupToken == "" {
		log.Printf("[%s] ConfirmEmail ERROR missing token\n", h.GetReqId(r))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing token"))
		return
	}

	reserve, err := db.GetReservation(signupToken)
	if err != nil {
		log.Printf("[%s] ConfirmEmail ERROR error reading reservation, %s\n", h.GetReqId(r), err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("missing token"))
		return
	} else if reserve == nil {
		log.Printf("[%s] ConfirmEmail WARN reservation not found\n", h.GetReqId(r))
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	} else if reserve.Expire.Before(time.Now()) {
		log.Printf("[%s] ConfirmEmail WARN reservation[%d] expired\n", h.GetReqId(r), reserve.Id)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("token expired"))
		return
	} else if reserve.UserId <= 0 {
		log.Printf("[%s] ConfirmEmail WARN reservation[%d] corrupted, userId[%d]\n", h.GetReqId(r), reserve.Id, reserve.UserId)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("corrupted token"))
		return
	}

	user, err := app.GetUser(reserve.UserId)
	if err != nil {
		dbUser, err := db.GetUser(reserve.UserId)
		if err != nil {
			log.Printf("[%s] ConfirmEmail ERROR retrieving user[%d], %s\n", h.GetReqId(r), reserve.UserId, err.Error())
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("corrupted token"))
			return
		}
		tmp := handler.UserFromDB(*dbUser)
		user = &tmp
		err = app.TrackUser(user)
		if err != nil {
			log.Printf("[%s] ConfirmEmail ERROR tracking user[%d], %s\n", h.GetReqId(r), user.Id, err.Error())
		}
	} else if user == nil {
		log.Printf("[%s] ConfirmEmail ERROR user[%d] not found\n", h.GetReqId(r), reserve.UserId)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}
	if user.Status != a.UserStatusPending {
		log.Printf("[%s] ConfirmEmail ERROR user[%d] status[%s] is not pending\n", h.GetReqId(r), user.Id, user.Status)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid user status"))
		return
	}

	err = db.UpdateUserStatus(user.Id, string(a.UserStatusActive))
	if err != nil {
		log.Printf("[%s] ConfirmEmail ERROR failed to update user[%d] status\n", h.GetReqId(r), user.Id)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to update user status"))
		return
	}

	RenderLogin(w, r, &template.InfoMessage{
		Header: "Congrats! " + user.Email + " is confirmed",
		Body:   "Your user name is " + user.Name + " until you decide to change it",
		Footer: "Please, login using your signup credentials",
	})
}
