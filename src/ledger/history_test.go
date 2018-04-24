package ledger

import (
	"testing"
)

func TestHistory(t *testing.T) {
	history := NewHistory("n0", 5, 4)
	if history.quorumThreshold != 4 {
		t.Errorf("history.quorumThreshold == %q, want 4", history.quorumThreshold)
	}
	if history.blockingThreshold != 2 {
		t.Errorf("history.blockingThreshold == %q, want 2", history.blockingThreshold)
	}

	history2 := NewHistory("n0", 5, 5)
	if history2.quorumThreshold != 5 {
		t.Errorf("history2.quorumThreshold == %q, want 4", history.quorumThreshold)
	}
	if history2.blockingThreshold != 1 {
		t.Errorf("history2.blockingThreshold == %q, want 2", history.blockingThreshold)
	}

	history3 := NewHistory("n0", 5, 3)
	if history3.quorumThreshold != 3 {
		t.Errorf("history3.quorumThreshold == %q, want 4", history.quorumThreshold)
	}
	if history3.blockingThreshold != 3 {
		t.Errorf("history3.blockingThreshold == %q, want 2", history.blockingThreshold)
	}
}
