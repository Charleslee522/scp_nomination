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
var roundChannels ChannelType

const Delay time.Duration = 100

func init() {
	log.SetPrefix("[Trace] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	channels = make(ChannelType)
	roundChannels = make(ChannelType)
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
	Round             int
	roundTime         time.Time
	quit              chan bool

	votes    *VotingBox
	accepted *VotingBox
	confirm  *VotingBox

	selfMessageState map[Value]FederatedVotingState
	confirmValues    []Value
	Validators       []Node
	channels         ChannelType
	roundChannels    ChannelType
	ValuePool        []Value

	msgQueue MutexQueue

	faultyRound []int
	isFaulty    bool

	isInTest bool
}

func NewConsensus(nName string, qTh int, validators []Node) Consensus {
	p := Consensus{nodeName: nName, nodePriority: GetPriority(0, nName), quorumThreshold: qTh,
		blockingThreshold: len(validators) + 1 - qTh + 1, Round: 1}
	p.votes = NewVotingBox()
	p.accepted = NewVotingBox()
	p.confirm = NewVotingBox()
	p.selfMessageState = make(map[Value]FederatedVotingState)
	p.confirmValues = []Value{}
	p.Validators = validators
	p.leaders = append(p.leaders, p.GetRoundLeader())
	channels[nName] = make(ChannelValueType)
	p.channels = channels
	roundChannels[nName] = make(ChannelValueType)
	p.roundChannels = roundChannels
	p.msgQueue = MutexQueue{M: []SCPNomination{}}
	p.roundTime = time.Now()
	p.RoundReset()
	p.quit = make(chan bool)

	return p
}
func NewFaultyConsensus(nName string, qTh int, validators []Node, faultyRound []int) Consensus {
	p := NewConsensus(nName, qTh, validators)
	p.faultyRound = faultyRound
	p.isFaulty = true

	return p
}

func (c *Consensus) RoundReset() {
	// c.Round++
	// log.Println(c.nodeName, "Round change to", c.Round)
	// c.leaders = append(c.leaders, c.GetRoundLeader())
	// log.Println(c.nodeName, "add Round leader ", c.GetRoundLeader())
	// c.Nominate()
	// time.AfterFunc(c.GetRoundDuration(), func() {
	// 	// msg := SCPNomination{NodeName: "ROUND"}
	// 	// c.channels[c.nodeName] <- msg
	// })
}

func (c *Consensus) GetRoundDuration() time.Duration {
	return time.Duration(c.Round+1) * Delay * time.Millisecond
}

func (c *Consensus) InsertValues(messages []string) {
	for _, msg := range messages {
		c.ValuePool = append(c.ValuePool, Value{msg})
	}
}

func (c *Consensus) Nominate() {
	// if leader
	if c.isLeader(c.nodeName) {
		if c.Round == 1 {
			time.Sleep(100 * time.Microsecond)
		}
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

func (c *Consensus) isOutsider() bool {
	if c.isInTest {
		return true
	}

	if c.isFaulty {
		for _, faultyRound := range c.faultyRound {
			if c.Round == faultyRound {
				return true
			}
		}
	}
	return false
}

func (c *Consensus) echo() {
	log.Println(c.nodeName, " echo")
	if c.isOutsider() {
		log.Println(c.nodeName, " in testing or no voting Round, so do not echo")
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
	log.Println(c.nodeName, " broadcast")
	if c.isOutsider() {
		log.Println(c.nodeName, " in testing or no voting Round, so do not broadcast")
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
	maxPriority := GetPriority(c.Round, c.nodeName)
	leaderNodeName := c.nodeName
	for _, node := range c.Validators {
		priority := GetPriority(c.Round, node.Name)
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

func PreWrite(queue *MutexQueue, msg SCPNomination) {
	queue.Lock.Lock()
	defer queue.Lock.Unlock()
	queue.M = append([]SCPNomination{msg}, queue.M...)
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

	go func() {
		start := time.Now()
		for {
			msg := Read(&c.msgQueue)
			if msg.NodeName == "" {
				log.Println("read time: ", time.Since(start))
				if time.Since(start) > 1000*time.Millisecond {
					c.quit <- true
					return
				}
				time.Sleep(200 * time.Millisecond)
			} else if msg.NodeName == "ROUND" {
				log.Println(c.nodeName, "New Round Message Receive")
				c.RoundReset()
			} else {
				c.ReceiveMessage(msg)
				start = time.Now()
			}
		}
	}()

	for {
		select {
		case msg := <-c.channels[c.nodeName]:
			Write(&c.msgQueue, msg)
		case msg := <-c.roundChannels[c.nodeName]:
			PreWrite(&c.msgQueue, msg)
		case <-c.quit:
			log.Println("quit Start()")
			return
		default:
		}
	}
}
