package sse

import (
	"log"

	"neon-chat/src/state"
	m "neon-chat/src/utils/maintenance"
)

// TODO queue mechanism for delta updates
func PollUpdates(state *state.State, conn *state.Conn, pollingUserId uint) bool {
	log.Printf("TRACE [%s] Live updates triggered by user[%d]\n", conn.Origin, conn.User.Id)
	done := false
	//ticker := time.NewTicker(1 * time.Second)

	for !done && !m.MaintenanceManager.IsInMaintenance() {
		select {
		case <-conn.Reader.Context().Done():
			//log.Printf("TRACE [%s] user[%d] disconnects as done:[%t]\n", conn.Origin, conn.User.Id, done)
			done = true
		case up := <-conn.In:
			//log.Printf("TRACE [%s] user[%d] is receiving update[%s]\n", conn.Origin, conn.User.Id, up.Event)
			conn.SendUpdates(up, conn.User.Id)
			//case <-ticker.C:
			//log.Printf("TRACE [%s] user[%d] still polling updates\n", conn.Origin, conn.User.Id)
		}
	}
	log.Printf("TRACE [%s] user[%d] polling updates stopped\n", conn.Origin, conn.User.Id)
	return done
}
