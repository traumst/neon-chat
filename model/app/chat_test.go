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
	u1 := User{Id: 1, Name: "John", Type: UserType(Free)}
	u2 := User{Id: 2, Name: "Jill", Type: UserType(Free)}
	c := Chat{users: []*User{&u1}}
	if c.isUserInChat(u2.Id) {
		t.Errorf("TestIsUserInChat [%d] should not have been in the chat with [%d]", u1.Id, u2.Id)
	}
	if !c.isUserInChat(u1.Id) {
		t.Errorf("TestIsUserInChat [%d] should have been in the chat", u1.Id)
	}
	t.Logf("TestIsUserInChat finished")
}

func TestIsAuthor(t *testing.T) {
	t.Logf("TestIsAuthor started")
	//u := "test-user"
	//o := "other-user"
	u1 := User{Id: 1, Name: "John", Type: UserType(Free)}
	u2 := User{Id: 2, Name: "Jill", Type: UserType(Free)}
	o := User{Id: 3, Name: "Mr Bill", Type: UserType(Free)}
	c := Chat{Owner: &o, users: []*User{&u1, &u2}}
	m := Message{ID: 1, Owner: &o, Author: &u1, Text: "test message"}
	mwid, err := c.history.Add(&m)
	if err != nil {
		t.Errorf("TestIsAuthor failed to add message: %s", err)
	}
	if c.isAuthor(o.Id, mwid.ID) {
		t.Errorf("TestIsAuthor [%d] should not have been the author of the message", o.Id)
	}
	if !c.isAuthor(u1.Id, mwid.ID) {
		t.Errorf("TestIsAuthor [%d] should have been the author of the message", u1.Id)
	}
	t.Logf("TestIsOwner finished")
}

func TestIsOwner(t *testing.T) {
	t.Logf("TestIsOwner started")
	o := User{Id: 3, Name: "Mr Bill", Type: UserType(Free)}
	c := Chat{Owner: &o}
	if c.isOwner(123) {
		t.Errorf("TestIsOwner [%d] should have been the owner", o.Id)
	}
	if !c.isOwner(o.Id) {
		t.Errorf("TestIsOwner [%d] should have been the owner", o.Id)
	}
	t.Logf("TestIsOwner finished")
}
