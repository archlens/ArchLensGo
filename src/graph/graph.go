package graph

import (
	"errors"
	"maps"
	"slices"
)

type Graph struct {
	nodes      map[string]*node
	TotalNodes   uint
	TotalEdges uint
}

type node struct {
	name       string
	incoming  uint
	outgoing   uint
	successors map[string]*node // Name of node to node pointer
}

func NewGraph() *Graph {
	return &Graph{
		nodes: map[string]*node{},
		TotalNodes: 0,
		TotalEdges: 0,
	}
}

func newNode(name string) *node {
	return &node{
		name: name,
		incoming: 0,
		outgoing: 0,
		successors: map[string]*node{},
	}
}

func (g *Graph) GetNodeNames() []string {
	return slices.Collect(maps.Keys(g.nodes))
}

// Add nodes to graph. Nodes that already are in the graph are ignored
func (g *Graph) AddNodes(nodes ...string) {
	for _, nodeName := range nodes {
		if _, ok := g.nodes[nodeName]; !ok {
			g.nodes[nodeName] = newNode(nodeName)
			g.TotalNodes++
		}
	}
}

// Removes a given node from the graph recursively walking all children to make sure it is no longer in the graph
// Function doesn't return error as removing non existing nodes changes nothing so this function shouldn't be able to fail
func (g *Graph) RemoveNodes(nodes ...string) {
	for _, nodeName := range nodes {
		n, ok := g.nodes[nodeName]
		if !ok {
			continue
		}

		// Walk children to remove imports
		for _, children := range n.successors {
			children.incoming--
			g.TotalEdges--
		}

		// Walk graph to find parent nodes and remove self from them
		// This part is very slow and it might be possible to make it faster if we build and maintain the graph differently
		// At the same time this should be a rare operation so maybe it doesn't matter too much
		for _, parent := range g.nodes {
			if _, ok := parent.successors[nodeName]; ok {
				delete(parent.successors, nodeName)
				parent.outgoing--
				g.TotalEdges--
			}
		}

		// Finally remove node from graph
		delete(g.nodes, nodeName)
		g.TotalNodes--
	}
}

// Adds nodes as predecessors to the source node. Predecessor nodes that
// aren't already in the graph will be added.
//
// Returns an error if the source node doesn't exist in the graph
func (g *Graph) AddPredecessors(source string, predecessors ...string) error {
	var sourceNode *node
	var ok bool
	if sourceNode, ok = g.nodes[source]; !ok {
		return errors.New("source node is not in graph")
	}

	var preNode *node
	for _, predecessorNode := range predecessors {
		if preNode, ok = g.nodes[predecessorNode]; !ok {
			// Create and add node if not in graph
			preNode = newNode(predecessorNode)
			g.nodes[predecessorNode] = preNode
			g.TotalNodes++
		}

		if _, ok = preNode.successors[source]; !ok {
			preNode.successors[source] = sourceNode
			preNode.outgoing++
			sourceNode.incoming++
			g.TotalEdges++
		}
	}
	return nil
}

// Adds nodes as successors to the source node. Successor nodes that
// aren't already in the graph will be added.
//
// Returns an error if the source node doesn't exist in the graph
func (g *Graph) AddSuccessors(source string, successors ...string) error {
	var sourceNode *node
	var ok bool
	if sourceNode, ok = g.nodes[source]; !ok {
		return errors.New("source node is not in graph")
	}
	var sucNode *node
	for _, successorNode := range successors {
		if sucNode, ok = g.nodes[successorNode]; !ok {
			// Create and add node if not in graph
			sucNode = newNode(successorNode)
			g.nodes[successorNode] = sucNode
			g.TotalNodes++
		}

		if _, ok = sourceNode.successors[successorNode]; !ok {
			sourceNode.successors[successorNode] = sucNode
			sucNode.incoming++
			sourceNode.outgoing++
			g.TotalEdges++
		}
	}
	return nil
}
