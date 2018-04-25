package common

type VotingBox struct {
	Voting map[Value]map[string]bool
}

func NewVotingBox() *VotingBox {
	p := new(VotingBox)
	p.Voting = make(map[Value]map[string]bool)
	return p
}

func (v *VotingBox) HasValue(value Value) bool {
	return v.Voting[value] != nil
}

func (v *VotingBox) Count(value Value) int {
	if v.Voting[value] == nil {
		return 0
	} else {
		return len(v.Voting[value])
	}
}

func (v *VotingBox) Add(value Value, nodeName string) {
	if v.Voting[value] == nil {
		v.Voting[value] = make(map[string]bool)
	}
	v.Voting[value][nodeName] = true
}
