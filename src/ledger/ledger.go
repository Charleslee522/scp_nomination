package ledger

import (
	. "github.com/Charleslee522/scp_nomination/src/common"
)

type Ledger struct {
	Node        Node
	N_validator int
	Consensus   Consensus
}

func NewLedger(node Node, validators []Node, quorumThreshold int) *Ledger {
	p := new(Ledger)
	p.Node = node
	p.N_validator = quorumThreshold
	p.Consensus = NewConsensus(node.Name, node.Priority, quorumThreshold, validators)
	return p
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
	l.Consensus.Start()
}
