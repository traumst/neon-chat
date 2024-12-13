package quorum

import (
	"fmt"
	"log"
	"math/rand"
	"neon-chat/src/consts"
	"neon-chat/src/utils"
	"time"
)

type Vote struct {
	NodeId string
	VoteId string
	Yay    bool
	Stamp  time.Time
}

func (v Vote) String() string {
	return fmt.Sprintf("node[%s],vote[%s],yay[%t],%s",
		v.NodeId, v.VoteId, v.Yay, v.Stamp.UTC().Format(consts.Timestamp))
}

type NodeStatus string

const (
	Lead = "lead"
	Obey = "obey"
)

type Node struct {
	id         string
	status     NodeStatus
	timeout    time.Duration
	heartbeat  *time.Timer
	castVote   chan Vote
	countVotes chan Vote
}

func InitNode(
	castVote chan Vote,
	countVotes chan Vote,
) (*Node, error) {
	offs := fmt.Sprintf("1%d%dms", rand.Intn(5)+1, rand.Intn(10))
	offset, err := time.ParseDuration(offs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse offset[%s]", offs)
	}

	return &Node{
		id:         utils.RandStringBytes(7),
		status:     Obey,
		timeout:    offset,
		heartbeat:  time.NewTimer(offset),
		castVote:   castVote,
		countVotes: countVotes,
	}, nil
}

func (n *Node) Status() NodeStatus {
	return n.status
}

func (n *Node) CastVote(candidateId string, yay bool) {
	n.castVote <- Vote{
		NodeId: n.id,
		VoteId: candidateId,
		Yay:    yay,
		Stamp:  time.Now(),
	}
}

func (n *Node) CountVotes() (string, error) {
	counts := map[string]int{}
	voters := []string{}
	timeout := time.NewTimer(n.timeout)
	done := false
	for !done {
		select {
		case <-timeout.C:
			done = true
		case vote := <-n.countVotes:
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

	elect, count := "", 0
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
