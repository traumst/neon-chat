package state

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"neon-chat/src/model/app"
	h "neon-chat/src/utils/http"
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
	uc := make(OpenConnections, 0)
	user := app.User{Id: 1 /*, Name: "John", Type: app.UserType(app.UserTypeBasic) */}
	isConn := uc.IsConn(user.Id)
	if isConn {
		t.Errorf("TestIsConnEmpty user was not supposed to be connected")
	}
}

func TestIsConn(t *testing.T) {
	t.Logf("TestIsConn started")
	r := httptest.NewRequest("GET", "/some-route", nil)
	w := httptest.NewRecorder()
	uc := OpenConnections{}
	user := app.User{Id: 1, Name: "John", Type: app.UserType(app.UserTypeBasic)}
	uc.Add(&user, "test_origin", w, *r)
	isConn := uc.IsConn(user.Id)
	if !isConn {
		t.Errorf("TestIsConn user was supposed to be conn")
	}
}

func TestAdd(t *testing.T) {
	t.Logf("TestAdd started")
	r := httptest.NewRequest("GET", "/some-route", nil)
	w := httptest.NewRecorder()
	uc := OpenConnections{}
	user := app.User{Id: 1, Name: "John", Type: app.UserType(app.UserTypeBasic)}
	conn := uc.Add(&user, "test_origin", w, *r)
	if conn == nil {
		t.Errorf("TestAdd expected conn, got NIL")
	} else if conn.User.Id != user.Id || conn.Origin != "test_origin" {
		t.Errorf("TestAdd unexpected user[%d] and origin[%s]", conn.User.Id, conn.Origin)
	}
}

func TestGetEmpty(t *testing.T) {
	t.Logf("TestGetEmpty started")
	uc := OpenConnections{}
	user := app.User{Id: uint(rand.Uint32()) /*Name: "John", Type: app.UserType(app.UserTypeBasic)*/}
	conns := uc.Get(user.Id)
	if len(conns) != 0 {
		t.Errorf("TestGetEmpty expected empty, got [%+v]", conns)
	}
}

func TestGet(t *testing.T) {
	t.Logf("TestGet started")
	r := httptest.NewRequest("GET", "/some-route", nil)
	reqId := "test-req-id"
	h.SetReqId(r, &reqId)
	w := httptest.NewRecorder()
	uc := OpenConnections{}
	user := app.User{Id: uint(rand.Uint32()), Name: "John", Type: app.UserType(app.UserTypeBasic)}
	conn := uc.Add(&user, "test_origin", w, *r)
	if conn == nil {
		t.Errorf("TestGet expected conn, got NIL")
	}
	conn2 := uc.Get(user.Id)
	if len(conn2) == 0 {
		t.Errorf("TestGetEmpty expected conn2, got [%+v]", conn2)
	}
	for _, c := range conn2 {
		if conn.User != c.User ||
			conn.Origin != c.Origin ||
			conn.In != c.In ||
			conn.Writer != c.Writer {
			t.Errorf("TestGetEmpty expected equality, got [%+v], [%+v]", conn, conn2)
		}
	}
}

func TestDropEmpty(t *testing.T) {
	t.Logf("TestDropEmpty started")
	uc := OpenConnections{}
	conn := Conn{}
	uc.Drop(&conn)
}

func TestDrop(t *testing.T) {
	t.Logf("TestDrop started")
	r, err := http.NewRequest("GET", "/some-route", nil)
	reqId := "test-req-id"
	h.SetReqId(r, &reqId)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	uc := make(OpenConnections, 0)
	user := app.User{Id: uint(rand.Uint32()), Name: "John", Type: app.UserType(app.UserTypeBasic)}
	addedConn := uc.Add(&user, "test_origin", w, *r)
	if addedConn == nil {
		t.Errorf("failed to add conn")
	}
	conns := uc.Get(user.Id)
	if len(conns) == 0 {
		t.Errorf("TestDrop unexpected error, %s", err)
	}
	uc.Drop(conns[0])
	conns = uc.Get(user.Id)
	if len(conns) != 0 {
		t.Errorf("TestDrop expected empty, got [%+v]", conns)
	}
}

func TestUserConns(t *testing.T) {
	t.Logf("TestUserConns started")
	r := httptest.NewRequest("GET", "/some-route", nil)
	w := httptest.NewRecorder()
	uc := OpenConnections{}
	user := app.User{Id: uint(rand.Uint32()), Name: "John", Type: app.UserType(app.UserTypeBasic)}
	conn1 := uc.Add(&user, "test_origin_1", w, *r)
	if conn1 == nil {
		t.Errorf("TestUserConns expected conn1, got NIL")
	}
	user2 := app.User{Id: 2, Name: "Jill", Type: app.UserType(app.UserTypeBasic)}
	conn2 := uc.Add(&user2, "test_origin_2", w, *r)
	if conn2 == nil {
		t.Errorf("TestUserConns expected conn2, got NIL")
	}
	user3 := app.User{Id: 3, Name: "Bill", Type: app.UserType(app.UserTypeBasic)}
	conn3 := uc.Add(&user3, "test_origin_3", w, *r)
	if conn3 == nil {
		t.Errorf("TestUserConns expected conn3, got NIL")
	}

	conns := uc[user2.Id]
	if len(conns) != 1 {
		t.Errorf("TestUserConns expected 1 conn, got [%d]", len(conns))
	} else if conns[0].Id != conn2.Id ||
		conns[0].Origin != conn2.Origin ||
		conns[0].User.Id != conn2.User.Id {
		t.Errorf("TestUserConns expected conn[%d|%s|%d], got conn[%d|%s|%d]",
			conn2.Id, conn2.Origin, conn2.User.Id,
			conns[0].Id, conns[0].Origin, conns[0].User.Id)
	}
}
