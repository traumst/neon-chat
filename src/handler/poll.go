package handler

import (
	"log"
	"sync"
)

func PollUpdatesForUser(conn *Conn, pollingUserId uint) {
	log.Printf("∞--%s--> APP.PollUpdatesForUser TRACE IN, triggered by [%d]\n", conn.Origin, conn.User.Id)
	var wg sync.WaitGroup
	done := false
	for !done {
		select {
		case <-conn.Reader.Context().Done():
			log.Printf("<--%s--∞ APP.PollUpdatesForUser WARN user[%d] conn[%v] disonnected\n",
				conn.Origin, pollingUserId, conn.Origin)
			done = true
		case up := <-conn.In:
			wg.Add(1)
			go func() {
				defer wg.Done()
				log.Printf("∞--%s--> APP.PollUpdatesForUser TRACE user[%d] is receiving update[%s]\n",
					conn.Origin, conn.User.Id, up.Event)
				conn.sendUpdates(up, pollingUserId)
				//conn.Out <- up
			}()
		}
	}
	wg.Wait()
}
