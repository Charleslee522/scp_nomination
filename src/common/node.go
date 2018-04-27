package common

type Node struct {
	Kind        int
	FaultyRound []int
	Name        string
	Priority    string
	Port        int
	Validators  []string
	Messages    []string
	Threshold   int
}

func (n *Node) GetName() string {
	return n.Name
}
