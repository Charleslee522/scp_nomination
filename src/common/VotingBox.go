package common

import "strings"

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

func (v *VotingBox) ContainValueString(value Value) bool {
	for key, _ := range v.Voting {
		if key.Data != value.Data && strings.Contains(key.Data, value.Data) {
			// log.Printf("%s contain %s", key.Data, value.Data)
			return true
		}
	}
	return false
}

func (v *VotingBox) Count(value Value) int {
	if v.Voting[value] == nil {
		return 0
	} else {
		return len(v.Voting[value])
	}
}

func (v *VotingBox) Add(value Value, nodeName string) {
	if !v.ContainValueString(value) {
		if v.Voting[value] == nil {
			v.Voting[value] = make(map[string]bool)
		}
		v.Voting[value][nodeName] = true
	}
}
