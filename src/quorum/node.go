package quorum

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"neon-chat/src/utils"
)

type NodeId string

type NodeStatus string

const (
	NodeStatusLead = "lead"
	NodeStatusObey = "obey"
)

const voteBufferLength = 64

type Node struct {
	id       NodeId
	status   NodeStatus
	timeout  time.Duration
	election chan Vote
}

func NewNode() (*Node, error) {
	// between 150 and 199
	randomOffset := fmt.Sprintf("%dms", rand.Intn(50)+150)
	timeout, err := time.ParseDuration(randomOffset)
	if err != nil {
		return nil, fmt.Errorf("failed to parse offset[%s], %s", randomOffset, err)
	}

	return &Node{
		id:       NodeId(utils.RandStringBytes(7)),
		status:   NodeStatusObey,
		timeout:  timeout,
		election: make(chan Vote, voteBufferLength),
	}, nil
}

func (n *Node) Id() NodeId {
	return n.id
}

func (n *Node) BecomeLeader() {
	n.status = NodeStatusLead
}

func (n *Node) BecomeFollower() {
	n.status = NodeStatusObey
}

func (n *Node) Status() NodeStatus {
	return n.status
}

func (n *Node) CastVote(candidateId NodeId, yay bool) {
	n.election <- Vote{
		NodeId: n.id,
		VoteId: candidateId,
		Yay:    yay,
		Stamp:  time.Now(),
	}
}

func (n *Node) CountVotes() (NodeId, error) {
	counts := map[NodeId]int{}
	voters := []NodeId{}
	timeout := time.NewTimer(n.timeout)
	done := false
	for !done {
		select {
		case <-timeout.C:
			done = true
		case vote := <-n.election:
			if utils.Contains(voters, vote.NodeId) {
				log.Printf("WARN Node already voted [%s]\n", vote.NodeId)
				continue
			}

			voters = append(voters, vote.NodeId)
			if _, ok := counts[vote.VoteId]; !ok {
				counts[vote.VoteId] = 0
			}
			counts[vote.VoteId] += 1
		}
	}

	elect, count := NodeId(""), 0
	for id, v := range counts {
		if v > count {
			elect, count = id, v
		}
	}
	if count == 0 || elect == "" {
		return "", fmt.Errorf("no votes received")
	}
	if count < (len(voters)/2)+1 {
		return "", fmt.Errorf("no quorum achieved")
	}

	return elect, nil
}
