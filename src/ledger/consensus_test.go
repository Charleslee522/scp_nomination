package ledger

import (
	"testing"

	. "github.com/Charleslee522/scp_nomination/src/common"
)

func TestConsensus(t *testing.T) {
	nodes := []Node{Node{}, Node{}, Node{}, Node{}}
	consensus := NewConsensus("n0", 1, 4, nodes)
	if consensus.quorumThreshold != 4 {
		t.Errorf("consensus.quorumThreshold == %d, want 4", consensus.quorumThreshold)
	}
	if consensus.blockingThreshold != 2 {
		t.Errorf("consensus.blockingThreshold == %d, want 2", consensus.blockingThreshold)
	}

	consensus2 := NewConsensus("n0", 1, 5, nodes)
	if consensus2.quorumThreshold != 5 {
		t.Errorf("consensus2.quorumThreshold == %d, want 4", consensus.quorumThreshold)
	}
	if consensus2.blockingThreshold != 1 {
		t.Errorf("consensus2.blockingThreshold == %d, want 2", consensus.blockingThreshold)
	}

	consensus3 := NewConsensus("n0", 1, 3, nodes)
	if consensus3.quorumThreshold != 3 {
		t.Errorf("consensus3.quorumThreshold == %d, want 4", consensus.quorumThreshold)
	}
	if consensus3.blockingThreshold != 3 {
		t.Errorf("consensus3.blockingThreshold == %d, want 2", consensus.blockingThreshold)
	}
}
