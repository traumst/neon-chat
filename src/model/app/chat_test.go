package app

import "testing"

func TestDefaultChat(t *testing.T) {
	t.Logf("TestDefaultChat started")
	c := Chat{}
	if c.Id != 0 {
		t.Errorf("TestDefaultChat expected UnknownUpdate, got [%d]", c.Id)
	} else if c.Name != "" {
		t.Errorf("TestDefaultChat expected empty user, got [%s]", c.Name)
	} else if c.Owner != nil {
		t.Errorf("TestDefaultChat expected empty owner, got [%v]", c.Owner)
	} else if len(c.users) != 0 {
		t.Errorf("TestDefaultChat expected empty users, got [%v]", c.users)
	}
	t.Logf("TestDefaultChat finished")
}

func TestIsUserInChat(t *testing.T) {
	t.Logf("TestIsUserInChat started")
	u1 := User{Id: 1, Name: "John", Type: UserType(UserTypeBasic)}
	u2 := User{Id: 2 /*Name: "Jill", Type: UserType(UserTypeBasic)*/}
	c := Chat{users: []*User{&u1}}
	if c.IsUserInChat(u2.Id) {
		t.Errorf("TestIsUserInChat [%d] should not have been in the chat with [%d]", u1.Id, u2.Id)
	}
	if !c.IsUserInChat(u1.Id) {
		t.Errorf("TestIsUserInChat [%d] should have been in the chat", u1.Id)
	}
	t.Logf("TestIsUserInChat finished")
}

func TestIsAuthor(t *testing.T) {
	t.Logf("TestIsAuthor started")
	u1 := User{Id: 1, Name: "John", Type: UserType(UserTypeBasic)}
	u2 := User{Id: 2, Name: "Jill", Type: UserType(UserTypeBasic)}
	o := User{Id: 3, Name: "Mr Bill", Type: UserType(UserTypeBasic)}
	c := Chat{Id: 11, Owner: &o, users: []*User{&u1, &u2}}
	m := Message{ChatId: c.Id, Author: &u1, Text: "test message"}
	mwid, err := c.AddMessage(u1.Id, m)
	//mwid, err := c.history.Add(&m)
	if err != nil {
		t.Errorf("TestIsAuthor failed to add message: %s", err)
	}
	if mwid.Id != m.Id {
		t.Errorf("TestIsAuthor expected message id [%d], got [%d]", m.Id, mwid.Id)
		return
	}
	if mwid.ChatId != m.ChatId {
		t.Errorf("TestIsAuthor expected chat id [%d], got [%d]", m.ChatId, mwid.ChatId)
		return
	}
	if c.isAuthor(o.Id, mwid.Id) {
		t.Errorf("TestIsAuthor [%d] should not have been the author of the message, [%v]", o.Id, mwid)
		return
	}
	if mwid.Author.Id != u1.Id {
		t.Errorf("TestIsAuthor expected author id [%d], got [%d]", u1.Id, mwid.Author.Id)
		return
	}
	if !c.isAuthor(u1.Id, mwid.Id) {
		t.Errorf("TestIsAuthor expected [%d] but was [%d], [%+v]", u1.Id, mwid.Author.Id, mwid)
		return
	}
	t.Logf("TestIsOwner finished")
}

func TestIsOwner(t *testing.T) {
	t.Logf("TestIsOwner started")
	o := User{Id: 3, Name: "Mr Bill", Type: UserType(UserTypeBasic)}
	c := Chat{Owner: &o}
	if c.isOwner(123) {
		t.Errorf("TestIsOwner [%d] should have been the owner", o.Id)
	}
	if !c.isOwner(o.Id) {
		t.Errorf("TestIsOwner [%d] should have been the owner", o.Id)
	}
	t.Logf("TestIsOwner finished")
}
