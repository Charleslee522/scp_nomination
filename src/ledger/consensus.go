package ledger

import (
	"log"
	"sync"
	"time"

	. "github.com/Charleslee522/scp_nomination/src/common"
)

type ChannelValueType chan SCPNomination
type ChannelType map[string]ChannelValueType

var channels ChannelType

func init() {
	log.SetPrefix("[Trace] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	channels = make(ChannelType)
}

type MutexQueue struct {
	M    []SCPNomination
	Lock sync.RWMutex
}

type Consensus struct {
	nodeName          string
	nodePriority      string
	leaders           []string
	quorumThreshold   int
	blockingThreshold int
	round             int
	roundClock        time.Time

	votes    *VotingBox
	accepted *VotingBox
	confirm  *VotingBox

	selfMessageState map[Value]FederatedVotingState
	confirmValues    []Value
	Validators       []Node
	channels         ChannelType
	ValuePool        []Value

	msgQueue MutexQueue
	isInTest bool
}

func NewConsensus(nName string, qTh int,
	validators []Node) Consensus {
	p := Consensus{nodeName: nName, nodePriority: GetPriority(0, nName), quorumThreshold: qTh,
		blockingThreshold: len(validators) + 1 - qTh + 1, round: 1}
	p.votes = NewVotingBox()
	p.accepted = NewVotingBox()
	p.confirm = NewVotingBox()
	p.selfMessageState = make(map[Value]FederatedVotingState)
	p.confirmValues = []Value{}
	p.Validators = validators
	p.leaders = append(p.leaders, p.GetRoundLeader())
	channels[nName] = make(ChannelValueType)
	p.channels = channels
	p.msgQueue = MutexQueue{M: []SCPNomination{}}
	p.roundClock = time.Now()

	return p
}

func (c *Consensus) InsertValues(messages []string) {
	for _, msg := range messages {
		c.ValuePool = append(c.ValuePool, Value{msg})
	}
}

func (c *Consensus) Nominate() {
	// if leader
	if c.isLeader(c.nodeName) {
		time.Sleep(100 * time.Microsecond)
		c.AppendVotes(c.ValuePool, c.nodeName)
		c.broadcast()
	}
}

func (c *Consensus) AppendMessage(msg SCPNomination) {
	c.AppendVotes(msg.Votes, msg.NodeName)
	c.AppendAccepted(msg.Accepted, msg.NodeName)
}

func (c *Consensus) AppendVotes(values []Value, nodeName string) {
	for _, value := range values {
		if c.accepted.HasValue(value) ||
			c.confirm.HasValue(value) {
			log.Println(c.nodeName, "has value", value)
			continue
		}

		c.votes.Add(value, nodeName)
		if c.selfMessageState[value] == NONE {
			c.selfMessageState[value] = VOTES
		}

		if c.votes.Count(value) >= c.quorumThreshold {
			log.Println(c.nodeName, "value in votes", value, "exceed quorum threshold",
				c.quorumThreshold, ", so it is moved to accepted")
			c.accepted.Add(value, c.nodeName)
			c.selfMessageState[value] = ACCEPTED
			c.broadcast()
			c.logSelfState()
		}
	}
}

func (c *Consensus) logSelfState() {
	log.Println(c.nodeName, "self messages states are ", c.selfMessageState)
}

func (c *Consensus) AppendAccepted(values []Value, nodeName string) {
	for _, value := range values {
		if c.confirm.HasValue(value) {
			continue
		}

		c.accepted.Add(value, nodeName)

		if c.selfMessageState[value] < ACCEPTED && c.accepted.Count(value) >= c.blockingThreshold {
			log.Println(c.nodeName, "value in accepted", value, " exceed blocking threshold",
				c.quorumThreshold, ", so it is moved to accept")
			c.accepted.Add(value, c.nodeName)
			c.selfMessageState[value] = ACCEPTED
			c.broadcast()
			c.logSelfState()
		}

		if c.accepted.Count(value) >= c.quorumThreshold {
			log.Println(c.nodeName, "value in accepted", value, " exceed quorum threshold",
				c.quorumThreshold, ", so it is moved to confirm")
			c.confirm.Add(value, c.nodeName)
			c.selfMessageState[value] = CONFIRM
			c.confirmValues = append(c.confirmValues, value)
			c.broadcast()
			c.logSelfState()
		}
	}
}

func (c *Consensus) echo() {
	if c.isInTest {
		return
	}
	votes := []Value{}
	accepted := []Value{}
	for value, state := range c.selfMessageState {
		if state == VOTES || state == ACCEPTED {
			votes = append(votes, value)
		} else {
			// do nothing
		}
	}
	msg := SCPNomination{Votes: votes, Accepted: accepted, NodeName: c.nodeName}
	for _, node := range c.Validators {
		if node.Name == c.nodeName {
			continue
		}
		c.sendMessage(msg, node.Name)
	}
}

func (c *Consensus) sendMessage(msg SCPNomination, toNodeName string) {
	log.Println(c.nodeName, " to ", toNodeName, "send message", msg)
	c.channels[toNodeName] <- msg
}

func (c *Consensus) broadcast() {
	if c.isInTest {
		return
	}
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
	if len(votes)+len(accepted) == 0 {
		return
	}
	msg := SCPNomination{Votes: votes, Accepted: accepted, NodeName: c.nodeName}
	for _, node := range c.Validators {
		if node.Name == c.nodeName {
			continue
		}
		c.sendMessage(msg, node.Name)
	}
}

func (c *Consensus) GetRoundLeader() string {
	maxPriority := GetPriority(c.round, c.nodeName)
	leaderNodeName := c.nodeName
	for _, node := range c.Validators {
		priority := GetPriority(c.round, node.Name)
		if maxPriority < priority {
			maxPriority = priority
			leaderNodeName = node.Name
		}
	}
	return leaderNodeName
}

func (c *Consensus) isLeader(nodeName string) bool {
	for _, leaderName := range c.leaders {
		if leaderName == nodeName {
			return true
		}
	}
	return false
}

func (c *Consensus) ReceiveMessage(msg SCPNomination) {
	log.Println(c.nodeName, "receive message ", msg)
	if time.Since(c.roundClock) > time.Duration(100*c.round)*time.Millisecond {
		c.round++
	}
	c.AppendMessage(msg)
	if c.isLeader(msg.NodeName) {
		c.echo()
	}
}

func Read(queue *MutexQueue) SCPNomination {
	queue.Lock.RLock()
	defer queue.Lock.RUnlock()
	if len(queue.M) == 0 {
		return SCPNomination{}
	} else {
		msg := queue.M[0]
		queue.M = queue.M[1:]
		return msg
	}
}

func Write(queue *MutexQueue, msg SCPNomination) {
	queue.Lock.Lock()
	defer queue.Lock.Unlock()
	queue.M = append(queue.M, msg)
}

func (c *Consensus) GetConfirmedValues() []Value {
	result := []Value{}
	for value, state := range c.selfMessageState {
		if state == CONFIRM {
			result = append(result, value)
		}
	}
	return result
}

func (c *Consensus) Start() {
	if len(c.votes.Voting) == 0 &&
		len(c.accepted.Voting) == 0 {
		c.Nominate()
	}
	quit := make(chan bool)
	go func() {
		start := time.Now()
		for {
			msg := Read(&c.msgQueue)
			if msg.NodeName == "" {
				log.Println("read time: ", time.Since(start))
				if time.Since(start) > 1000*time.Millisecond {
					quit <- true
					return
				}
				time.Sleep(200 * time.Millisecond)
				continue
			}

			c.ReceiveMessage(msg)
			start = time.Now()
		}
	}()

	for {
		select {
		case msg := <-c.channels[c.nodeName]:
			Write(&c.msgQueue, msg)
		case <-quit:
			log.Println("quit event")
			return
		default:
		}
	}
}
