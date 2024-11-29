package state

import (
	"fmt"
	"math/rand"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"neon-chat/src/app"
	"neon-chat/src/app/enum"
	"neon-chat/src/utils"
	"neon-chat/src/utils/config"
	h "neon-chat/src/utils/http"
)

func TestStateDefaults(t *testing.T) {
	app1 := &GlobalAppState
	app1.Init(config.Config{})
	if app1.isInit != true {
		t.Errorf("expected isInit true, got [%v]", app1.isInit)
	}
	if app1.config.CacheSize != 1024 {
		t.Errorf("expected cache size 1024, got [%d]", app1.config.CacheSize)
	}
}

func TestAddConn(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/some-route", nil)
	user := app.User{Id: 1, Name: "John", Type: enum.UserType(enum.UserTypeBasic)}
	app1 := &GlobalAppState
	app1.Init(config.Config{})
	conn1 := app1.AddConn(w, *r, &user, nil)
	if conn1 == nil {
		t.Fatalf("expected a conn1, got nil")
	}
	conn2 := app1.AddConn(w, *r, &user, nil)
	if conn2 == nil {
		t.Fatalf("expected a conn2, got nil")
	}
	if conn1 == conn2 {
		t.Fatalf("expected conn1 and conn2 to be different")
	}
}

func TestGetConn(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/some-route", nil)
	reqId := "test-req-id"
	reqId = h.SetReqId(r, &reqId)
	user := app.User{Id: uint(rand.Uint32()), Name: "John", Type: enum.UserType(enum.UserTypeBasic)}
	app1 := &GlobalAppState
	app1.Init(config.Config{})
	expect := app1.AddConn(w, *r, &user, nil)
	if expect == nil {
		t.Fatalf("expected a conn, got nil")
	}
	conns := app1.GetConn(user.Id)
	if len(conns) == 0 {
		t.Fatalf("expected connections were empty, count [%d]", len(conns))
	}
	var conn *Conn
	for _, c := range conns {
		if c.Id == expect.Id {
			conn = c
			break
		}
	}
	if expect == nil || conn == nil {
		t.Fatalf("expected connections, got\n expected[%v]\n conn[%v]", expect, conn)
	}
	if expect.User != conn.User {
		t.Fatalf("expected equality,\nexpected user[%v],\nconn origin[%v]", expect.User, conn.User)
	}
	if expect.Origin != conn.Origin {
		t.Fatalf("expected equality,\nexpected origin[%v],\nconn origin[%v]", expect.Origin, conn.Origin)
	}
	if expect.In != conn.In {
		t.Fatalf("expected equality,\nexpected in[%v],\nconn in[%v]", expect.In, conn.In)
	}
	if expect.Writer != conn.Writer {
		t.Fatalf("expected equality,\nexpected writer[%v],\nconn writer[%v]", expect.Writer, conn.Writer)
	}
}

func TestDropConn(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/some-route", nil)
	reqId := "test-req-id"
	reqId = h.SetReqId(r, &reqId)
	user := app.User{Id: uint(rand.Uint32()), Name: "John", Type: enum.UserType(enum.UserTypeBasic)}
	app1 := &GlobalAppState
	app1.Init(config.Config{})
	expect := app1.AddConn(w, *r, &user, nil)
	if expect == nil {
		t.Fatalf("expected a conn, got nil")
	}
	conns := app1.GetConn(user.Id)
	if len(conns) == 0 {
		t.Fatalf("expected connections were empty, count [%d]", len(conns))
	}
	err := app1.DropConn(expect)
	if err != nil {
		t.Fatalf("expected no error, got [%s]", err.Error())
	}
	conns = app1.GetConn(user.Id)
	if len(conns) != 0 {
		t.Fatalf("user still has [%d] connections[%v]", len(conns), conns)
	}
}

func TestOpenChat(t *testing.T) {
	user := app.User{
		Id: uint(rand.Uint32()),
		// Name: "John",
		// Type: enum.UserType(enum.UserTypeBasic),
	}
	app1 := &GlobalAppState
	app1.Init(config.Config{
		Port:      0,
		Sqlite:    "",
		Smtp:      config.SmtpConfig{User: "", Pass: "", Host: "", Port: ""},
		CacheSize: 0,
		TestUsers: []*config.TestUser{},
	})
	err := app1.OpenChat(user.Id, 22)
	if err != nil {
		t.Fatalf("expected no error, %s", err.Error())
	}
	openChatId := app1.GetOpenChat(user.Id)
	if openChatId != 22 {
		t.Fatalf("expected chat[22], got [%d]", openChatId)
	}
}

func TestGetOpenChatEmpty(t *testing.T) {
	user := app.User{
		Id: uint(rand.Uint32()),
		// Name: "John",
		// Type: enum.UserType(enum.UserTypeBasic),
	}
	app1 := &GlobalAppState
	app1.Init(config.Config{})
	open := app1.GetOpenChat(user.Id)
	if open != 0 {
		t.Fatalf("expected 0, got [%v]", open)
	}
}

