package common

type VotingBox struct {
	voting map[Value]map[string]bool
}

func NewVotingBox() *VotingBox {
	p := new(VotingBox)
	p.voting = make(map[Value]map[string]bool)
	return p
}

func (v *VotingBox) HasValue(value Value) bool {
	return v.voting[value] != nil
}

func (v *VotingBox) Count(value Value) int {
	if v.voting[value] == nil {
		return 0
	} else {
		return len(v.voting[value])
	}
}

func (v *VotingBox) Add(value Value, nodeName string) {
	if v.voting[value] == nil {
		v.voting[value] = make(map[string]bool)
	}
	v.voting[value][nodeName] = true
}
