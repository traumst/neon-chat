package app

import (
	"math/rand"
	"testing"
)

func TestAddChat(t *testing.T) {
	t.Logf("TestAddChat started")
	cl := NewOpenChats()
	user := User{Id: 1, Name: "John", Type: UserType(UserTypeBasic)}
	cl.AddChat(11, &user, "test-chat")
	user2 := User{Id: 2, Name: "Jane", Type: UserType(UserTypeBasic)}
	cl.AddChat(22, &user2, "test-chat")
	_, err := cl.GetChat(2, 11)
	if err == nil {
		t.Errorf("TestAddChat failed to get chat [%d] for user [], [%s]", 11, err)
		return
	}
	_, err = cl.GetChat(2, 22)
	if err != nil {
		t.Errorf("TestAddChat failed to get chat [%s]", err)
		return

	}
}

func TestOpenChat(t *testing.T) {
	t.Logf("TestOpenChat started")
	cl := NewOpenChats()
	user := User{Id: 1, Name: "John", Type: UserType(UserTypeBasic)}
	cl.AddChat(11, &user, "test-chat")

	user2 := User{Id: 2, Name: "Jill", Type: UserType(UserTypeBasic)}
	cl.AddChat(22, &user2, "test-chat-2")
	openChat2 := cl.GetOpenChat(user2.Id)
	if openChat2 == nil {
		t.Errorf("TestOpenChat openChat was NIL")
		return
	}
	openChat1, err := cl.OpenChat(user.Id, 11)
	if err != nil {
		t.Errorf("TestOpenChat failed to open chat [%s]", err)
		return
	}
	if openChat1 == nil {
		t.Errorf("TestOpenChat openChat was NIL")
		return
	}
	if openChat1.Id == openChat2.Id || openChat1.Id != 11 {
		t.Errorf("TestOpenChat expected chatId %d, got %d", 11, openChat1.Id)
	}
}

func TestGetChat(t *testing.T) {
	t.Logf("TestAddChat started")
	cl := NewOpenChats()
	user := User{Id: uint(rand.Uint32()), Name: "John", Type: UserType(UserTypeBasic)}
	cl.AddChat(11, &user, "test-chat")
	chats := cl.GetChats(user.Id)
	if chats == nil {
		t.Errorf("TestAddChat chats were nil")
	} else if len(chats) != 1 {
		t.Errorf("TestAddChat expected 1 chat, got %d", len(chats))
	}
}

func TestGetOpenChatEmpty(t *testing.T) {
	t.Logf("TestGetOpenChatEmpty started")
	cl := NewOpenChats()
	user := User{Id: 1 /*Name: "John", Type: UserType(UserTypeBasic)*/}
	chat := cl.GetOpenChat(user.Id)
	if chat != nil {
		t.Errorf("TestGetOpenChatEmpty expected chat to be NIL")
	}
}

func TestGetOpenChat(t *testing.T) {
	t.Logf("TestGetOpenChat started")
	cl := NewOpenChats()
	user := User{Id: 1, Name: "John", Type: UserType(UserTypeBasic)}
	cl.AddChat(11, &user, "test-chat")
	chat := cl.GetOpenChat(user.Id)
	if chat == nil {
		t.Errorf("TestGetOpenChat chat was NIL [%+v]", chat)
		return
	}
	cl.AddChat(22, &user, "test-chat-2")
	chat = cl.GetOpenChat(user.Id)
	if chat == nil {
		t.Errorf("TestGetOpenChat chat was NIL [%+v]", chat)
	} else if chat.Id != 22 {
		t.Errorf("TestGetOpenChat expected chatId 1, got [%d]", chat.Id)
	}
}

func TestGetChats(t *testing.T) {
	t.Logf("TestGetChats started")
	cl := NewOpenChats()
	user := User{Id: uint(rand.Uint32()), Name: "John", Type: UserType(UserTypeBasic)}
	cl.AddChat(11, &user, "test-chat")
	cl.AddChat(11, &user, "test-chat2")
	cl.AddChat(11, &user, "test-chat3")
	chats := cl.GetChats(user.Id)
	if chats == nil {
		t.Errorf("TestGetChats chats were nil")
	} else if len(chats) != 1 {
		t.Errorf("TestGetChats expected 1 chat, got [%d]", len(chats))
	}
}

