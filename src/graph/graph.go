package graph

import (
	"errors"
	"maps"
	"slices"
)

type Graph struct {
	nodes      map[string]*node
	TotalNodes uint
	TotalEdges uint
}

type node struct {
	name       string
	incoming   uint
	outgoing   uint
	successors map[string]*node // Name of node to node pointer
	weights    map[string]uint
}

// Edge is a directed edge: From depends on To.
type Edge struct {
	From  string
	To    string
	Count uint
}

func NewGraph() *Graph {
	return &Graph{
		nodes:      map[string]*node{},
		TotalNodes: 0,
		TotalEdges: 0,
	}
}

func newNode(name string) *node {
	return &node{
		name:       name,
		incoming:   0,
		outgoing:   0,
		successors: map[string]*node{},
		weights:    map[string]uint{},
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

		// Outgoing edges (n -> child): children lose an incoming edge.
		// A self-loop is counted here, so the parent walk below skips n.
		for _, child := range n.successors {
			child.incoming--
			g.TotalEdges--
		}

		// Incoming edges (parent -> n): remove n from every parent.
		// This scans the whole graph, which is fine for a rare operation.
		for _, parent := range g.nodes {
			if parent == n {
				continue
			}
			if _, ok := parent.successors[nodeName]; ok {
				delete(parent.successors, nodeName)
				delete(parent.weights, nodeName)
				parent.outgoing--
				g.TotalEdges--
			}
		}

		delete(g.nodes, nodeName)
		g.TotalNodes--
	}
}

// Edges returns every edge in the graph.
func (g *Graph) Edges() []Edge {
	edges := make([]Edge, 0, g.TotalEdges)
	for name, n := range g.nodes {
		for succ := range n.successors {
			edges = append(edges, Edge{From: name, To: succ, Count: n.weights[succ]})
		}
	}
	return edges
}

// AddPackage records that one file in pkg imports deps. Call it once per parsed
// file: nodes are created as needed and each pkg -> dep edge weight goes up by 1.
func (g *Graph) AddPackage(pkg string, deps ...string) {
	g.AddNodes(pkg)

	seen := make(map[string]struct{}, len(deps))
	filtered := make([]string, 0, len(deps))
	for _, d := range deps {
		if d == pkg {
			// No loops
			continue
		}
		if _, dup := seen[d]; dup {
			continue
		}
		seen[d] = struct{}{}
		filtered = append(filtered, d)
	}
	_ = g.AddSuccessors(pkg, filtered...)
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
		preNode.weights[source]++
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
		sourceNode.weights[successorNode]++
	}
	return nil
}
