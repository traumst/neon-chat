package state

// import (
// 	"math/rand"
// 	"net/http/httptest"
// 	"testing"

// 	"prplchat/src/model/app"
// 	"prplchat/src/utils"
// 	h "prplchat/src/utils/http"
// )

// // CONN
// func TestAddConn(t *testing.T) {
// 	t.Logf("TestAddConn started")
// 	w := httptest.NewRecorder()
// 	r := httptest.NewRequest("GET", "/some-route", nil)
// 	user := app.User{Id: 1, Name: "John", Type: app.UserType(app.UserTypeBasic)}
// 	app1 := &Application
// 	app1.Init(utils.Config{})
// 	conn1 := app1.AddConn(w, *r, &user, nil)
// 	if conn1 == nil {
// 		t.Fatalf("TestAddConn expected a conn1, got nil")
// 	}
// 	conn2 := app1.AddConn(w, *r, &user, nil)
// 	if conn2 == nil {
// 		t.Fatalf("TestAddConn expected a conn2, got nil")
// 	}
// 	if conn1 == conn2 {
// 		t.Fatalf("TestAddConn expected conn1 and conn2 to be different")
// 	}
// }

// func TestGetConn(t *testing.T) {
// 	t.Logf("TestGetConn started")
// 	w := httptest.NewRecorder()
// 	r := httptest.NewRequest("GET", "/some-route", nil)
// 	reqId := "test-req-id"
// 	reqId = h.SetReqId(r, &reqId)
// 	user := app.User{Id: uint(rand.Uint32()), Name: "John", Type: app.UserType(app.UserTypeBasic)}
// 	app1 := &Application
// 	app1.Init(utils.Config{})
// 	expect := app1.AddConn(w, *r, &user, nil)
// 	if expect == nil {
// 		t.Fatalf("TestGetConn expected a conn, got nil")
// 	}
// 	conns := app1.GetConn(user.Id)
// 	if len(conns) == 0 {
// 		t.Fatalf("TestGetConn expected connections were empty, count [%d]", len(conns))
// 	}
// 	var conn *Conn
// 	for _, c := range conns {
// 		if c.Id == expect.Id {
// 			conn = c
// 			break
// 		}
// 	}
// 	if expect == nil || conn == nil {
// 		t.Fatalf("TestGetConn expected connections, got\n expected[%v]\n conn[%v]", expect, conn)
// 	}
// 	if expect.User != conn.User {
// 		t.Fatalf("TestGetConn expected equality,\nexpected user[%v],\nconn origin[%v]", expect.User, conn.User)
// 	}
// 	if expect.Origin != conn.Origin {
// 		t.Fatalf("TestGetConn expected equality,\nexpected origin[%v],\nconn origin[%v]", expect.Origin, conn.Origin)
// 	}
// 	if expect.In != conn.In {
// 		t.Fatalf("TestGetConn expected equality,\nexpected in[%v],\nconn in[%v]", expect.In, conn.In)
// 	}
// 	if expect.Writer != conn.Writer {
// 		t.Fatalf("TestGetConn expected equality,\nexpected writer[%v],\nconn writer[%v]", expect.Writer, conn.Writer)
// 	}
// }

// func TestDropConn(t *testing.T) {
// 	t.Logf("TestDropConn started")
// 	w := httptest.NewRecorder()
// 	r := httptest.NewRequest("GET", "/some-route", nil)
// 	reqId := "test-req-id"
// 	reqId = h.SetReqId(r, &reqId)
// 	user := app.User{Id: uint(rand.Uint32()), Name: "John", Type: app.UserType(app.UserTypeBasic)}
// 	app1 := &Application
// 	app1.Init(utils.Config{})
// 	expect := app1.AddConn(w, *r, &user, nil)
// 	if expect == nil {
// 		t.Fatalf("TestDropConn expected a conn, got nil")
// 	}
// 	conns := app1.GetConn(user.Id)
// 	if len(conns) == 0 {
// 		t.Fatalf("TestDropConn expected connections were empty, count [%d]", len(conns))
// 	}
// 	err := app1.DropConn(expect)
// 	if err != nil {
// 		t.Fatalf("TestDropConn expected no error, got [%s]", err.Error())
// 	}
// 	conns = app1.GetConn(user.Id)
// 	if len(conns) != 0 {
// 		t.Fatalf("TestGetConn user still has [%d] connections[%v]", len(conns), conns)
// 	}
// }

