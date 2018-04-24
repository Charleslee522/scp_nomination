package ledger

import "log"

import . "github.com/Charleslee522/scp_nomination/src/common"

type Ledger struct {
	Node         Node
	Validators   []Node
	N_validator  uint16
	ValuePool    []Value
	ValueHistory History
}

func NewLedger(node Node, validators []Node, n_validator uint16) *Ledger {
	p := new(Ledger)
	p.Node = node
	p.Validators = validators
	p.N_validator = n_validator
	p.ValueHistory = NewHistory(n_validator)
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
		// l.ValueHistory = append(l.ValueHistory, l.ValuePool...)
		l.ValueHistory.AppendVotes(l.ValuePool)
	}
}

func (l *Ledger) ReceiveMessage(msg SCPNomination) {
	l.ValueHistory.AppendMessage(msg)
	log.Println("[VOTES]: ", l.ValueHistory.votes)
	log.Println("[ACCEPT]: ", l.ValueHistory.accepted)
}

func (l *Ledger) Echo() {

}

type FederatedVotingState uint16

const (
	Votes FederatedVotingState = 1 + iota
	Accepted
)

func (l *Ledger) GetValueState(value Value) FederatedVotingState {
	if l.ValueHistory.accepted[value] > 0 {
		return Accepted
	} else {
		return Votes
	}
}
