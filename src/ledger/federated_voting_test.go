package ledger

import (
	"testing"

	. "github.com/Charleslee522/scp_nomination/src/common"
)

func TestLedgerFederatedVotingBasic(t *testing.T) {
	var node0 Node = Node{Name: "n0", Priority: 3}
	var node1 Node = Node{Name: "n1", Priority: 2}
	var node2 Node = Node{Name: "n2", Priority: 1}
	var node3 Node = Node{Name: "n3", Priority: 2}
	var node4 Node = Node{Name: "n4", Priority: 1}

	nodes := []Node{node1, node2, node3, node4}

	v11 := Value{Data: "value11"}
	v12 := Value{Data: "value12"}

	var ledger1 *Ledger = NewLedger(node0, nodes, 4)
	ledger1.Consensus.isInTest = true
	vPool1 := []Value{v11, v12}

	ledger1.Consensus.InsertValues([]string{"value11", "value12"})
	ledger1.Consensus.Nominate()

	msgFrom1 := SCPNomination{Votes: vPool1, NodeName: node1.GetName()}
	msgFrom2 := SCPNomination{Votes: vPool1, NodeName: node2.GetName()}
	msgFrom3 := SCPNomination{Votes: vPool1, NodeName: node3.GetName()}
	msgFrom4 := SCPNomination{Votes: vPool1, NodeName: node4.GetName()}

	ledger1.Consensus.ReceiveMessage(msgFrom1)
	ledger1.Consensus.ReceiveMessage(msgFrom2)
	ledger1.Consensus.ReceiveMessage(msgFrom3)

	if ledger1.GetValueState(v11) != ACCEPTED {
		t.Errorf("v11 State == %q, want %q", ledger1.GetValueState(v11), ACCEPTED)
	}
	if ledger1.GetValueState(v12) != ACCEPTED {
		t.Errorf("v12 State == %q, want %q", ledger1.GetValueState(v12), ACCEPTED)
	}

	ledger1.Consensus.ReceiveMessage(msgFrom4) // do nothing

	accpetedMsgFrom1 := SCPNomination{Accepted: vPool1, NodeName: node1.GetName()}
	accpetedMsgFrom2 := SCPNomination{Accepted: vPool1, NodeName: node2.GetName()}
	accpetedMsgFrom3 := SCPNomination{Accepted: vPool1, NodeName: node3.GetName()}
	accpetedMsgFrom4 := SCPNomination{Accepted: vPool1, NodeName: node4.GetName()}

	ledger1.Consensus.ReceiveMessage(accpetedMsgFrom1)
	ledger1.Consensus.ReceiveMessage(accpetedMsgFrom2)
	ledger1.Consensus.ReceiveMessage(accpetedMsgFrom3)

	if ledger1.GetValueState(v11) != CONFIRM {
		t.Errorf("v11 State == %q, want %q", ledger1.GetValueState(v11), ACCEPTED)
	}
	if ledger1.GetValueState(v12) != CONFIRM {
		t.Errorf("v12 State == %q, want %q", ledger1.GetValueState(v12), ACCEPTED)
	}

	ledger1.Consensus.ReceiveMessage(accpetedMsgFrom4) // do nothing
}

