package handler

import (
	"log"
)

// TODO queue mechanism for delta updates
func PollUpdatesForUser(app *AppState, conn *Conn, pollingUserId uint) {
	log.Printf("[%s] APP.PollUpdatesForUser TRACE IN, triggered by [%d]\n",
		conn.Origin, conn.User.Id)
	done := false
	for !done {
		select {
		case <-conn.Reader.Context().Done():
			log.Printf("[%s] APP.PollUpdatesForUser DEBUG user[%d] conn[%v] disonnected\n",
				conn.Origin, pollingUserId, conn.Origin)
			done = true
		case up := <-conn.In:
			log.Printf("[%s] APP.PollUpdatesForUser TRACE user[%d] is receiving update[%s]\n",
				conn.Origin, conn.User.Id, up.Event)
			conn.SendUpdates(up, pollingUserId)
		}
	}
}
