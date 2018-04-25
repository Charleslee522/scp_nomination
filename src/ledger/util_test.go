package ledger

import (
	"testing"
)

func TestUtilGetPriority(t *testing.T) {
	priority_n0 := GetPriority(1, "n0")
	priority_n1 := GetPriority(1, "n1")
	priority_n2 := GetPriority(1, "n2")
	priority_n3 := GetPriority(1, "n3")
	priority_n4 := GetPriority(1, "n4")

	if priority_n4 <= priority_n1 {
		t.Errorf("%s hash %s < %s hash %s, want n4 > n1", "n0", priority_n4, "n1", priority_n1)
	}
	if priority_n4 <= priority_n2 {
		t.Errorf("%s hash %s < %s hash %s, want n4 > n2", "n0", priority_n4, "n2", priority_n2)
	}
	if priority_n4 <= priority_n3 {
		t.Errorf("%s hash %s < %s hash %s, want n4 > n3", "n0", priority_n4, "n3", priority_n3)
	}
	if priority_n4 <= priority_n0 {
		t.Errorf("%s hash %s < %s hash %s, want n4 > n0", "n0", priority_n4, "n4", priority_n0)
	}
}
func TestUtilGetPriorityNextRound(t *testing.T) {
	priority_n0 := GetPriority(2, "n0")
	priority_n1 := GetPriority(2, "n1")
	priority_n2 := GetPriority(2, "n2")
	priority_n3 := GetPriority(2, "n3")
	priority_n4 := GetPriority(2, "n4")

	if priority_n0 <= priority_n1 {
		t.Errorf("%s hash %s < %s hash %s, want n0 > n1", "n0", priority_n0, "n1", priority_n1)
	}
	if priority_n0 <= priority_n2 {
		t.Errorf("%s hash %s < %s hash %s, want n0 > n1", "n0", priority_n0, "n2", priority_n2)
	}
	if priority_n0 <= priority_n3 {
		t.Errorf("%s hash %s < %s hash %s, want n0 > n1", "n0", priority_n0, "n3", priority_n3)
	}
	if priority_n0 <= priority_n4 {
		t.Errorf("%s hash %s < %s hash %s, want n0 > n1", "n0", priority_n0, "n4", priority_n4)
	}
}
