package quorum

import (
	"fmt"
	"time"

	"neon-chat/src/consts"
)

type Vote struct {
	NodeId NodeId
	VoteId NodeId
	Yay    bool
	Stamp  time.Time
}

func (v Vote) String() string {
	return fmt.Sprintf("node[%s],vote[%s],yay[%t],%s",
		v.NodeId, v.VoteId, v.Yay, v.Stamp.UTC().Format(consts.Timestamp))
}
