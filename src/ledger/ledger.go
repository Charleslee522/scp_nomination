package ledger

import (
	"fmt"

	. "github.com/Charleslee522/scp_nomination/src/common"
)

type Ledger struct {
	Node         Node
	Validators   []Node
	N_Validator  uint32
	ValuePool    []Value
	ValueHistory []Value
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
		l.ValueHistory = append(l.ValueHistory, l.ValuePool...)
	}
}

func (l *Ledger) ReceiveMessage(msg SCPNomination) {
	fmt.Println(l.Node.GetName())
	fmt.Println(msg)
}

type FederatedVotingState int

const (
	Votes FederatedVotingState = 1 + iota
	Accepted
)

func (l *Ledger) GetValueState(value *Value) FederatedVotingState {
	return Accepted
}
