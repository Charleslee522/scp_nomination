package ledger

import (
	"log"

	. "github.com/Charleslee522/scp_nomination/src/common"
)

type Ledger struct {
	Node      Node
	Consensus *Consensus
}

func NewLedger(node Node, validators []Node, quorumThresholdPercent int) *Ledger {
	p := new(Ledger)
	p.Node = node
	n_validator := (len(validators) + 1) * quorumThresholdPercent / 100
	log.Println(node.Name, "has", n_validator, "validators ")
	if node.Kind == 0 {
		p.Consensus = NewConsensus(node.Name, n_validator, validators)
	} else if node.Kind == 1 {
		p.Consensus = NewFaultyConsensus(node.Name, n_validator, validators, node.FaultyRound)
	} else {
		p.Consensus = new(Consensus)
	}
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
