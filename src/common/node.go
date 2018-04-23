package common

type Node struct {
	Name string
	// message  []string
	Priority int
}

func (n *Node) GetName() string {
	return n.Name
}
