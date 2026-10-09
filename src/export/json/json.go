// Package json exports graphs and graph diffs in ArchLens' JSON format, the same
// shape the Python implementation writes for `render` and `render-diff-json`:
//
//	{
//	    "title": "...",
//	    "packages": [{"name": "...", "state": "NEUTRAL|CREATED|DELETED"}],
//	    "edges": [{"state": "...", "fromPackage": "...", "toPackage": "...",
//	               "label": "3 (+1)", "relations": [...]}]
//	}
package json

import (
	"cmp"
	stdjson "encoding/json"
	"fmt"
	"slices"
	"strconv"

	"github.com/archlens/ArchLens/graph"
)

// Entity states. The names match the Python EntityState enum.
const (
	stateNeutral = "NEUTRAL"
	stateCreated = "CREATED"
	stateDeleted = "DELETED"
)

type Exporter struct{}

type document struct {
	Title    string        `json:"title"`
	Packages []packageJSON `json:"packages"`
	Edges    []edgeJSON    `json:"edges"`
}

type packageJSON struct {
	Name  string `json:"name"`
	State string `json:"state"`
}

type edgeJSON struct {
	State       string     `json:"state"`
	FromPackage string     `json:"fromPackage"`
	ToPackage   string     `json:"toPackage"`
	Label       string     `json:"label"`
	Relations   []relation `json:"relations"`
}

// relation is a file-level dependency behind a package-level edge.
type relation struct {
	FromFile fileRef `json:"from_file"`
	ToFile   fileRef `json:"to_file"`
}

type fileRef struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func (e Exporter) Extension() string {
	return ".json"
}

// ExportGraph writes a plain graph: every package and edge is NEUTRAL and an
// edge's label is its dependency count.
func (e Exporter) ExportGraph(g *graph.Graph, title string) []byte {
	doc := newDocument(title)

	names := g.GetNodeNames()
	slices.Sort(names)
	for _, n := range names {
		doc.Packages = append(doc.Packages, packageJSON{Name: n, State: stateNeutral})
	}

	edges := g.Edges()
	slices.SortFunc(edges, func(a, b graph.Edge) int {
		return cmp.Or(cmp.Compare(a.From, b.From), cmp.Compare(a.To, b.To))
	})
	for _, ed := range edges {
		doc.Edges = append(doc.Edges, edgeJSON{
			State:       stateNeutral,
			FromPackage: ed.From,
			ToPackage:   ed.To,
			Label:       strconv.FormatUint(uint64(ed.Count), 10),
			Relations:   []relation{},
		})
	}
	return marshal(doc)
}

// ExportDiff writes a diff: added packages are CREATED, removed ones DELETED, and
// edges carry the state and label described on diffEdge.
func (e Exporter) ExportDiff(d *graph.Diff, title string) []byte {
	doc := newDocument(title)

	nodes := slices.Clone(d.Nodes)
	slices.Sort(nodes)
	for _, n := range nodes {
		state := stateNeutral
		if d.Added[n] {
			state = stateCreated
		} else if d.Removed[n] {
			state = stateDeleted
		}
		doc.Packages = append(doc.Packages, packageJSON{Name: n, State: state})
	}

	edges := slices.Clone(d.Edges)
	slices.SortFunc(edges, func(a, b graph.EdgeDiff) int {
		return cmp.Or(cmp.Compare(a.From, b.From), cmp.Compare(a.To, b.To))
	})
	for _, ed := range edges {
		state, label := diffEdge(d, ed)
		doc.Edges = append(doc.Edges, edgeJSON{
			State:       state,
			FromPackage: ed.From,
			ToPackage:   ed.To,
			Label:       label,
			Relations:   []relation{},
		})
	}
	return marshal(doc)
}

// diffEdge returns an edge's state and label, following the Python implementation:
//   - edges out of an added/removed package are CREATED/DELETED with the plain count;
//   - otherwise the label is "<count on compare branch> (+/-<change>)" and the state
//     CREATED or DELETED, or just "<count>" and NEUTRAL when nothing changed.
func diffEdge(d *graph.Diff, e graph.EdgeDiff) (state, label string) {
	switch {
	case d.Added[e.From]:
		return stateCreated, strconv.FormatUint(uint64(e.After), 10)
	case d.Removed[e.From]:
		return stateDeleted, strconv.FormatUint(uint64(e.Before), 10)
	}
	switch delta := e.Delta(); {
	case delta > 0:
		return stateCreated, fmt.Sprintf("%d (+%d)", e.After, delta)
	case delta < 0:
		return stateDeleted, fmt.Sprintf("%d (-%d)", e.After, -delta)
	default:
		return stateNeutral, strconv.FormatUint(uint64(e.After), 10)
	}
}

// newDocument uses non-nil slices so empty lists encode as [] rather than null.
func newDocument(title string) document {
	return document{Title: title, Packages: []packageJSON{}, Edges: []edgeJSON{}}
}

func marshal(doc document) []byte {
	// Plain strings and slices can't fail to encode.
	data, err := stdjson.MarshalIndent(doc, "", "    ")
	if err != nil {
		panic("json export: " + err.Error())
	}
	return data
}