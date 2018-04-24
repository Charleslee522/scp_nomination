package common

type Node struct {
	Name       string
	Priority   int
	Port       int
	Validators []string
}

func (n *Node) GetName() string {
	return n.Name
}
