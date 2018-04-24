package ledger

import (
	"log"
	"time"

	. "github.com/Charleslee522/scp_nomination/src/common"
)

func init() {
	log.SetPrefix("[Trace] ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}

type Consensus struct {
	nodeName          string
	leaderName        string
	quorumThreshold   int
	blockingThreshold int

	votes    *VotingBox
	accepted *VotingBox
	confirm  *VotingBox

	selfMessageState map[Value]FederatedVotingState
	confirmValues    []Value
	Validators       []Node
	channels         ChannelType
}

func NewConsensus(nName string, n_validator int,
	qTh int, validators []Node, channels *ChannelType) Consensus {
	p := Consensus{nodeName: nName, quorumThreshold: qTh,
		blockingThreshold: n_validator - qTh + 1}
	p.votes = NewVotingBox()
	p.accepted = NewVotingBox()
	p.confirm = NewVotingBox()
	p.selfMessageState = make(map[Value]FederatedVotingState)
	p.confirmValues = []Value{}
	p.Validators = validators
	p.leaderName = p.GetLeaderNodeName()
	return p
}

func (c *Consensus) AppendMessage(msg SCPNomination) {
	c.AppendVotes(msg.Votes, msg.NodeName)
	c.AppendAccepted(msg.Accepted, msg.NodeName)
}

func (c *Consensus) AppendVotes(values []Value, nodeName string) {
	for _, value := range values {
		if c.accepted.HasValue(value) ||
			c.confirm.HasValue(value) {
			continue
		}

		c.votes.Add(value, nodeName)
		if c.selfMessageState[value] == NONE {
			c.selfMessageState[value] = VOTES
		}

		if c.votes.Count(value) >= c.quorumThreshold {
			log.Println(value, "in votes exceed quorum threshold",
				c.quorumThreshold, ", so it is moved to accepted")
			c.accepted.Add(value, c.nodeName)
			c.selfMessageState[value] = ACCEPTED
		}
		go c.broadcast()
	}
}

func (c *Consensus) AppendAccepted(values []Value, nodeName string) {
	for _, value := range values {
		if c.confirm.HasValue(value) {
			continue
		}

		c.accepted.Add(value, nodeName)

		if c.accepted.Count(value) >= c.blockingThreshold {
			log.Println(value, "in accepted exceed blocking threshold",
				c.quorumThreshold, ", so it is moved to accept")
			c.accepted.Add(value, c.nodeName)
			c.selfMessageState[value] = ACCEPTED
		}

		if c.accepted.Count(value) >= c.quorumThreshold {
			log.Println(value, "in accepted exceed quorum threshold",
				c.quorumThreshold, ", so it is moved to confirm")
			c.confirm.Add(value, c.nodeName)
			c.selfMessageState[value] = CONFIRM
			c.confirmValues = append(c.confirmValues, value)
		}
	}
}

func (c *Consensus) echo() {
	msg := SCPNomination{Votes: []Value{}, Accepted: []Value{}, NodeName: c.nodeName}
	for _, node := range c.Validators {
		if node.Name == c.nodeName {
			continue
		}
		c.channels[node.Name] <- msg
	}
}

func (c *Consensus) sendMessage(msg *SCPNomination, toNodeName string) {
	select {
	case c.channels[toNodeName] <- *msg:
		log.Println(c.nodeName, "send message in ", msg, " to ", toNodeName)
		break
	default:
		log.Println("nope")
	}
}

func (c *Consensus) broadcast() {
	log.Println(c.nodeName, "broadcast")
	votes := []Value{}
	accepted := []Value{}
	for value, state := range c.selfMessageState {
		if state == VOTES {
			votes = append(votes, value)
		} else if state == ACCEPTED {
			accepted = append(accepted, value)
		} else {
			// do nothing
		}
	}

	msg := SCPNomination{Votes: votes, Accepted: accepted, NodeName: c.nodeName}
	for _, node := range c.Validators {
		if node.Name == c.nodeName {
			continue
		}
		go c.sendMessage(&msg, node.Name)
		time.Sleep(200 * time.Millisecond)
	}
}

func (c *Consensus) GetLeaderNodeName() string {
	maxPriority := 0
	leaderNodeName := c.nodeName
	for _, node := range c.Validators {
		if maxPriority < node.Priority {
			maxPriority = node.Priority
			leaderNodeName = node.Name
		}
	}
	return leaderNodeName
}

func (c *Consensus) isSelfLeader() bool {
	return c.GetLeaderNodeName() == c.nodeName
}
