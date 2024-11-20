package sse

import (
	"log"
	"neon-chat/src/state"
)

// TODO queue mechanism for delta updates
func PollUpdates(state *state.State, conn *state.Conn, pollingUserId uint) {
	log.Printf("TRACE [%s] Live updates triggered by user[%d]\n", conn.Origin, conn.User.Id)
	done := false

	for !done {
		select {
		case <-conn.Reader.Context().Done():
			log.Printf("TRACE [%s] user[%d] conn disonnected\n", conn.Origin, pollingUserId)
			done = true
		case up := <-conn.In:
			log.Printf("TRACE [%s] user[%d] is receiving update[%s]\n", conn.Origin, conn.User.Id, up.Event)
			conn.SendUpdates(up, pollingUserId)
		}
	}
	log.Printf("TRACE [%s] user[%d] disconnects as done:[%t]\n", conn.Origin, conn.User.Id, done)
}