// // USER
// func TestDropUser(t *testing.T) {
// 	t.Logf("TestDropUser started")
// 	user := app.User{
// 		Id:   uint(rand.Uint32()),
// 		Name: "John",
// 		Type: app.UserType(app.UserTypeBasic),
// 	}
// 	invitee := app.User{Id: uint(rand.Uint32()),
// 		Name: "Jane",
// 		Type: app.UserType(app.UserTypeBasic),
// 	}
// 	app1 := &Application
// 	app1.Init(utils.Config{})
// 	chatId := uint(11)
// 	app1.AddChat(chatId, "TestChat", &user)
// 	chat, _ := app1.GetChat(user.Id, chatId)
// 	err := chat.AddUser(user.Id, &invitee)
// 	if err != nil {
// 		t.Fatalf("TestDropUser failed to add user to chat, %s", err.Error())
// 	}
// 	users, err := chat.GetUsers(invitee.Id)
// 	if err != nil || users == nil {
// 		t.Fatalf("TestDropUser expected user[%d] to be invited, but: %s", invitee.Id, err)
// 	}
// 	err = app1.ExpelFromChat(user.Id, chatId, invitee.Id)
// 	if err != nil {
// 		t.Fatalf("TestDropUser expected no error, %s", err.Error())
// 	}
// 	users, err = chat.GetUsers(invitee.Id)
// 	if err == nil || users != nil {
// 		t.Fatalf("TestDropUser expected user[%d] to be expelled, got [%v]", invitee.Id, users)
// 	}
// }

// func TestUpdateUser(t *testing.T) {
// 	user := app.User{
// 		Id:     uint(rand.Uint32()),
// 		Name:   "John",
// 		Email:  "aaa@aaa.aaa",
// 		Type:   app.UserType(app.UserTypeBasic),
// 		Status: app.UserStatusPending,
// 	}
// 	invitee := app.User{Id: uint(rand.Uint32()),
// 		Name: "Jane",
// 		Type: app.UserType(app.UserTypeBasic),
// 	}
// 	app1 := &Application
// 	app1.Init(utils.Config{})
// 	chatId := uint(11)
// 	app1.AddChat(chatId, "TestChat", &user)
// 	chat, _ := app1.GetChat(user.Id, chatId)
// 	_ = chat.AddUser(user.Id, &invitee)
// 	err := app1.UpdateUser(
// 		user.Id,
// 		&app.User{
// 			Id:     user.Id,
// 			Name:   "Johnny",
// 			Email:  "bbb@bbb.bbb",
// 			Type:   app.UserType(app.UserTypeBasic),
// 			Status: app.UserStatusActive,
// 			Salt:   "CANNOT_CHANGE",
// 		})
// 	if err != nil {
// 		t.Fatalf("TestUpdateUser failed to update user[%d], %s", user.Id, err.Error())
// 	}
// 	updatedUsers, err := chat.GetUsers(user.Id)
// 	if err != nil {
// 		t.Fatalf("TestUpdateUser failed to get users from chat[%d], %s", chatId, err.Error())
// 	}
// 	for _, uu := range updatedUsers {
// 		if uu.Id != user.Id {
// 			continue
// 		}
// 		if uu != &user {
// 			t.Fatalf("TestUpdateUser broken pointer user[%s] and user[%s] should be the same", uu.Name, user.Name)
// 		}
// 		if uu.Name != "Johnny" {
// 			t.Fatalf("TestUpdateUser expected name[Johnny] got name[%s]", uu.Name)
// 		}
// 		if uu.Email != "bbb@bbb.bbb" {
// 			t.Fatalf("TestUpdateUser expected name[Johnny] got name[%s]", uu.Name)
// 		}
// 		if uu.Salt != user.Salt {
// 			t.Fatalf("TestUpdateUser expected salt[%s] changed[%s]", user.Salt, uu.Salt)
// 		}
// 		return
// 	}
// 	t.Fatalf("TestUpdateUser user[%d] not found", user.Id)
// }

// // CHAT
// func TestAddChat(t *testing.T) {
// 	t.Logf("TestAddChat started")
// 	user := app.User{
// 		Id:   uint(rand.Uint32()),
// 		Name: "John",
// 		Type: app.UserType(app.UserTypeBasic),
// 	}
// 	app1 := &Application
// 	app1.Init(utils.Config{})
// 	chatId := uint(11)
// 	err := app1.AddChat(chatId, "TestChat", &user)
// 	if err != nil {
// 		t.Fatalf("TestAddChat failed to add user[%d] into chat[%d], %s", user.Id, chatId, err.Error())
// 	}
// 	chat, err := app1.GetChat(user.Id, chatId)
// 	if err != nil {
// 		t.Fatalf("TestAddChat failed to get chat[%d] for user[%d], [%s]",
// 			chatId, user.Id, err.Error())
// 	}
// 	if chat == nil || chat.Id != chatId {
// 		t.Fatalf("TestAddChat expected a chat[%d] for user[%d] got nil", chatId, user.Id)
// 	}
// }

