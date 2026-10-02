package mermaid

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/TyphonHill/go-mermaid/diagrams/flowchart"
	"github.com/archlens/ArchLens/graph"
)

const (
	colorIncrease = "#2da44e" // green
	colorDecrease = "#cf222e" // red
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

	for _, e := range g.Edges() {
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
func Render(g *graph.Graph, title string, diagramLocation string) {
	err := FromGraph(g, title).RenderToFile(filepath.Join(diagramLocation, title))
	if err != nil {
		panic("diagram could not be saved to location: " + err.Error())
	}
}

// diffLabel is the link text: the current count, plus the change when there is one.
func diffLabel(e graph.EdgeDiff) string {
	switch d := e.Delta(); {
	case d > 0:
		return fmt.Sprintf(`"%d (+%d)"`, e.After, d)
	case d < 0:
		return fmt.Sprintf(`"%d (-%d)"`, e.After, -d)
	default:
		return fmt.Sprintf(`"%d"`, e.After)
	}
}

// DiffString converts a graph diff into a flowchart. Links whose dependency count
// went up are green, links where it went down are red, unchanged links keep the
// default style. The flowchart library has no linkStyle support, so those lines
// are appended by hand (mermaid addresses links by their order of definition).
func DiffString(d *graph.Diff, title string) string {
	fc := flowchart.NewFlowchart()
	fc.Title = title
	fc.SetDirection(flowchart.FlowchartDirectionLeftRight)

	lbl := labels(d.Nodes)
	nodes := make(map[string]*flowchart.Node, len(d.Nodes))
	for _, name := range d.Nodes {
		nodes[name] = fc.AddNode(lbl[name])
	}

	var styles strings.Builder
	for i, e := range d.Edges {
		fc.AddLink(nodes[e.From], nodes[e.To]).SetText(diffLabel(e))

		color := ""
		if delta := e.Delta(); delta > 0 {
			color = colorIncrease
		} else if delta < 0 {
			color = colorDecrease
		}
		if color != "" {
			fmt.Fprintf(&styles, "    linkStyle %d stroke:%s,stroke-width:2px,color:%s\n", i, color, color)
		}
	}
	return fc.String() + styles.String()
}

// RenderDiff writes the diff diagram to diagramLocation/fileName.
func RenderDiff(d *graph.Diff, title, diagramLocation, fileName string) error {
	if err := os.MkdirAll(diagramLocation, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(diagramLocation, fileName), []byte(DiffString(d, title)), 0o644)
}