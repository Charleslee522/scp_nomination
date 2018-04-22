package ledger

import (
	"fmt"
	"testing"
)

type Node struct {
	name     string
	priority int
}

func (n Node) GetName() string {
	return n.name
}

type Value struct {
	data string
}

type SCPNomination struct {
	votes    []Value
	accepted []Value
	nodeName string
}

// func newSCPNomination(data string, node Node) Value {
// 	m := SCPNomination{data, node.GetName()}
// 	return m
// }

type Ledger struct {
	node        Node
	validators  []Node
	n_validator uint32
}

func (l Ledger) receiveMessage(msg SCPNomination) {
	fmt.Println(l.node.GetName())
	fmt.Println(msg)
}

type FederatedVotingState int

const (
	Votes FederatedVotingState = 1 + iota
	Accepted
)

func (l Ledger) getValueState(value *Value) FederatedVotingState {
	return Accepted
}

// 노드 세 개 정의(n1, n2, n3)
// 이 테스트의 주인공은 n1
// n2와 n3을 validator 로 지정
// 리더 선출(n1) - TestLedgerSelfLeader
// 	다른 노드의 메시지를 받으면 저장만 하고, echoing 하지 않음
// 리더 선출(n2) - TestLedgerOtherLeader
// 	n2의 메시지를 받으면 저장 후에 echoing 함
// 	n3의 메시지를 받으면 저장만 하고, echoing 하지 않음

func TestLedgerSelfLeader(t *testing.T) {
	var node1 Node = Node{"n1", 3}
	var node2 Node = Node{"n2", 2}
	var node3 Node = Node{"n3", 1}

	nodes := []Node{node1, node2, node3}

	v11 := Value{"value11"}
	v12 := Value{"value12"}
	// v21 := Value{"value21"}
	// v22 := Value{"value22"}
	// v31 := Value{"value31"}
	// v32 := Value{"value32"}

	var ledger1 Ledger = Ledger{node1, nodes, 3}
	// var ledger2 Ledger = Ledger{node2, nodes, 3}
	// var ledger3 Ledger = Ledger{node3, nodes, 3}
	vPool1 := []Value{v11, v12}
	// vPool2 := []Value{v21, v22}
	// vPool3 := []Value{v31, v32}

	msgFrom2 := SCPNomination{}
	msgFrom2.votes = append(msgFrom2.votes, vPool1[0])
	msgFrom2.votes = append(msgFrom2.votes, vPool1[1])
	msgFrom2.nodeName = node2.GetName()

	msgFrom3 := SCPNomination{}
	msgFrom3.votes = append(msgFrom3.votes, vPool1[0])
	msgFrom3.votes = append(msgFrom3.votes, vPool1[1])
	msgFrom3.nodeName = node3.GetName()

	ledger1.receiveMessage(msgFrom2)
	ledger1.receiveMessage(msgFrom3)

	if ledger1.getValueState(&v11) != Accepted {
		t.Errorf("v11 State == %q, want %q", ledger1.getValueState(&v11), Accepted)
	}

}
