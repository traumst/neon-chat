package sse

import (
	"context"
	"log"
	"neon-chat/src/state"
	"neon-chat/src/utils"
	"time"
)

// TODO queue mechanism for delta updates
func PollUpdates(ctx context.Context, state *state.State, conn *state.Conn, pollingUserId uint) bool {
	log.Printf("TRACE [%s] Live updates triggered by user[%d]\n", conn.Origin, conn.User.Id)
	done := false
	ticker := time.NewTicker(5 * time.Second)

	for !done && !utils.MaintenanceManager.IsInMaintenance() && ctx.Err() == nil {
		select {
		case <-conn.Reader.Context().Done():
			log.Printf("TRACE [%s] user[%d] disconnects as done:[%t]\n", conn.Origin, conn.User.Id, done)
			return true
		case up := <-conn.In:
			log.Printf("TRACE [%s] user[%d] is receiving update[%s]\n", conn.Origin, conn.User.Id, up.Event)
			conn.SendUpdates(up, conn.User.Id)
		case <-ticker.C:
			//log.Printf("TRACE [%s] user[%d] still polling updates\n", conn.Origin, conn.User.Id)
		}
	}
	log.Printf("TRACE [%s] user[%d] polling updates stopped\n", conn.Origin, conn.User.Id)
	return false
}
