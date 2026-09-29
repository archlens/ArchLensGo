package graph

type Graph struct {
	Nodes []*Node
}

type Node struct {
	Name string
	Incoming uint
	Outgoing uint
	Successors map[string]*Node // Name of node to node pointer
}

func (g *Graph) GetNodeNames() []string {
	names := make([]string, len(g.Nodes))
	for i, node := range g.Nodes {
		names[i] = node.Name
	}
	return names
}

func (n *Node) AddSuccessors(nodes ...*Node) {
	for _, node := range nodes {
		_, ok := n.Successors[node.Name]
		if !ok { // New node is not already a successor
			node.Incoming++
			n.Outgoing++
			n.Successors[node.Name] = node
		}
	}
}