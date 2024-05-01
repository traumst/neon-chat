package handler

import (
	"log"
)

// TODO
// I want that updates, coming in the second case of select,
// to be first put into a queue and only then processed.
// if processing fails I want to keep them
// if processing succeeds I want to remove them
// I want to keep them for a maximum of 10 seconds
// if the same pollingUserId is connected again - try and send the queue first
func PollUpdatesForUser(app *AppState, conn *Conn, pollingUserId uint) {
	log.Printf("∞--%s--> APP.PollUpdatesForUser TRACE IN, triggered by [%d]\n",
		conn.Origin, conn.User.Id)
	done := false
	for !done {
		select {
		case <-conn.Reader.Context().Done():
			log.Printf("<--%s--∞ APP.PollUpdatesForUser WARN user[%d] conn[%v] disonnected\n",
				conn.Origin, pollingUserId, conn.Origin)
			done = true
		case up := <-conn.In:
			log.Printf("∞--%s--> APP.PollUpdatesForUser TRACE user[%d] is receiving update[%s]\n",
				conn.Origin, conn.User.Id, up.Event)
			conn.sendUpdates(up, pollingUserId)
		}
	}
}
