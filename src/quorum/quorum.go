package quorum

import (
	"context"
	"fmt"
	"neon-chat/src/utils"
)

type queryType string

const (
	QueryTypeWrite = "write"
)

type channelName string

const (
	ChannelNameRequest  channelName = "request"
	ChannelNameResponse channelName = "response"
)

type QuorumMessage struct {
	kind queryType
}

type Quorum struct {
	ctx  context.Context
	brkr *Broker
	node *Node
}

func NewQuorum(ctx context.Context, brkr *Broker, node *Node) *Quorum {
	return &Quorum{
		ctx:  ctx,
		brkr: brkr,
		node: node,
	}
}

// TODO concurrent requests
func (q *Quorum) RequestWrite(msg map[string]any) (chan QuorumMessage, error) {
	msgId := utils.RandStringBytes(10)
	msg[QueryTypeWrite] = true
	err := q.brkr.Write(q.ctx, string(ChannelNameRequest), msgId, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to request write, %s", err)
	}

	var res []map[string]any
	for {
		res, err = q.brkr.Read(q.ctx, string(ChannelNameResponse), 1, q.node.timeout)
		if err != nil {
			return nil, fmt.Errorf("failed to assert write, %s", err)
		}
		if res == nil || res[0]["msgId"] != msgId {

		}
	}
}
