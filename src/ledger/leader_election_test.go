package ledger

import (
	"testing"

	. "github.com/Charleslee522/scp_nomination/src/common"
)

func TestLedgerLeaderElection(t *testing.T) {
	var node0 Node = Node{Name: "n0"}
	var node1 Node = Node{Name: "n1"}
	var node2 Node = Node{Name: "n2"}
	var node3 Node = Node{Name: "n3"}
	var node4 Node = Node{Name: "n4"}

	nodes := []Node{node1, node2, node3, node4}

	var ledger0 *Ledger = NewLedger(node0, nodes, 4)
	ledger0.Consensus.isInTest = true
	ledger0.Consensus.Round = 1
	if ledger0.Consensus.GetRoundLeader() != "n4" {
		t.Errorf("ledger0 Leader %s, want %s", ledger0.Consensus.GetRoundLeader(), "n4")
	}

	ledger0.Consensus.Round = 2
	if ledger0.Consensus.GetRoundLeader() != "n0" {
		t.Errorf("ledger0 Leader %s, want %s", ledger0.Consensus.GetRoundLeader(), "n0")
	}

	ledger0.Consensus.Round = 3
	if ledger0.Consensus.GetRoundLeader() != "n3" {
		t.Errorf("ledger0 Leader %s, want %s", ledger0.Consensus.GetRoundLeader(), "n3")
	}

	ledger0.Consensus.Round = 4
	if ledger0.Consensus.GetRoundLeader() != "n4" {
		t.Errorf("ledger0 Leader %s, want %s", ledger0.Consensus.GetRoundLeader(), "n4")
	}

	ledger0.Consensus.Round = 5
	if ledger0.Consensus.GetRoundLeader() != "n4" {
		t.Errorf("ledger0 Leader %s, want %s", ledger0.Consensus.GetRoundLeader(), "n4")
	}

	ledger0.Consensus.Round = 6
	if ledger0.Consensus.GetRoundLeader() != "n0" {
		t.Errorf("ledger0 Leader %s, want %s", ledger0.Consensus.GetRoundLeader(), "n0")
	}

	ledger0.Consensus.Round = 7
	if ledger0.Consensus.GetRoundLeader() != "n4" {
		t.Errorf("ledger0 Leader %s, want %s", ledger0.Consensus.GetRoundLeader(), "n4")
	}

	ledger0.Consensus.Round = 8
	if ledger0.Consensus.GetRoundLeader() != "n0" {
		t.Errorf("ledger0 Leader %s, want %s", ledger0.Consensus.GetRoundLeader(), "n0")
	}

	ledger0.Consensus.Round = 9
	if ledger0.Consensus.GetRoundLeader() != "n2" {
		t.Errorf("ledger0 Leader %s, want %s", ledger0.Consensus.GetRoundLeader(), "n2")
	}
}
func TestLedgerLeaderElection2(t *testing.T) {
	var node0 Node = Node{Name: "n0"}
	var node1 Node = Node{Name: "n1"}
	var node2 Node = Node{Name: "n2"}
	var node3 Node = Node{Name: "n3"}

	var ledger0 *Ledger = NewLedger(node0, []Node{node1, node2}, 3)
	ledger0.Consensus.isInTest = true
	ledger0.Consensus.Round = 1
	if ledger0.Consensus.GetRoundLeader() != "n2" {
		t.Errorf("ledger0 Leader %s, want %s", ledger0.Consensus.GetRoundLeader(), "n2")
	}

	var ledger1 *Ledger = NewLedger(node1, []Node{node0, node2, node3}, 3)
	ledger1.Consensus.isInTest = true
	ledger1.Consensus.Round = 1
	if ledger1.Consensus.GetRoundLeader() != "n2" {
		t.Errorf("ledger1 Leader %s, want %s", ledger1.Consensus.GetRoundLeader(), "n2")
	}

	var ledger3 *Ledger = NewLedger(node3, []Node{node1, node2}, 3)
	ledger3.Consensus.isInTest = true
	ledger3.Consensus.Round = 1
	if ledger3.Consensus.GetRoundLeader() != "n2" {
		t.Errorf("ledger3 Leader %s, want %s", ledger3.Consensus.GetRoundLeader(), "n2")
	}
}

