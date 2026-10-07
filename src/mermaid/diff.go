package mermaid

import (
	"cmp"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/TyphonHill/go-mermaid/diagrams/flowchart"
	"github.com/archlens/ArchLens/graph"
)

const (
	colorIncrease = "#2da44e" // green
	colorDecrease = "#cf222e" // red

	fillAdded   = "#dafbe1" // light green node background
	fillRemoved = "#ffebe9" // light red node background
	textOnFill  = "#1f2328" // dark text, readable on both fills
)

// nodeStyle returns the style for an added or removed package, or nil if unchanged.
func nodeStyle(d *graph.Diff, name string) *flowchart.NodeStyle {
	var fill, stroke string
	switch {
	case d.Added[name]:
		fill, stroke = fillAdded, colorIncrease
	case d.Removed[name]:
		fill, stroke = fillRemoved, colorDecrease
	default:
		return nil
	}
	s := flowchart.NewNodeStyle()
	s.Fill, s.Stroke, s.Color, s.StrokeWidth = fill, stroke, textOnFill, 2
	return s
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

// DiffString converts a graph diff into a flowchart. Packages that were added are
// green and packages that were removed are red. Links whose dependency count
// went up are green, links where it went down are red, unchanged links keep the
// default style. The flowchart library has no linkStyle support, so those lines
// are appended by hand (mermaid addresses links by their order of definition).
func DiffString(d *graph.Diff, title string) string {
	fc := flowchart.NewFlowchart()
	fc.Title = title
	fc.SetDirection(flowchart.FlowchartDirectionLeftRight)

	lbl := labels(d.Nodes)
	nodes := make(map[string]*flowchart.Node, len(d.Nodes))
	graphNodes := d.Nodes
	slices.Sort(graphNodes)
	for _, name := range graphNodes {
		nodes[name] = fc.AddNode(lbl[name])
		if st := nodeStyle(d, name); st != nil {
			nodes[name].SetStyle(st)
		}
	}

	var styles strings.Builder
	edges := d.Edges
	slices.SortFunc(edges, func(a, b graph.EdgeDiff) int {
		return cmp.Or(cmp.Compare(a.From, b.From), cmp.Compare(a.To, b.To))
	})
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
func RenderDiff(d *graph.Diff, title, diagramLocation, fileName string, markdown bool) error {
	if err := os.MkdirAll(diagramLocation, 0o755); err != nil {
		return err
	}
	diagram := DiffString(d, title)
	if markdown {
		diagram = "```mermaid\n" + diagram  + "```"
		fileName += ".md"
	}
	return os.WriteFile(filepath.Join(diagramLocation, fileName), []byte(diagram), 0o644)
}