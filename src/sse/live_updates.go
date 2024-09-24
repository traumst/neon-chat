package sse

import (
	"log"
	"neon-chat/src/state"
	"time"
)

// TODO queue mechanism for delta updates
func LiveUpdates(state *state.State, conn *state.Conn, pollingUserId uint) {
	log.Printf("TRACE IN LiveUpdates [%s] triggered by [%d]\n", conn.Origin, conn.User.Id)
	done := false
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for !done {
		select {
		case <-conn.Reader.Context().Done():
			log.Printf("DEBUG LiveUpdates [%s] user[%d] conn[%v] disonnected\n",
				conn.Origin, pollingUserId, conn.Origin)
			done = true
		case up := <-conn.In:
			log.Printf("TRACE LiveUpdates [%s] user[%d] is receiving update[%s]\n",
				conn.Origin, conn.User.Id, up.Event)
			conn.SendUpdates(up, pollingUserId)
		}
	}
	log.Printf("TRACE OUT LiveUpdates [%s] user[%d] disconnects as done [%t]\n",
		conn.Origin, conn.User.Id, done)
}
