package ledger

import (
	"log"

	. "github.com/Charleslee522/scp_nomination/src/common"
)

type Ledger struct {
	Node         Node
	Validators   []Node
	N_validator  int
	ValuePool    []Value
	ValueHistory History
}

func NewLedger(node Node, validators []Node, quorumThreshold int) *Ledger {
	p := new(Ledger)
	p.Node = node
	p.Validators = validators
	p.N_validator = quorumThreshold
	p.ValueHistory = NewHistory(node.Name, p.GetLeaderNodeName(), len(validators), quorumThreshold)
	return p
}

func (l *Ledger) InsertValues(vPool []Value) {
	l.ValuePool = vPool
}

func (l *Ledger) isSelfLeader() bool {
	return true
}

func (l *Ledger) Nominate() {
	// if leader
	if l.isSelfLeader() {
		l.ValueHistory.AppendVotes(l.ValuePool, l.Node.Name)
		log.Println("[VOTES]: ", l.ValueHistory.votes)
		log.Println("[ACCEPTED]: ", l.ValueHistory.accepted)
		log.Println("[CONFIRM]: ", l.ValueHistory.confirm)
	}
}

func (l *Ledger) ReceiveMessage(msg SCPNomination) {
	l.ValueHistory.AppendMessage(msg)
	log.Println("[VOTES]: ", l.ValueHistory.votes)
	log.Println("[ACCEPTED]: ", l.ValueHistory.accepted)
	log.Println("[CONFIRM]: ", l.ValueHistory.confirm)
}

func (l *Ledger) Echo() {

}

type FederatedVotingState uint16

const (
	NONE FederatedVotingState = 0 + iota
	VOTES
	ACCEPTED
	CONFIRM
)

func (l *Ledger) GetValueState(value Value) FederatedVotingState {
	if l.ValueHistory.selfMessageState[value] == NONE {
		return NONE
	}
	return l.ValueHistory.selfMessageState[value]
}

func (l *Ledger) GetLeaderNodeName() string {
	maxPriority := 0
	leaderNodeName := l.Node.Name
	for _, node := range l.Validators {
		if maxPriority < node.Priority {
			maxPriority = node.Priority
			leaderNodeName = node.Name
		}
	}
	return leaderNodeName
}

func (l *Ledger) GetNominatedValues() []Value {
	return l.ValueHistory.confirmValues
}