func TestDeleteChatEmpty(t *testing.T) {
	t.Logf("TestDeleteChatEmpty started")
	cl := NewOpenChats()
	user := User{Id: 1, Name: "John", Type: UserType(UserTypeBasic)}
	err := cl.DeleteChat(user.Id, nil)
	if err == nil {
		t.Errorf("TestDeleteChatEmpty expected error")
	}
	chat1 := Chat{Id: 1, users: []*User{&user}, Owner: &user}
	err = cl.DeleteChat(user.Id, &chat1)
	if err == nil {
		t.Errorf("TestDeleteChatEmpty expected error, 1")
	}
	chat2 := Chat{Id: 0, users: []*User{&user}, Owner: &user}
	err = cl.DeleteChat(user.Id, &chat2)
	if err == nil {
		t.Errorf("TestDeleteChatEmpty expected error, 0")
	}
	chat3 := Chat{Id: 0, users: []*User{&user}, Owner: &user}
	err = cl.DeleteChat(user.Id, &chat3)
	if err == nil {
		t.Errorf("TestDeleteChatEmpty expected error, 0")
	}
}

func TestDeleteChat(t *testing.T) {
	t.Logf("TestDeleteChatEmpty started")
	cl := NewOpenChats()
	user := User{Id: 1, Name: "John", Type: UserType(UserTypeBasic)}
	cl.AddChat(11, &user, "test-chat")
	openChat := cl.GetOpenChat(user.Id)
	if openChat == nil {
		t.Errorf("TestDeleteChatEmpty openChat was NIL")
		return
	}
	chat, err := cl.GetChat(user.Id, 11)
	if err != nil {
		t.Errorf("TestDeleteChatEmpty get failed [%s]", err)
		return
	}
	err = cl.DeleteChat(user.Id, chat)
	if err != nil {
		t.Errorf("TestDeleteChatEmpty delete failed [%s]", err)
		return
	}
	openChat = cl.GetOpenChat(user.Id)
	if openChat != nil {
		t.Errorf("TestDeleteChatEmpty openChat was expected to be NIL, but was [%+v]", openChat)
	}
}

func TestDeleteChatNotOwner(t *testing.T) {
	t.Logf("TestDeleteChatEmpty started")
	cl := NewOpenChats()
	user := User{Id: 1, Name: "John", Type: UserType(UserTypeBasic)}
	cl.AddChat(11, &user, "test-chat")
	openChat := cl.GetOpenChat(user.Id)
	if openChat == nil {
		t.Errorf("TestDeleteChatEmpty openChat was NIL")
		return
	}
	chat, err := cl.GetChat(user.Id, 11)
	if err != nil {
		t.Errorf("TestDeleteChatEmpty get failed [%s]", err)
		return
	}
	err = cl.DeleteChat(333, chat)
	if err == nil {
		t.Errorf("TestDeleteChatEmpty delete should not have been allowed")
	}
}

func TestInviteUser(t *testing.T) {
	t.Logf("TestInviteUserEmpty started")
	cl := NewOpenChats()
	owner := User{Id: uint(rand.Uint32()), Name: "John", Type: UserType(UserTypeBasic)}
	cl.AddChat(11, &owner, "test-chat")
	invitee := User{Id: uint(rand.Uint32()), Name: "Jill", Type: UserType(UserTypeBasic)}
	err := cl.InviteUser(owner.Id, 11, &invitee)
	if err != nil {
		t.Errorf("TestInviteUserEmpty fail to invite user [%v], [%v]", invitee, err.Error())
		return
	}
	inviteeChats := cl.GetChats(invitee.Id)
	if inviteeChats == nil {
		t.Errorf("TestInviteUserEmpty inviteeChats were nil")
	} else if len(inviteeChats) != 1 {
		t.Errorf("TestInviteUserEmpty expected 1 chat, got %d", len(inviteeChats))
	} else if inviteeChats[0].Id != 11 {
		t.Errorf("TestInviteUserEmpty expected chatId %d, got %d", 11, inviteeChats[0].Id)
	} else if !inviteeChats[0].IsUserInChat(invitee.Id) {
		t.Errorf("TestInviteUserEmpty user [%v] is not in chat", invitee)
	} else if !inviteeChats[0].IsUserInChat(owner.Id) {
		t.Errorf("TestInviteUserEmpty owner [%v] is not in chat", owner)
	}
}
