package ledger

import (
	"testing"

	. "github.com/Charleslee522/scp_nomination/src/common"
)

// func newSCPNomination(data string, node Node) Value {
// 	m := SCPNomination{data, node.GetName()}
// 	return m
// }

// 노드 세 개 정의(n1, n2, n3)
// 이 테스트의 주인공은 n1
// n2와 n3을 validator 로 지정
// 리더 선출(n1) - TestLedgerSelfLeader
// 	다른 노드의 메시지를 받으면 저장만 하고, echoing 하지 않음
// 리더 선출(n2) - TestLedgerOtherLeader
// 	n2의 메시지를 받으면 저장 후에 echoing 함
// 	n3의 메시지를 받으면 저장만 하고, echoing 하지 않음

func TestLedgerSelfLeader(t *testing.T) {
	var node1 Node = Node{Name: "n1", Priority: 3}
	var node2 Node = Node{Name: "n2", Priority: 2}
	var node3 Node = Node{Name: "n3", Priority: 1}

	nodes := []Node{node1, node2, node3}

	v11 := Value{Data: "value11"}
	v12 := Value{Data: "value12"}
	// v21 := Value{"value21"}
	// v22 := Value{"value22"}
	// v31 := Value{"value31"}
	// v32 := Value{"value32"}

	var ledger1 Ledger = Ledger{Node: node1, Validators: nodes, N_Validator: 3}
	// var ledger2 Ledger = Ledger{node2, nodes, 3}
	// var ledger3 Ledger = Ledger{node3, nodes, 3}
	vPool1 := []Value{v11, v12}
	// vPool2 := []Value{v21, v22}
	// vPool3 := []Value{v31, v32}

	ledger1.InsertValues(vPool1)
	ledger1.Nominate()

	msgFrom2 := SCPNomination{}
	msgFrom2.Votes = append(msgFrom2.Votes, vPool1[0])
	msgFrom2.Votes = append(msgFrom2.Votes, vPool1[1])
	msgFrom2.NodeName = node2.GetName()

	msgFrom3 := SCPNomination{}
	msgFrom3.Votes = append(msgFrom3.Votes, vPool1[0])
	msgFrom3.Votes = append(msgFrom3.Votes, vPool1[1])
	msgFrom3.NodeName = node3.GetName()

	ledger1.ReceiveMessage(msgFrom2)
	ledger1.ReceiveMessage(msgFrom3)

	if ledger1.GetValueState(&v11) != Accepted {
		t.Errorf("v11 State == %q, want %q", ledger1.GetValueState(&v11), Accepted)
	}

}