// func TestGetChats(t *testing.T) {
// 	t.Logf("TestGetChats started")
// 	user := app.User{
// 		Id:   uint(rand.Uint32()),
// 		Name: "John",
// 		Type: app.UserType(app.UserTypeBasic),
// 	}
// 	app1 := &Application
// 	app1.Init(utils.Config{})
// 	chatId1 := app1.AddChat(11, "TestChat1", &user)
// 	chatId2 := app1.AddChat(22, "TestChat2", &user)
// 	chats := app1.GetChats(user.Id)
// 	if len(chats) != 2 {
// 		t.Fatalf("TestGetChats expected 2 chats, got [%d]", len(chats))
// 	}
// 	if chats[0].Id != 11 || chats[1].Id != 22 {
// 		t.Fatalf("TestGetChats expected chat ids [%d, %d], got [%d, %d]",
// 			chatId1, chatId2, chats[0].Id, chats[1].Id)
// 	}
// }

// func TestGetChat(t *testing.T) {
// 	t.Logf("TestGetChats started")
// 	user := app.User{
// 		Id:   uint(rand.Uint32()),
// 		Name: "John",
// 		Type: app.UserType(app.UserTypeBasic),
// 	}
// 	app1 := &Application
// 	app1.Init(utils.Config{})
// 	app1.AddChat(11, "TestChat1", &user)
// 	chat1, err := app1.GetChat(user.Id, 11)
// 	if err != nil {
// 		t.Fatalf("TestGetChat failed to get chat[%d] for user[%d], [%s]",
// 			11, user.Id, err.Error())
// 	}
// 	if chat1 == nil {
// 		t.Fatalf("TestGetChat expected a chat[%d] for user[%d] got nil", 11, user.Id)
// 	}
// 	app1.AddChat(22, "TestChat2", &user)
// 	chat2, err := app1.GetChat(user.Id, 22)
// 	if err != nil {
// 		t.Fatalf("TestGetChat failed to get chat[%d] for user[%d], %s",
// 			11, user.Id, err.Error())
// 	}
// 	if chat2 == nil {
// 		t.Fatalf("TestGetChat expected a chat[%d] for user[%d] got nil", 22, user.Id)
// 	}
// }

// func TestOpenChat(t *testing.T) {
// 	t.Logf("TestOpenChat started")
// 	user := app.User{
// 		Id:   uint(rand.Uint32()),
// 		Name: "John",
// 		Type: app.UserType(app.UserTypeBasic),
// 	}
// 	app1 := &Application
// 	app1.Init(utils.Config{})
// 	chatId1 := app1.AddChat(11, "TestChat1", &user)
// 	chatId2 := app1.AddChat(22, "TestChat2", &user)
// 	chatId3 := app1.AddChat(33, "TestChat3", &user)
// 	open, err := app1.OpenChat(user.Id, 22)
// 	if err != nil {
// 		t.Fatalf("TestOpenChat expected no error, %s", err.Error())
// 	}
// 	if open == nil {
// 		t.Fatalf("TestOpenChat expected a chat, got nil")
// 	}
// 	if open.Id != 22 {
// 		t.Fatalf("TestOpenChat expected chat[%d], got [%d]", chatId2, open.Id)
// 	}
// 	if open.Id == 11 || open.Id == 33 {
// 		t.Fatalf("TestOpenChat expected open chat[%d] to be different from [%d, %d]",
// 			chatId2, chatId1, chatId3)
// 	}
// }

// func TestGetOpenChatEmpty(t *testing.T) {
// 	t.Logf("TestGetOpenChatEmpty started")
// 	user := app.User{
// 		Id: uint(rand.Uint32()),
// 		// Name: "John",
// 		// Type: app.UserType(app.UserTypeBasic),
// 	}
// 	app1 := &Application
// 	app1.Init(utils.Config{})
// 	open := app1.GetOpenChat(user.Id)
// 	if open != nil {
// 		t.Fatalf("TestGetOpenChatEmpty expected nil, got [%v]", open)
// 	}
// }

