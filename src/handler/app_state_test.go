package handler

import (
	"math/rand"
	"net/http/httptest"
	"testing"

	"prplchat/src/model/app"
	h "prplchat/src/utils/http"
)

var app1 = &ApplicationState

// CONN
func TestAddConn(t *testing.T) {
	t.Logf("TestAddConn started")
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/some-route", nil)
	user := app.User{Id: 1, Name: "John", Type: app.UserType(app.UserTypeBasic)}
	conn1 := app1.AddConn(w, *r, &user)
	if conn1 == nil {
		t.Fatalf("TestAddConn expected a conn1, got nil")
	}
	conn2 := app1.AddConn(w, *r, &user)
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
	expect := app1.AddConn(w, *r, &user)
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
	expect := app1.AddConn(w, *r, &user)
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

// USER
func TestInviteUser(t *testing.T) {
	t.Logf("TestInviteUser started")
	user := app.User{
		Id:   uint(rand.Uint32()),
		Name: "John",
		Type: app.UserType(app.UserTypeBasic),
	}
	invitee := app.User{Id: uint(rand.Uint32()),
		Name: "Jane",
		Type: app.UserType(app.UserTypeBasic),
	}
	chatId := app1.AddChat(&user, "TestChat")
	err := app1.InviteUser(user.Id, chatId, &invitee)
	if err != nil {
		t.Fatalf("TestInviteUser failed to invite user[%d] into chat[%d], [%s]",
			user.Id, chatId, err.Error())
	}
	chat, err := app1.GetChat(invitee.Id, chatId)
	if err != nil {
		t.Fatalf("TestInviteUser failed to get chat[%d] for user[%d], [%s]",
			chatId, user.Id, err.Error())
	}
	if chat == nil || chat.Id != chatId {
		t.Fatalf("TestInviteUser expected a chat[%d] for user[%d] got nil", chatId, user.Id)
	}
}

func TestDropUser(t *testing.T) {
	t.Logf("TestDropUser started")
	user := app.User{
		Id:   uint(rand.Uint32()),
		Name: "John",
		Type: app.UserType(app.UserTypeBasic),
	}
	invitee := app.User{Id: uint(rand.Uint32()),
		Name: "Jane",
		Type: app.UserType(app.UserTypeBasic),
	}
	chatId := app1.AddChat(&user, "TestChat")
	err := app1.InviteUser(user.Id, chatId, &invitee)
	if err != nil {
		t.Fatalf("TestDropUser failed to invite, %s", err.Error())
	}
	chat, err := app1.GetChat(invitee.Id, chatId)
	if err != nil {
		t.Fatalf("TestDropUser failed to get chat, %s", err.Error())
	}
	if chat == nil {
		t.Fatalf("TestDropUser expected chat[%d], got nil", chatId)
	}
	if chat.Id != chatId {
		t.Fatalf("TestDropUser expected chat[%d], got [%d]", chatId, chat.Id)
	}
	err = app1.DropUser(user.Id, chatId, invitee.Id)
	if err != nil {
		t.Fatalf("TestDropUser expected no error, %s", err.Error())
	}
	chat, err = app1.GetChat(invitee.Id, chatId)
	if err == nil {
		t.Fatalf("TestDropUser expected error, got chat[%v]", chat)
	}
	if chat != nil {
		t.Fatalf("TestDropUser expected error, got chat[%d]", chat.Id)
	}
}

// CHAT
func TestAddChat(t *testing.T) {
	t.Logf("TestAddChat started")
	user := app.User{
		Id:   uint(rand.Uint32()),
		Name: "John",
		Type: app.UserType(app.UserTypeBasic),
	}
	chatId := app1.AddChat(&user, "TestChat")
	chat, err := app1.GetChat(user.Id, chatId)
	if err != nil {
		t.Fatalf("TestAddChat failed to get chat[%d] for user[%d], [%s]",
			chatId, user.Id, err.Error())
	}
	if chat == nil || chat.Id != chatId {
		t.Fatalf("TestAddChat expected a chat[%d] for user[%d] got nil", chatId, user.Id)
	}
}

func TestGetChats(t *testing.T) {
	t.Logf("TestGetChats started")
	user := app.User{
		Id:   uint(rand.Uint32()),
		Name: "John",
		Type: app.UserType(app.UserTypeBasic),
	}
	chatId1 := app1.AddChat(&user, "TestChat1")
	chatId2 := app1.AddChat(&user, "TestChat2")
	chats := app1.GetChats(user.Id)
	if len(chats) != 2 {
		t.Fatalf("TestGetChats expected 2 chats, got [%d]", len(chats))
	}
	if chats[0].Id != chatId1 || chats[1].Id != chatId2 {
		t.Fatalf("TestGetChats expected chat ids [%d, %d], got [%d, %d]",
			chatId1, chatId2, chats[0].Id, chats[1].Id)
	}
}

func TestGetChat(t *testing.T) {
	t.Logf("TestGetChats started")
	user := app.User{
		Id:   uint(rand.Uint32()),
		Name: "John",
		Type: app.UserType(app.UserTypeBasic),
	}
	chatId1 := app1.AddChat(&user, "TestChat1")
	chat1, err := app1.GetChat(user.Id, chatId1)
	if err != nil {
		t.Fatalf("TestGetChat failed to get chat[%d] for user[%d], [%s]",
			chatId1, user.Id, err.Error())
	}
	if chat1 == nil {
		t.Fatalf("TestGetChat expected a chat[%d] for user[%d] got nil", chatId1, user.Id)
	}
	chatId2 := app1.AddChat(&user, "TestChat2")
	chat2, err := app1.GetChat(user.Id, chatId2)
	if err != nil {
		t.Fatalf("TestGetChat failed to get chat[%d] for user[%d], %s",
			chatId1, user.Id, err.Error())
	}
	if chat2 == nil {
		t.Fatalf("TestGetChat expected a chat[%d] for user[%d] got nil", chatId1, user.Id)
	}
}

func TestOpenChat(t *testing.T) {
	t.Logf("TestOpenChat started")
	user := app.User{
		Id:   uint(rand.Uint32()),
		Name: "John",
		Type: app.UserType(app.UserTypeBasic),
	}
	chatId1 := app1.AddChat(&user, "TestChat1")
	chatId2 := app1.AddChat(&user, "TestChat2")
	chatId3 := app1.AddChat(&user, "TestChat3")
	open, err := app1.OpenChat(user.Id, chatId2)
	if err != nil {
		t.Fatalf("TestOpenChat expected no error, %s", err.Error())
	}
	if open == nil {
		t.Fatalf("TestOpenChat expected a chat, got nil")
	}
	if open.Id != chatId2 {
		t.Fatalf("TestOpenChat expected chat[%d], got [%d]", chatId2, open.Id)
	}
	if open.Id == chatId1 || open.Id == chatId3 {
		t.Fatalf("TestOpenChat expected open chat[%d] to be different from [%d, %d]",
			chatId2, chatId1, chatId3)
	}
}

func TestGetOpenChatEmpty(t *testing.T) {
	t.Logf("TestGetOpenChatEmpty started")
	user := app.User{
		Id:   uint(rand.Uint32()),
		Name: "John",
		Type: app.UserType(app.UserTypeBasic),
	}
	open := app1.GetOpenChat(user.Id)
	if open != nil {
		t.Fatalf("TestGetOpenChatEmpty expected nil, got [%v]", open)
	}
}

func TestGetOpenChat(t *testing.T) {
	t.Logf("TestGetOpenChat started")
	user := app.User{
		Id:   uint(rand.Uint32()),
		Name: "John",
		Type: app.UserType(app.UserTypeBasic),
	}
	_ = app1.AddChat(&user, "TestChat1")
	chatId2 := app1.AddChat(&user, "TestChat2")
	chatId3 := app1.AddChat(&user, "TestChat3")
	open3 := app1.GetOpenChat(user.Id)
	if open3 == nil {
		t.Fatalf("TestGetOpenChat expected a chat[%d], got nil", chatId3)
	}
	if open3.Id != chatId3 {
		t.Fatalf("TestGetOpenChat expected chat[%d], got [%d]", chatId3, open3.Id)
	}
	open2, err := app1.OpenChat(user.Id, chatId2)
	if err != nil {
		t.Fatalf("TestGetOpenChat expected no error, %s", err.Error())
	}
	if open2 == nil {
		t.Fatalf("TestGetOpenChat expected a chat, got nil")
	}
	if open2.Id != chatId2 {
		t.Fatalf("TestGetOpenChat expected chat[%d], got [%d]", chatId2, open2.Id)
	}
	open := app1.GetOpenChat(user.Id)
	if open == nil {
		t.Fatalf("TestGetOpenChat expected a chat[%d], got nil", chatId3)
	}
	if open.Id != chatId2 {
		t.Fatalf("TestGetOpenChat expected chat[%d], got [%d]", chatId2, open.Id)
	}
}

func TestCloseChat(t *testing.T) {
	t.Logf("TestCloseChat started")
	user := app.User{
		Id:   uint(rand.Uint32()),
		Name: "John",
		Type: app.UserType(app.UserTypeBasic),
	}
	chatId := app1.AddChat(&user, "TestChat1")
	open := app1.GetOpenChat(user.Id)
	if open == nil {
		t.Fatalf("TestCloseChat expected chat[%d], got nil", chatId)
	}
	if open.Id != chatId {
		t.Fatalf("TestCloseChat expected chat[%d], got [%d]", chatId, open.Id)
	}
	err := app1.CloseChat(user.Id, chatId)
	if err != nil {
		t.Fatalf("TestCloseChat failed to close chat, %s", err.Error())
	}
	open = app1.GetOpenChat(user.Id)
	if open != nil {
		t.Fatalf("TestCloseChat expected chat[%d] to be closed, but got[%v]", chatId, open)
	}
}

func TestDeleteChat(t *testing.T) {
	t.Logf("TestDeleteChat started")
	user := app.User{
		Id:   uint(rand.Uint32()),
		Name: "John",
		Type: app.UserType(app.UserTypeBasic),
	}
	chatId := app1.AddChat(&user, "TestChat1")
	open := app1.GetOpenChat(user.Id)
	if open == nil {
		t.Fatalf("TestDeleteChat expected open chat[%d], got nil", chatId)
	}
	if open.Id != chatId {
		t.Fatalf("TestDeleteChat expected open chat[%d], got [%d]", chatId, open.Id)
	}
	chat, err := app1.GetChat(user.Id, chatId)
	if err != nil {
		t.Fatalf("TestDeleteChat failed to get chat[%d] for user[%d], %s",
			chatId, user.Id, err.Error())
	}
	if chat == nil {
		t.Fatalf("TestDeleteChat got nil instead of chat[%d] for user[%d]", chatId, user.Id)
	}
	err = app1.DeleteChat(user.Id, open)
	if err != nil {
		t.Fatalf("TestDeleteChat failed to delete chat, %s", err.Error())
	}
	open = app1.GetOpenChat(user.Id)
	if open != nil {
		t.Fatalf("TestDeleteChat expected chat[%d] to be closed, but got[%v]", chatId, open)
	}
	chat, err = app1.GetChat(user.Id, chatId)
	if err == nil {
		t.Fatalf("TestDeleteChat expected nil but got chat[%d]", chatId)
	}
	if chat != nil {
		t.Fatalf("TestDeleteChat expected nil but got chat[%d] for user[%d]", chatId, user.Id)
	}
}
