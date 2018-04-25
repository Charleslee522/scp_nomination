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
	if ledger0.Consensus.GetLeaderNodeName() != "n0" {
		t.Errorf("ledger0 Leader %s, want %s", ledger0.Consensus.GetLeaderNodeName(), "n0")
	}
}