func TestLedgerLeaderElection3(t *testing.T) {
	var node0 Node = Node{Name: "n0"}
	var node1 Node = Node{Name: "n1"}
	var node2 Node = Node{Name: "n2"}
	var node3 Node = Node{Name: "n3"}
	var node4 Node = Node{Name: "n4"}
	var node5 Node = Node{Name: "n5"}

	var ledger0 *Ledger = NewLedger(node0, []Node{node1, node2}, 3)
	ledger0.Consensus.isInTest = true
	ledger0.Consensus.Round = 1
	if ledger0.Consensus.GetRoundLeader() != "n2" {
		t.Errorf("ledger0 Leader %s, want %s", ledger0.Consensus.GetRoundLeader(), "n2")
	}
	ledger0.Consensus.Round = 2
	if ledger0.Consensus.GetRoundLeader() != "n0" {
		t.Errorf("ledger0 Leader %s, want %s", ledger0.Consensus.GetRoundLeader(), "n0")
	}

	var ledger1 *Ledger = NewLedger(node1, []Node{node0, node2, node3, node4}, 3)
	ledger1.Consensus.isInTest = true
	ledger1.Consensus.Round = 1
	if ledger1.Consensus.GetRoundLeader() != "n4" {
		t.Errorf("ledger1 Leader %s, want %s", ledger1.Consensus.GetRoundLeader(), "n4")
	}
	ledger1.Consensus.Round = 2
	if ledger1.Consensus.GetRoundLeader() != "n0" {
		t.Errorf("ledger1 Leader %s, want %s", ledger1.Consensus.GetRoundLeader(), "n0")
	}

	var ledger3 *Ledger = NewLedger(node3, []Node{node1, node2, node4, node5}, 3)
	ledger3.Consensus.isInTest = true
	ledger3.Consensus.Round = 1
	if ledger3.Consensus.GetRoundLeader() != "n5" {
		t.Errorf("ledger3 Leader %s, want %s", ledger3.Consensus.GetRoundLeader(), "n5")
	}
	ledger3.Consensus.Round = 2
	if ledger3.Consensus.GetRoundLeader() != "n1" {
		t.Errorf("ledger3 Leader %s, want %s", ledger3.Consensus.GetRoundLeader(), "n1")
	}

	var ledger5 *Ledger = NewLedger(node5, []Node{node3, node4}, 3)
	ledger5.Consensus.isInTest = true
	ledger5.Consensus.Round = 1
	if ledger5.Consensus.GetRoundLeader() != "n5" {
		t.Errorf("ledger5 Leader %s, want %s", ledger5.Consensus.GetRoundLeader(), "n5")
	}
	ledger5.Consensus.Round = 2
	if ledger5.Consensus.GetRoundLeader() != "n3" {
		t.Errorf("ledger5 Leader %s, want %s", ledger5.Consensus.GetRoundLeader(), "n3")
	}
}

func TestLedgerLeaderElection4(t *testing.T) {
	var node0 Node = Node{Name: "n0"}
	var node1 Node = Node{Name: "n1"}
	var node2 Node = Node{Name: "n2"}
	var node3 Node = Node{Name: "n3"}
	var node4 Node = Node{Name: "n4"}
	var node5 Node = Node{Name: "n5"}

	var ledger0 *Ledger = NewLedger(node0, []Node{node1, node2, node3, node5}, 3)
	ledger0.Consensus.isInTest = true
	ledger0.Consensus.Round = 1
	if ledger0.Consensus.GetRoundLeader() != "n5" {
		t.Errorf("ledger0 Leader %s, want %s", ledger0.Consensus.GetRoundLeader(), "n5")
	}
	ledger0.Consensus.Round = 2
	if ledger0.Consensus.GetRoundLeader() != "n0" {
		t.Errorf("ledger0 Leader %s, want %s", ledger0.Consensus.GetRoundLeader(), "n0")
	}

	var ledger1 *Ledger = NewLedger(node1, []Node{node0, node2, node3, node4}, 3)
	ledger1.Consensus.isInTest = true
	ledger1.Consensus.Round = 1
	if ledger1.Consensus.GetRoundLeader() != "n4" {
		t.Errorf("ledger1 Leader %s, want %s", ledger1.Consensus.GetRoundLeader(), "n4")
	}
	ledger1.Consensus.Round = 2
	if ledger1.Consensus.GetRoundLeader() != "n0" {
		t.Errorf("ledger1 Leader %s, want %s", ledger1.Consensus.GetRoundLeader(), "n0")
	}

	var ledger3 *Ledger = NewLedger(node3, []Node{node1, node2, node0, node5}, 3)
	ledger3.Consensus.isInTest = true
	ledger3.Consensus.Round = 1
	if ledger3.Consensus.GetRoundLeader() != "n5" {
		t.Errorf("ledger3 Leader %s, want %s", ledger3.Consensus.GetRoundLeader(), "n5")
	}
	ledger3.Consensus.Round = 2
	if ledger3.Consensus.GetRoundLeader() != "n0" {
		t.Errorf("ledger3 Leader %s, want %s", ledger3.Consensus.GetRoundLeader(), "n0")
	}

	var ledger4 *Ledger = NewLedger(node4, []Node{node1, node2}, 3)
	ledger4.Consensus.isInTest = true
	ledger4.Consensus.Round = 1
	if ledger4.Consensus.GetRoundLeader() != "n4" {
		t.Errorf("ledger4 Leader %s, want %s", ledger4.Consensus.GetRoundLeader(), "n4")
	}
	ledger4.Consensus.Round = 2
	if ledger4.Consensus.GetRoundLeader() != "n1" {
		t.Errorf("ledger4 Leader %s, want %s", ledger4.Consensus.GetRoundLeader(), "n1")
	}

	var ledger5 *Ledger = NewLedger(node5, []Node{node0, node3}, 3)
	ledger5.Consensus.isInTest = true
	ledger5.Consensus.Round = 1
	if ledger5.Consensus.GetRoundLeader() != "n5" {
		t.Errorf("ledger5 Leader %s, want %s", ledger5.Consensus.GetRoundLeader(), "n5")
	}
	ledger5.Consensus.Round = 2
	if ledger5.Consensus.GetRoundLeader() != "n0" {
		t.Errorf("ledger5 Leader %s, want %s", ledger5.Consensus.GetRoundLeader(), "n0")
	}
}
