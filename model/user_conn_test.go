package model

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

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
	isConn, conn := uc.IsConn("test_user")
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
	uc.Add("test_user", "test_origin", w, *r)
	isConn, conn := uc.IsConn("test_user")
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

	conn := uc.Add("test_user", "test_origin", w, *r)
	if conn == nil {
		t.Errorf("TestAdd expected conn, got NIL")
	} else if conn.User != "test_user" || conn.Origin != "test_origin" {
		t.Errorf("TestAdd unexpected user and origin [%s] [%s]", conn.User, conn.Origin)
	}
}

func TestGetEmpty(t *testing.T) {
	t.Logf("TestGetEmpty started")
	uc := UserConn{}
	conn, err := uc.Get("test_user")
	if conn != nil {
		t.Errorf("TestGetEmpty expected empty, got [%+v]", conn)
	} else if err == nil {
		t.Errorf("TestGetEmpty expected error, got NIL")
	} else if !strings.Contains(err.Error(), "not connected") || !strings.Contains(err.Error(), "test_user") {
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

	conn := uc.Add("test_user", "test_origin", w, *r)
	if conn == nil {
		t.Errorf("TestGet expected conn, got NIL")
	}

	conn2, err := uc.Get("test_user")
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
	testUser := "test_user"
	addedConn := uc.Add(testUser, "test_origin", w, *r)
	if addedConn == nil {
		t.Errorf("failed to add conn")
	}

	conn, err := uc.Get("test_user")
	if err != nil || conn == nil {
		t.Errorf("TestDrop unexpected error, %s", err)
	}

	uc.Drop(conn)
	conn, err = uc.Get("test_user")
	if conn != nil {
		t.Errorf("TestDrop expected empty, got [%+v]", conn)
	} else if err == nil {
		t.Errorf("TestDrop expected error, got NIL")
	} else if !strings.Contains(err.Error(), "not connected") || !strings.Contains(err.Error(), testUser) {
		t.Errorf("TestDrop unexpected error [%s]", err)
	}
}

func TestUserConns(t *testing.T) {
	t.Logf("TestUserConns started")
	r := httptest.NewRequest("GET", "/some-route", nil)
	w := httptest.NewRecorder()
	uc := UserConn{}

	conn1 := uc.Add("test_user_1", "test_origin_1", w, *r)
	if conn1 == nil {
		t.Errorf("TestUserConns expected conn1, got NIL")
	}
	conn2 := uc.Add("test_user_2", "test_origin_2", w, *r)
	if conn2 == nil {
		t.Errorf("TestUserConns expected conn2, got NIL")
	}
	conn3 := uc.Add("test_user_3", "test_origin_3", w, *r)
	if conn3 == nil {
		t.Errorf("TestUserConns expected conn3, got NIL")
	}

	conns := uc.userConns("test_user_2")
	if len(conns) != 1 {
		t.Errorf("TestUserConns expected 1 conn, got [%d]", len(conns))
	} else if conns[0].ID != conn2.ID ||
		conns[0].Origin != conn2.Origin ||
		conns[0].User != conn2.User {
		t.Errorf("TestUserConns expected conn[%d|%s|%s], got conn[%d|%s|%s]",
			conn2.ID, conn2.Origin, conn2.User,
			conns[0].ID, conns[0].Origin, conns[0].User)
	}
}