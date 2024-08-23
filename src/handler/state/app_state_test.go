package state

import (
	"math/rand"
	"net/http/httptest"
	"testing"

	"prplchat/src/model/app"
	"prplchat/src/utils"
	h "prplchat/src/utils/http"
)

func TestStateDefaults(t *testing.T) {
	app1 := &GlobalAppState
	app1.Init(utils.Config{})
	if app1.isInit != true {
		t.Errorf("TestStateDefaults expected isInit true, got [%v]", app1.isInit)
	}
	if app1.config.CacheSize != 1024 {
		t.Errorf("TestStateDefaults expected cache size 1024, got [%d]", app1.config.CacheSize)
	}
}

func TestAddConn(t *testing.T) {
	t.Logf("TestAddConn started")
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/some-route", nil)
	user := app.User{Id: 1, Name: "John", Type: app.UserType(app.UserTypeBasic)}
	app1 := &GlobalAppState
	app1.Init(utils.Config{})
	conn1 := app1.AddConn(w, *r, &user, nil)
	if conn1 == nil {
		t.Fatalf("TestAddConn expected a conn1, got nil")
	}
	conn2 := app1.AddConn(w, *r, &user, nil)
	if conn2 == nil {
		t.Fatalf("TestAddConn expected a conn2, got nil")
	}
	if conn1 == conn2 {
		t.Fatalf("TestAddConn expected conn1 and conn2 to be different")
	}
}

func TestGetConn(t *testing.T) {
	t.Logf("TestGetConn started")
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/some-route", nil)
	reqId := "test-req-id"
	reqId = h.SetReqId(r, &reqId)
	user := app.User{Id: uint(rand.Uint32()), Name: "John", Type: app.UserType(app.UserTypeBasic)}
	app1 := &GlobalAppState
	app1.Init(utils.Config{})
	expect := app1.AddConn(w, *r, &user, nil)
	if expect == nil {
		t.Fatalf("TestGetConn expected a conn, got nil")
	}
	conns := app1.GetConn(user.Id)
	if len(conns) == 0 {
		t.Fatalf("TestGetConn expected connections were empty, count [%d]", len(conns))
	}
	var conn *Conn
	for _, c := range conns {
		if c.Id == expect.Id {
			conn = c
			break
		}
	}
	if expect == nil || conn == nil {
		t.Fatalf("TestGetConn expected connections, got\n expected[%v]\n conn[%v]", expect, conn)
	}
	if expect.User != conn.User {
		t.Fatalf("TestGetConn expected equality,\nexpected user[%v],\nconn origin[%v]", expect.User, conn.User)
	}
	if expect.Origin != conn.Origin {
		t.Fatalf("TestGetConn expected equality,\nexpected origin[%v],\nconn origin[%v]", expect.Origin, conn.Origin)
	}
	if expect.In != conn.In {
		t.Fatalf("TestGetConn expected equality,\nexpected in[%v],\nconn in[%v]", expect.In, conn.In)
	}
	if expect.Writer != conn.Writer {
		t.Fatalf("TestGetConn expected equality,\nexpected writer[%v],\nconn writer[%v]", expect.Writer, conn.Writer)
	}
}

func TestDropConn(t *testing.T) {
	t.Logf("TestDropConn started")
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/some-route", nil)
	reqId := "test-req-id"
	reqId = h.SetReqId(r, &reqId)
	user := app.User{Id: uint(rand.Uint32()), Name: "John", Type: app.UserType(app.UserTypeBasic)}
	app1 := &GlobalAppState
	app1.Init(utils.Config{})
	expect := app1.AddConn(w, *r, &user, nil)
	if expect == nil {
		t.Fatalf("TestDropConn expected a conn, got nil")
	}
	conns := app1.GetConn(user.Id)
	if len(conns) == 0 {
		t.Fatalf("TestDropConn expected connections were empty, count [%d]", len(conns))
	}
	err := app1.DropConn(expect)
	if err != nil {
		t.Fatalf("TestDropConn expected no error, got [%s]", err.Error())
	}
	conns = app1.GetConn(user.Id)
	if len(conns) != 0 {
		t.Fatalf("TestGetConn user still has [%d] connections[%v]", len(conns), conns)
	}
}

func TestOpenChat(t *testing.T) {
	t.Logf("TestOpenChat started")
	user := app.User{
		Id: uint(rand.Uint32()),
		// Name: "John",
		// Type: app.UserType(app.UserTypeBasic),
	}
	app1 := &GlobalAppState
	app1.Init(utils.Config{
		Port:   0,
		Sqlite: "",
		Smtp: utils.SmtpConfig{
			User: "",
			Pass: "",
			Host: "",
			Port: "",
		},
		CacheSize: 0,
	})
	err := app1.OpenChat(user.Id, 22)
	if err != nil {
		t.Fatalf("TestOpenChat expected no error, %s", err.Error())
	}
	openChatId := app1.GetOpenChat(user.Id)
	if openChatId != 22 {
		t.Fatalf("TestOpenChat expected chat[22], got [%d]", openChatId)
	}
}

func TestGetOpenChatEmpty(t *testing.T) {
	t.Logf("TestGetOpenChatEmpty started")
	user := app.User{
		Id: uint(rand.Uint32()),
		// Name: "John",
		// Type: app.UserType(app.UserTypeBasic),
	}
	app1 := &GlobalAppState
	app1.Init(utils.Config{})
	open := app1.GetOpenChat(user.Id)
	if open != 0 {
		t.Fatalf("TestGetOpenChatEmpty expected 0, got [%v]", open)
	}
}

func TestGetOpenChat(t *testing.T) {
	t.Logf("TestGetOpenChat started")
	user := app.User{
		Id: uint(rand.Uint32()),
		// Name: "John",
		// Type: app.UserType(app.UserTypeBasic),
	}
	app1 := &GlobalAppState
	app1.Init(utils.Config{})
	err := app1.OpenChat(user.Id, 33)
	if err != nil {
		t.Fatalf("TestGetOpenChat failed to open chat, %s", err.Error())
	}
	open3 := app1.GetOpenChat(user.Id)
	if open3 != 33 {
		t.Fatalf("TestGetOpenChat expected open chat[0], got [%d]", open3)
	}
	err = app1.OpenChat(user.Id, 22)
	if err != nil {
		t.Fatalf("TestGetOpenChat expected no error, %s", err.Error())
	}
	open := app1.GetOpenChat(user.Id)
	if open != 22 {
		t.Fatalf("TestGetOpenChat expected chat[%d], got [%d]", 22, open)
	}
}

func TestCloseChat(t *testing.T) {
	t.Logf("TestCloseChat started")
	user := app.User{
		Id: uint(rand.Uint32()),
		// Name: "John",
		// Type: app.UserType(app.UserTypeBasic),
	}
	app1 := &GlobalAppState
	app1.Init(utils.Config{})
	app1.OpenChat(user.Id, 11)
	open := app1.GetOpenChat(user.Id)
	if open != 11 {
		t.Fatalf("TestCloseChat expected chat[%d], got [%d]", 11, open)
	}
	err := app1.CloseChat(user.Id, 11)
	if err != nil {
		t.Fatalf("TestCloseChat failed to close chat, %s", err.Error())
	}
	open = app1.GetOpenChat(user.Id)
	if open != 0 {
		t.Fatalf("TestCloseChat expected chat[%d] to be closed, but got[%v]", 11, open)
	}
}
