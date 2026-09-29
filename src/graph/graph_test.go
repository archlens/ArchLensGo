package graph

import (
	"testing"
)

func Test_Graph(t *testing.T) {
	graph := NewGraph()
	if len(graph.nodes) != 0 {
		t.Errorf("Number of nodes in empty graph is not 0\nExpected:\t%d\nActual\t%d\n", 0, len(graph.nodes))
	}
	nodes := []string{"1", "2", "3", "4", "5", "6"}
	graph.AddNodes(nodes...)
	if len(graph.nodes) != len(nodes) {
		t.Errorf("Incorrect number of nodes in graph\nExpected:\t%d\nActual\t%d\n", len(nodes), len(graph.nodes))
	}
	if len(graph.nodes) != int(graph.TotalNodes) {
		t.Errorf("TotalNodes is not the same as length of graph.nodes\nTotalNodes:\t%d\nActual:\t%d\n", graph.TotalNodes, len(graph.nodes))
	}
	
	oneSuccessors := []string{"2", "3", "7"}
	graph.AddSuccessors("1", oneSuccessors...)
	if len(graph.nodes) != len(nodes)+1 {
		t.Errorf("Incorrect number of nodes in graph. AddSuccessor might not have added the new node\nExpected:\t%d\nActual\t%d\n", len(nodes)+1, len(graph.nodes))
	}
	if int(graph.TotalEdges) != len(oneSuccessors) {
		t.Errorf("Incorrect number of edges\nExpected:\t%d\nActual:\t%d\n", len(oneSuccessors), graph.TotalEdges)
	}

	fivePredecessors := []string{"2", "3", "8"}
	graph.AddSuccessors("5", fivePredecessors...)
	if len(graph.nodes) != len(nodes)+2 {
		t.Errorf("Incorrect number of nodes in graph. AddPredecessors might not have added the new node\nExpected:\t%d\nActual\t%d\n", len(nodes)+1, len(graph.nodes))
	}
	if int(graph.TotalEdges) != len(oneSuccessors) + len(fivePredecessors) {
		t.Errorf("Incorrect number of edges\nExpected:\t%d\nActual:\t%d\n", len(oneSuccessors) + len(fivePredecessors), graph.TotalEdges)
	}
}