func TestLedgerFederatedVotingReceiveDuplicatedMessage(t *testing.T) {
	var node0 Node = Node{Name: "n0", Priority: 3}
	var node1 Node = Node{Name: "n1", Priority: 2}
	var node2 Node = Node{Name: "n2", Priority: 1}
	var node3 Node = Node{Name: "n3", Priority: 2}
	var node4 Node = Node{Name: "n4", Priority: 1}

	nodes := []Node{node1, node2, node3, node4}

	v11 := Value{Data: "value11"}
	v12 := Value{Data: "value12"}

	var ledger1 *Ledger = NewLedger(node0, nodes, 4)
	ledger1.Consensus.isInTest = true
	vPool1 := []Value{v11, v12}

	ledger1.Consensus.InsertValues([]string{"value11", "value12"})
	ledger1.Consensus.Nominate()

	msgFrom1 := SCPNomination{Votes: vPool1, NodeName: node1.GetName()}
	msgFrom2 := SCPNomination{Votes: vPool1, NodeName: node2.GetName()}
	msgFrom3 := SCPNomination{Votes: vPool1, NodeName: node3.GetName()}
	msgFrom4 := SCPNomination{Votes: vPool1, NodeName: node4.GetName()}

	ledger1.Consensus.ReceiveMessage(msgFrom1)
	ledger1.Consensus.ReceiveMessage(msgFrom1)
	ledger1.Consensus.ReceiveMessage(msgFrom1)
	ledger1.Consensus.ReceiveMessage(msgFrom1)

	if ledger1.GetValueState(v11) != VOTES {
		t.Errorf("v11 State == %q, want %q", ledger1.GetValueState(v11), VOTES)
	}
	if ledger1.GetValueState(v12) != VOTES {
		t.Errorf("v12 State == %q, want %q", ledger1.GetValueState(v12), VOTES)
	}

	ledger1.Consensus.ReceiveMessage(msgFrom2)
	ledger1.Consensus.ReceiveMessage(msgFrom3)

	if ledger1.GetValueState(v11) != ACCEPTED {
		t.Errorf("v11 State == %q, want %q", ledger1.GetValueState(v11), ACCEPTED)
	}
	if ledger1.GetValueState(v12) != ACCEPTED {
		t.Errorf("v12 State == %q, want %q", ledger1.GetValueState(v12), ACCEPTED)
	}

	ledger1.Consensus.ReceiveMessage(msgFrom4) // do nothing

	accpetedMsgFrom1 := SCPNomination{Accepted: vPool1, NodeName: node1.GetName()}
	accpetedMsgFrom2 := SCPNomination{Accepted: vPool1, NodeName: node2.GetName()}
	accpetedMsgFrom3 := SCPNomination{Accepted: vPool1, NodeName: node3.GetName()}
	accpetedMsgFrom4 := SCPNomination{Accepted: vPool1, NodeName: node4.GetName()}

	ledger1.Consensus.ReceiveMessage(accpetedMsgFrom1)
	ledger1.Consensus.ReceiveMessage(accpetedMsgFrom2)
	ledger1.Consensus.ReceiveMessage(accpetedMsgFrom3)

	if ledger1.GetValueState(v11) != CONFIRM {
		t.Errorf("v11 State == %q, want %q", ledger1.GetValueState(v11), CONFIRM)
	}
	if ledger1.GetValueState(v12) != CONFIRM {
		t.Errorf("v12 State == %q, want %q", ledger1.GetValueState(v12), CONFIRM)
	}

	ledger1.Consensus.ReceiveMessage(accpetedMsgFrom4) // do nothing
}

func TestLedgerFederatedVotingByBlockingThreshold(t *testing.T) {
	var node0 Node = Node{Name: "n0", Priority: 3}
	var node1 Node = Node{Name: "n1", Priority: 2}
	var node2 Node = Node{Name: "n2", Priority: 1}
	var node3 Node = Node{Name: "n3", Priority: 2}
	var node4 Node = Node{Name: "n4", Priority: 1}

	nodes := []Node{node1, node2, node3, node4}

	v11 := Value{Data: "value11"}
	v12 := Value{Data: "value12"}

	var ledger1 *Ledger = NewLedger(node0, nodes, 4)
	ledger1.Consensus.isInTest = true
	vPool1 := []Value{v11, v12}

	ledger1.Consensus.InsertValues([]string{"value11", "value12"})
	ledger1.Consensus.Nominate()

	if ledger1.GetValueState(v11) != VOTES {
		t.Errorf("v11 State == %q, want %q", ledger1.GetValueState(v11), VOTES)
	}
	if ledger1.GetValueState(v12) != VOTES {
		t.Errorf("v12 State == %q, want %q", ledger1.GetValueState(v12), VOTES)
	}

	accpetedMsgFrom1 := SCPNomination{Accepted: vPool1, NodeName: node1.GetName()}
	accpetedMsgFrom2 := SCPNomination{Accepted: vPool1, NodeName: node2.GetName()}
	accpetedMsgFrom3 := SCPNomination{Accepted: vPool1, NodeName: node3.GetName()}
	accpetedMsgFrom4 := SCPNomination{Accepted: vPool1, NodeName: node4.GetName()}

	ledger1.Consensus.ReceiveMessage(accpetedMsgFrom1)
	ledger1.Consensus.ReceiveMessage(accpetedMsgFrom2)

	// v11, v12 self state is changed by blocking threshold
	if ledger1.GetValueState(v11) != ACCEPTED {
		t.Errorf("v11 State == %q, want %q", ledger1.GetValueState(v11), ACCEPTED)
	}
	if ledger1.GetValueState(v12) != ACCEPTED {
		t.Errorf("v12 State == %q, want %q", ledger1.GetValueState(v12), ACCEPTED)
	}

	ledger1.Consensus.ReceiveMessage(accpetedMsgFrom3)

	if ledger1.GetValueState(v11) != CONFIRM {
		t.Errorf("v11 State == %q, want %q", ledger1.GetValueState(v11), CONFIRM)
	}
	if ledger1.GetValueState(v12) != CONFIRM {
		t.Errorf("v12 State == %q, want %q", ledger1.GetValueState(v12), CONFIRM)
	}

	ledger1.Consensus.ReceiveMessage(accpetedMsgFrom4) // do nothing

}
