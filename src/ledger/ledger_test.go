package ledger

import (
	"testing"

	. "github.com/Charleslee522/scp_nomination/src/common"
)

// 노드 세 개 정의(n1, n2, n3)
// 이 테스트의 주인공은 n1
// n2와 n3을 validator 로 지정
// 리더 선출(n1) - TestLedgerSelfLeader
// 	다른 노드의 메시지를 받으면 저장만 하고, echoing 하지 않음
// 리더 선출(n2) - TestLedgerOtherLeader
// 	n2의 메시지를 받으면 저장 후에 echoing 함
// 	n3의 메시지를 받으면 저장만 하고, echoing 하지 않음

func TestLedgerFederatedVoting(t *testing.T) {
	var node1 Node = Node{Name: "n1", Priority: 3}
	var node2 Node = Node{Name: "n2", Priority: 2}
	var node3 Node = Node{Name: "n3", Priority: 1}
	var node4 Node = Node{Name: "n4", Priority: 2}
	var node5 Node = Node{Name: "n5", Priority: 1}

	nodes := []Node{node1, node2, node3, node4, node5}

	v11 := Value{Data: "value11"}
	v12 := Value{Data: "value12"}

	var ledger1 *Ledger = NewLedger(node1, nodes, 4)
	vPool1 := []Value{v11, v12}

	ledger1.InsertValues(vPool1)
	ledger1.Nominate()

	msgFrom2 := SCPNomination{Votes: vPool1, NodeName: node2.GetName()}
	msgFrom3 := SCPNomination{Votes: vPool1, NodeName: node3.GetName()}
	msgFrom4 := SCPNomination{Votes: vPool1, NodeName: node4.GetName()}
	msgFrom5 := SCPNomination{Votes: vPool1, NodeName: node5.GetName()}

	ledger1.ReceiveMessage(msgFrom2)
	ledger1.ReceiveMessage(msgFrom3)
	ledger1.ReceiveMessage(msgFrom4)

	if ledger1.GetValueState(v11) != ACCEPTED {
		t.Errorf("v11 State == %q, want %q", ledger1.GetValueState(v11), ACCEPTED)
	}
	if ledger1.GetValueState(v12) != ACCEPTED {
		t.Errorf("v12 State == %q, want %q", ledger1.GetValueState(v12), ACCEPTED)
	}

	ledger1.ReceiveMessage(msgFrom5) // do nothing

	accpetedMsgFrom2 := SCPNomination{Accepted: vPool1, NodeName: node2.GetName()}
	accpetedMsgFrom3 := SCPNomination{Accepted: vPool1, NodeName: node3.GetName()}
	accpetedMsgFrom4 := SCPNomination{Accepted: vPool1, NodeName: node4.GetName()}
	accpetedMsgFrom5 := SCPNomination{Accepted: vPool1, NodeName: node5.GetName()}

	ledger1.ReceiveMessage(accpetedMsgFrom2)
	ledger1.ReceiveMessage(accpetedMsgFrom3)
	ledger1.ReceiveMessage(accpetedMsgFrom4)

	if ledger1.GetValueState(v11) != CONFIRM {
		t.Errorf("v11 State == %q, want %q", ledger1.GetValueState(v11), ACCEPTED)
	}
	if ledger1.GetValueState(v12) != CONFIRM {
		t.Errorf("v12 State == %q, want %q", ledger1.GetValueState(v12), ACCEPTED)
	}

	ledger1.ReceiveMessage(accpetedMsgFrom5) // do nothing
}
