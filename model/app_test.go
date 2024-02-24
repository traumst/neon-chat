package model

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"log"
)

func TestApp_PollUpdatesForUser(t *testing.T) {
	// now := time.Now()
	// timestamp := now.Format(time.RFC3339)
	// date := strings.Split(timestamp, "T")[0]
	// logPath := fmt.Sprintf("test/from-%s.log", date)
	// file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY, 0666)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.SetOutput(file)

	app := &App{
		State: AppState{
			chats: ChatList{
				chats:  make([]*Chat, 0),
				userAt: make(map[string]*Chat),
			},
			userConn: make(UserConn, 0),
		},
	}

	user1 := "user1"
	chatID1 := app.State.AddChat(user1, "chat1")
	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()
	req1, _ := http.NewRequestWithContext(ctx1, "GET", "/some-route", nil)
	conn1 := app.State.ReplaceConn(
		httptest.NewRecorder(),
		*req1,
		user1,
	)

	user2 := "user2"
	chatID2 := app.State.AddChat(user2, "chat2")
	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()
	req2, _ := http.NewRequestWithContext(ctx2, "GET", "/some-route", nil)
	conn2 := app.State.ReplaceConn(
		httptest.NewRecorder(),
		*req2,
		user2,
	)

	go app.PollUpdatesForUser(conn1, user1)
	go app.PollUpdatesForUser(conn2, user2)

	conn1.In <- UserUpdate{
		Type:    ChatUpdate,
		ChatID:  chatID1,
		Author:  user1,
		RawHtml: "user1: chat1: message1",
	}
	conn2.In <- UserUpdate{
		Type:    ChatUpdate,
		ChatID:  chatID2,
		Author:  user2,
		RawHtml: "user2: chat2: message2",
	}

	tick := time.NewTicker(10 * time.Second)
outerLoop:
	for i := 0; i <= 1; i++ {
		select {
		case e := <-conn1.Out:
			log.Printf("conn1.Out <- update, %s", e.RawHtml)
			cancel1()
		case e := <-conn2.Out:
			log.Printf("conn2.Out <- update, %s", e.RawHtml)
			cancel2()
		case <-tick.C:
			t.Errorf("TestApp_PollUpdatesForUser expected 2 updates, got %d", i)
			break outerLoop
		}
	}
}
