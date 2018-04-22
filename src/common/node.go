package common

type Node struct {
	name string
	// message  []string
	priority int
}

func (n Node) getName() string {
	return n.name
}
