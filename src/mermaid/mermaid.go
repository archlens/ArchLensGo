package mermaid

import (
	"cmp"
	"path"
	"path/filepath"
	"slices"
	"strconv"

	"github.com/TyphonHill/go-mermaid/diagrams/flowchart"
	"github.com/archlens/ArchLens/graph"
)

// labels maps each package to its display label: the last path segment, or the
// full path when two packages would otherwise look the same.
func labels(names []string) map[string]string {
	count := make(map[string]int, len(names))
	for _, n := range names {
		count[path.Base(n)]++
	}
	out := make(map[string]string, len(names))
	for _, n := range names {
		switch {
		case n == ".":
			out[n] = "(root)"
		case count[path.Base(n)] > 1:
			out[n] = n // ambiguous, keep the full path
		default:
			out[n] = path.Base(n)
		}
	}
	return out
}

// FromGraph converts a dependency graph into a flowchart (edge = "depends on").
func FromGraph(g *graph.Graph, title string) *flowchart.Flowchart {
	fc := flowchart.NewFlowchart()
	fc.Title = title
	fc.SetDirection(flowchart.FlowchartDirectionLeftRight)

	names := g.GetNodeNames()
	slices.Sort(names)
	lbl := labels(names)

	nodes := make(map[string]*flowchart.Node, len(names))
	for _, name := range names {
		nodes[name] = fc.AddNode(lbl[name])
	}

	// Sort edges to ensure same diagram gets generated when using same data.
	edges := g.Edges()
	slices.SortFunc(edges, func(a, b graph.Edge) int {
		return cmp.Or(cmp.Compare(a.From, b.From), cmp.Compare(a.To, b.To))
	})
	for _, e := range edges {
		link := fc.AddLink(nodes[e.From], nodes[e.To])
		link.SetText(strconv.FormatUint(uint64(e.Count), 10))
	}
	return fc
}

// Returns the mermaid diagram from graph as a string
func FlowchartString(g *graph.Graph, title string) string {
	return FromGraph(g, title).String()
}

// Render writes the mermaid diagram to a file specified by the archlens.json file
func Render(g *graph.Graph, title string, diagramLocation string, markdown bool) {
	diagram := FromGraph(g, title)
	path := filepath.Join(diagramLocation, title)
	if markdown {
		diagram.EnableMarkdownFence()
		path = path + ".md"
	}
	err := diagram.RenderToFile(path)
	if err != nil {
		panic("diagram could not be saved to location: " + err.Error())
	}
}
