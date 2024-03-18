package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"log"

	"go.chat/model"
	"go.chat/model/event"
)

func Test_PollUpdatesForUser(t *testing.T) {
	user1 := "user1"
	app := &model.ApplicationState
	chatID1 := app.AddChat(user1, "chat1")
	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()
	req1, _ := http.NewRequestWithContext(ctx1, "GET", "/some-route", nil)
	conn1 := app.ReplaceConn(
		httptest.NewRecorder(),
		*req1,
		user1,
	)

	user2 := "user2"
	chatID2 := app.AddChat(user2, "chat2")
	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()
	req2, _ := http.NewRequestWithContext(ctx2, "GET", "/some-route", nil)
	conn2 := app.ReplaceConn(
		httptest.NewRecorder(),
		*req2,
		user2,
	)

	go PollUpdatesForUser(conn1, user1)
	go PollUpdatesForUser(conn2, user2)

	conn1.In <- event.LiveUpdate{
		Event:  event.ChatCreated,
		ChatID: chatID1,
		MsgID:  -2,
		Author: user1,
		Data:   "user1: chat1: message1",
	}
	conn2.In <- event.LiveUpdate{
		Event:  event.ChatCreated,
		ChatID: chatID2,
		MsgID:  -2,
		Author: user2,
		Data:   "user2: chat2: message2",
	}

	tick := time.NewTicker(10 * time.Second)
outerLoop:
	for i := 0; i <= 1; i++ {
		select {
		case e := <-conn1.Out:
			log.Printf("conn1.Out <- update, %s", e.Data)
			cancel1()
		case e := <-conn2.Out:
			log.Printf("conn2.Out <- update, %s", e.Data)
			cancel2()
		case <-tick.C:
			t.Errorf("TestApp_PollUpdatesForUser expected 2 updates, got %d", i)
			break outerLoop
		}
	}
}