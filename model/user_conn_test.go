package model

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	"go.chat/model/app"
	"go.chat/utils"
)

func TestChannels(t *testing.T) {
	t.Logf("TestChannel started")
	ch1 := make(chan string, 4)
	defer close(ch1)
	ch2 := make(chan string, 4)
	defer close(ch2)
	var wg sync.WaitGroup

	go func() {
		for m := range ch1 {
			t.Logf("TestChannel ch1 received [%s]", m)
			wg.Done()
		}
	}()
	go func() {
		for m := range ch2 {
			t.Logf("TestChannel ch2 received [%s]", m)
			wg.Done()
		}
	}()

	for i := 0; i < 32; i++ {
		msg := fmt.Sprintf("msg-%d", i)
		wg.Add(2)
		t.Logf("TestChannel sent [%s]", msg)
		go func() { ch1 <- msg }()
		go func() { ch2 <- msg }()
	}

	t.Logf("TestChannel all messages sent")
	wg.Wait()
	t.Logf("TestChannel all messages received")
}

func TestIsConnEmpty(t *testing.T) {
	t.Logf("TestIsConnEmpty started")
	uc := make(UserConn, 0)
	user := app.User{Id: 1, Name: "John", Type: app.UserType(app.Free)}
	isConn, conn := uc.IsConn(user.Id)
	if isConn {
		t.Errorf("TestIsConnEmpty user was not supposed to be conn [%+v]", conn)
	} else if conn != nil {
		t.Errorf("TestIsConnEmpty expected empty, got [%+v]", conn)
	}
}

func TestIsConn(t *testing.T) {
	t.Logf("TestIsConn started")
	r := httptest.NewRequest("GET", "/some-route", nil)
	w := httptest.NewRecorder()
	uc := UserConn{}
	user := app.User{Id: 1, Name: "John", Type: app.UserType(app.Free)}
	uc.Add(&user, "test_origin", w, *r)
	isConn, conn := uc.IsConn(user.Id)
	if !isConn {
		t.Errorf("TestIsConn user was supposed to be conn")
	} else if conn == nil {
		t.Errorf("TestIsConn expected to have connection")
	}
}

func TestAdd(t *testing.T) {
	t.Logf("TestAdd started")
	r := httptest.NewRequest("GET", "/some-route", nil)
	w := httptest.NewRecorder()
	uc := UserConn{}
	user := app.User{Id: 1, Name: "John", Type: app.UserType(app.Free)}
	conn := uc.Add(&user, "test_origin", w, *r)
	if conn == nil {
		t.Errorf("TestAdd expected conn, got NIL")
	} else if conn.User.Id != user.Id || conn.Origin != "test_origin" {
		t.Errorf("TestAdd unexpected user[%d] and origin[%s]", conn.User.Id, conn.Origin)
	}
}

func TestGetEmpty(t *testing.T) {
	t.Logf("TestGetEmpty started")
	uc := UserConn{}
	user := app.User{Id: 1, Name: "John", Type: app.UserType(app.Free)}
	conn, err := uc.Get(user.Id)
	if conn != nil {
		t.Errorf("TestGetEmpty expected empty, got [%+v]", conn)
	} else if err == nil {
		t.Errorf("TestGetEmpty expected error, got NIL")
	} else if !strings.Contains(err.Error(), "not connected") ||
		!strings.Contains(err.Error(), strconv.FormatUint(uint64(user.Id), 10)) {
		t.Errorf("TestGetEmpty unexpected error [%s]", err)
	}
}

func TestGet(t *testing.T) {
	t.Logf("TestGet started")
	r := httptest.NewRequest("GET", "/some-route", nil)
	reqId := "test-req-id"
	utils.SetReqId(r, &reqId)
	w := httptest.NewRecorder()
	uc := UserConn{}
	user := app.User{Id: 1, Name: "John", Type: app.UserType(app.Free)}
	conn := uc.Add(&user, "test_origin", w, *r)
	if conn == nil {
		t.Errorf("TestGet expected conn, got NIL")
	}
	conn2, err := uc.Get(user.Id)
	if conn2 == nil {
		t.Errorf("TestGetEmpty expected conn2, got [%+v]", conn)
	} else if err != nil {
		t.Errorf("TestGetEmpty unexpected exception [%s]", err)
	} else if conn.User != conn2.User ||
		conn.Origin != conn2.Origin ||
		conn.In != conn2.In ||
		conn.Writer != conn2.Writer {
		t.Errorf("TestGetEmpty expected equality, got [%+v], [%+v]", conn, conn2)
	}
}

func TestDropEmpty(t *testing.T) {
	t.Logf("TestDropEmpty started")
	uc := UserConn{}
	conn := Conn{}
	uc.Drop(&conn)
}

func TestDrop(t *testing.T) {
	t.Logf("TestDrop started")
	r, err := http.NewRequest("GET", "/some-route", nil)
	reqId := "test-req-id"
	utils.SetReqId(r, &reqId)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	uc := make(UserConn, 0)
	user := app.User{Id: 1, Name: "John", Type: app.UserType(app.Free)}
	addedConn := uc.Add(&user, "test_origin", w, *r)
	if addedConn == nil {
		t.Errorf("failed to add conn")
	}
	conn, err := uc.Get(user.Id)
	if err != nil || conn == nil {
		t.Errorf("TestDrop unexpected error, %s", err)
	}
	uc.Drop(conn)
	conn, err = uc.Get(user.Id)
	if conn != nil {
		t.Errorf("TestDrop expected empty, got [%+v]", conn)
	} else if err == nil {
		t.Errorf("TestDrop expected error, got NIL")
	} else if !strings.Contains(err.Error(), "not connected") ||
		!strings.Contains(err.Error(), strconv.FormatUint(uint64(user.Id), 10)) {
		t.Errorf("TestDrop unexpected error [%s]", err)
	}
}

func TestUserConns(t *testing.T) {
	t.Logf("TestUserConns started")
	r := httptest.NewRequest("GET", "/some-route", nil)
	w := httptest.NewRecorder()
	uc := UserConn{}
	user := app.User{Id: 1, Name: "John", Type: app.UserType(app.Free)}
	conn1 := uc.Add(&user, "test_origin_1", w, *r)
	if conn1 == nil {
		t.Errorf("TestUserConns expected conn1, got NIL")
	}
	user2 := app.User{Id: 2, Name: "Jill", Type: app.UserType(app.Free)}
	conn2 := uc.Add(&user2, "test_origin_2", w, *r)
	if conn2 == nil {
		t.Errorf("TestUserConns expected conn2, got NIL")
	}
	user3 := app.User{Id: 3, Name: "Bill", Type: app.UserType(app.Free)}
	conn3 := uc.Add(&user3, "test_origin_3", w, *r)
	if conn3 == nil {
		t.Errorf("TestUserConns expected conn3, got NIL")
	}

	conns := uc.userConns(user2.Id)
	if len(conns) != 1 {
		t.Errorf("TestUserConns expected 1 conn, got [%d]", len(conns))
	} else if conns[0].ID != conn2.ID ||
		conns[0].Origin != conn2.Origin ||
		conns[0].User.Id != conn2.User.Id {
		t.Errorf("TestUserConns expected conn[%d|%s|%d], got conn[%d|%s|%d]",
			conn2.ID, conn2.Origin, conn2.User.Id,
			conns[0].ID, conns[0].Origin, conns[0].User.Id)
	}
}
