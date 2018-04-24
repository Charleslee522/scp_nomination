package ledger

import (
	"testing"

	. "github.com/Charleslee522/scp_nomination/src/common"
)

func TestLedgerLeaderSelf(t *testing.T) {
	var node0 Node = Node{Name: "n0", Priority: 3}
	var node1 Node = Node{Name: "n1", Priority: 2}
	var node2 Node = Node{Name: "n2", Priority: 1}
	var node3 Node = Node{Name: "n3", Priority: 2}
	var node4 Node = Node{Name: "n4", Priority: 1}

	nodes := []Node{node0, node1, node2, node3, node4}

	var ledger1 *Ledger = NewLedger(node0, nodes, 4)
	if ledger1.GetLeaderNodeName() != "n0" {
		t.Errorf("Ledger.GetLeaderNodeName() == %q, want n0", ledger1.GetLeaderNodeName())
	}
}

func TestLedgerLeaderTheOther(t *testing.T) {
	var node0 Node = Node{Name: "n0", Priority: 3}
	var node1 Node = Node{Name: "n1", Priority: 2}
	var node2 Node = Node{Name: "n2", Priority: 6}
	var node3 Node = Node{Name: "n3", Priority: 2}
	var node4 Node = Node{Name: "n4", Priority: 1}

	nodes := []Node{node0, node1, node2, node3, node4}

	var ledger1 *Ledger = NewLedger(node0, nodes, 4)
	if ledger1.GetLeaderNodeName() != "n2" {
		t.Errorf("Ledger.GetLeaderNodeName() == %q, want n2", ledger1.GetLeaderNodeName())
	}
}