func TestGetOpenChat(t *testing.T) {
	user := app.User{
		Id: uint(rand.Uint32()),
		// Name: "John",
		// Type: enum.UserType(enum.UserTypeBasic),
	}
	app1 := &GlobalAppState
	app1.Init(config.Config{})
	err := app1.OpenChat(user.Id, 33)
	if err != nil {
		t.Fatalf("failed to open chat, %s", err.Error())
	}
	open3 := app1.GetOpenChat(user.Id)
	if open3 != 33 {
		t.Fatalf("expected open chat[0], got [%d]", open3)
	}
	err = app1.OpenChat(user.Id, 22)
	if err != nil {
		t.Fatalf("expected no error, %s", err.Error())
	}
	open := app1.GetOpenChat(user.Id)
	if open != 22 {
		t.Fatalf("expected chat[%d], got [%d]", 22, open)
	}
}

func TestCloseChat(t *testing.T) {
	t.Logf("started")
	user := app.User{
		Id: uint(rand.Uint32()),
		// Name: "John",
		// Type: enum.UserType(enum.UserTypeBasic),
	}
	app1 := &GlobalAppState
	app1.Init(config.Config{})
	app1.OpenChat(user.Id, 11)
	open := app1.GetOpenChat(user.Id)
	if open != 11 {
		t.Fatalf("expected chat[%d], got [%d]", 11, open)
	}
	err := app1.CloseChat(user.Id, 11)
	if err != nil {
		t.Fatalf("failed to close chat, %s", err.Error())
	}
	open = app1.GetOpenChat(user.Id)
	if open != 0 {
		t.Fatalf("expected chat[%d] to be closed, but got[%v]", 11, open)
	}
}

func TestSaveToFile(t *testing.T) {
	t.Logf("started")
	testFile := fmt.Sprintf("test-save-state-%d.json", rand.Uint32())
	defer func() {
		_ = os.Remove(testFile)
	}()

	app1 := &GlobalAppState
	app1.Init(config.Config{CacheSize: 8})
	_ = app1.OpenChat(1, 11)
	_ = app1.OpenChat(2, 22)
	_ = app1.OpenChat(3, 33)
	err := app1.SaveToFile(testFile)
	if err != nil {
		t.Fatalf("expected no error, got [%s]", err.Error())
	}
	if app1.chats.Count() != 0 {
		t.Fatalf("expected to empty chats, got [%d]", app1.chats.Count())
	}
	if openChat := app1.GetOpenChat(1); openChat != 0 {
		t.Fatalf("expected no open chat, got [%d]", openChat)
	}

	content, err := utils.ReadFileContent(testFile)
	if err != nil {
		t.Fatalf("failed to read file, %s", err.Error())
	}

	if !strings.Contains(content, "\"1\":11") {
		t.Fatalf("content missing '1':11, got [%s]", content)
	}
	if !strings.Contains(content, "\"2\":22") {
		t.Fatalf("content missing '2':22, got [%s]", content)
	}
	if !strings.Contains(content, "\"3\":33") {
		t.Fatalf("content missing '3':33, got [%s]", content)
	}
}

func TestLoadFromFile(t *testing.T) {
	t.Logf("TestLoadFromFile started")
	testFile := fmt.Sprintf("test-load-state-%d.json", rand.Uint32())
	defer func() {
		_ = os.Remove(testFile)
	}()

	app1 := &GlobalAppState
	app1.Init(config.Config{CacheSize: 8})
	_ = app1.OpenChat(1, 11)
	_ = app1.OpenChat(2, 22)
	_ = app1.OpenChat(3, 33)
	_ = app1.SaveToFile(testFile)
	count := app1.chats.Count()
	if count != 0 {
		t.Fatalf("expected to empty chats, got [%d]", count)
	}
	if openChat := app1.GetOpenChat(1); openChat != 0 {
		t.Fatalf("expected no open chat, got [%d]", openChat)
	}

	err := app1.LoadFromFile(testFile)
	if err != nil {
		t.Fatalf("failed to load file, %s", err.Error())
	}
	count = app1.chats.Count()
	if count != 3 {
		t.Fatalf("expected chat length to be 3 got [%d]", count)
	}
	if openChat := app1.GetOpenChat(1); openChat != 11 {
		t.Fatalf("expected open chat 11, got [%d]", openChat)
	}
	if openChat := app1.GetOpenChat(2); openChat != 22 {
		t.Fatalf("expected open chat 22, got [%d]", openChat)
	}
	if openChat := app1.GetOpenChat(3); openChat != 33 {
		t.Fatalf("expected open chat 33, got [%d]", openChat)
	}
}
