package ledger

import (
	"log"
	"time"

	. "github.com/Charleslee522/scp_nomination/src/common"
)

type ChannelValueType chan SCPNomination
type ChannelType map[string]ChannelValueType

type Ledger struct {
	Node        Node
	N_validator int
	ValuePool   []Value
	Consensus   Consensus
	channels    ChannelType
}

func NewLedger(node Node, validators []Node, quorumThreshold int, channels *ChannelType) *Ledger {
	p := new(Ledger)
	p.Node = node
	p.N_validator = quorumThreshold
	p.Consensus = NewConsensus(node.Name, len(validators), quorumThreshold, validators, channels)
	p.channels = *channels
	return p
}

func (l *Ledger) InsertValues(vPool []Value) {
	l.ValuePool = vPool
}

func (l *Ledger) Nominate() {
	log.Println(l.Node.Name, "nominate ")
	// if leader
	if l.Consensus.isSelfLeader() {
		l.Consensus.AppendVotes(l.ValuePool, l.Node.Name)
		log.Println("[VOTES]: ", l.Consensus.votes)
		log.Println("[ACCEPTED]: ", l.Consensus.accepted)
		log.Println("[CONFIRM]: ", l.Consensus.confirm)
	}
}

func (l *Ledger) ReceiveMessage(msg SCPNomination) {
	log.Println(l.Node.Name, "receive message ", msg)
	l.Consensus.AppendMessage(msg)
	log.Println("[VOTES]: ", l.Consensus.votes)
	log.Println("[ACCEPTED]: ", l.Consensus.accepted)
	log.Println("[CONFIRM]: ", l.Consensus.confirm)
}

type FederatedVotingState uint16

const (
	NONE FederatedVotingState = 0 + iota
	VOTES
	ACCEPTED
	CONFIRM
)

func (l *Ledger) GetValueState(value Value) FederatedVotingState {
	if l.Consensus.selfMessageState[value] == NONE {
		return NONE
	}
	return l.Consensus.selfMessageState[value]
}

func (l *Ledger) GetNominatedValues() []Value {
	return l.Consensus.confirmValues
}

func (l *Ledger) Start() {
	for {
		log.Println(l.Node.Name, " Node Start! ")
		for {
			select {
			case msg := <-l.channels[l.Node.Name]:
				l.ReceiveMessage(msg)
			default:
			}
			time.Sleep(1 * time.Second)
		}
	}
}