// func TestGetOpenChat(t *testing.T) {
// 	t.Logf("TestGetOpenChat started")
// 	user := app.User{
// 		Id:   uint(rand.Uint32()),
// 		Name: "John",
// 		Type: app.UserType(app.UserTypeBasic),
// 	}
// 	app1 := &Application
// 	app1.Init(utils.Config{})
// 	app1.AddChat(11, "TestChat1", &user)
// 	app1.AddChat(22, "TestChat2", &user)
// 	app1.AddChat(33, "TestChat3", &user)
// 	open3 := app1.GetOpenChat(user.Id)
// 	if open3 == nil {
// 		t.Fatalf("TestGetOpenChat expected a chat[%d], got nil", 33)
// 	}
// 	if open3.Id != 33 {
// 		t.Fatalf("TestGetOpenChat expected chat[%d], got [%d]", 33, open3.Id)
// 	}
// 	open2, err := app1.OpenChat(user.Id, 22)
// 	if err != nil {
// 		t.Fatalf("TestGetOpenChat expected no error, %s", err.Error())
// 	}
// 	if open2 == nil {
// 		t.Fatalf("TestGetOpenChat expected a chat, got nil")
// 	}
// 	if open2.Id != 22 {
// 		t.Fatalf("TestGetOpenChat expected chat[%d], got [%d]", 22, open2.Id)
// 	}
// 	open := app1.GetOpenChat(user.Id)
// 	if open == nil {
// 		t.Fatalf("TestGetOpenChat expected a chat[%d], got nil", 22)
// 	}
// 	if open.Id != 22 {
// 		t.Fatalf("TestGetOpenChat expected chat[%d], got [%d]", 22, open.Id)
// 	}
// }

// func TestCloseChat(t *testing.T) {
// 	t.Logf("TestCloseChat started")
// 	user := app.User{
// 		Id:   uint(rand.Uint32()),
// 		Name: "John",
// 		Type: app.UserType(app.UserTypeBasic),
// 	}
// 	app1 := &Application
// 	app1.Init(utils.Config{})
// 	app1.AddChat(11, "TestChat1", &user)
// 	open := app1.GetOpenChat(user.Id)
// 	if open == nil {
// 		t.Fatalf("TestCloseChat expected chat[%d], got nil", 11)
// 	}
// 	if open.Id != 11 {
// 		t.Fatalf("TestCloseChat expected chat[%d], got [%d]", 11, open.Id)
// 	}
// 	err := app1.CloseChat(user.Id, 11)
// 	if err != nil {
// 		t.Fatalf("TestCloseChat failed to close chat, %s", err.Error())
// 	}
// 	open = app1.GetOpenChat(user.Id)
// 	if open != nil {
// 		t.Fatalf("TestCloseChat expected chat[%d] to be closed, but got[%v]", 11, open)
// 	}
// }

// func TestDeleteChat(t *testing.T) {
// 	t.Logf("TestDeleteChat started")
// 	user := app.User{
// 		Id:   uint(rand.Uint32()),
// 		Name: "John",
// 		Type: app.UserType(app.UserTypeBasic),
// 	}
// 	app1 := &Application
// 	app1.Init(utils.Config{})
// 	app1.AddChat(11, "TestChat1", &user)
// 	open := app1.GetOpenChat(user.Id)
// 	if open == nil {
// 		t.Fatalf("TestDeleteChat expected open chat[%d], got nil", 11)
// 	}
// 	if open.Id != 11 {
// 		t.Fatalf("TestDeleteChat expected open chat[%d], got [%d]", 11, open.Id)
// 	}
// 	chat, err := app1.GetChat(user.Id, 11)
// 	if err != nil {
// 		t.Fatalf("TestDeleteChat failed to get chat[%d] for user[%d], %s",
// 			11, user.Id, err.Error())
// 	}
// 	if chat == nil {
// 		t.Fatalf("TestDeleteChat got nil instead of chat[%d] for user[%d]", 11, user.Id)
// 	}
// 	err = app1.DeleteChat(user.Id, open)
// 	if err != nil {
// 		t.Fatalf("TestDeleteChat failed to delete chat, %s", err.Error())
// 	}
// 	open = app1.GetOpenChat(user.Id)
// 	if open != nil {
// 		t.Fatalf("TestDeleteChat expected chat[%d] to be closed, but got[%v]", 11, open)
// 	}
// 	chat, err = app1.GetChat(user.Id, 11)
// 	if err == nil {
// 		t.Fatalf("TestDeleteChat expected nil but got chat[%d]", 11)
// 	}
// 	if chat != nil {
// 		t.Fatalf("TestDeleteChat expected nil but got chat[%d] for user[%d]", 11, user.Id)
// 	}
// }
