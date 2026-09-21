package parsers

import (
	"fmt"
	"unsafe"

	ts "github.com/tree-sitter/go-tree-sitter"
)

func Parse(rawLang unsafe.Pointer, code []byte) {
	language := ts.NewLanguage(rawLang)

	parser := ts.NewParser()
	defer parser.Close()
	parser.SetLanguage(language)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	// tree and code are both still valid here — that's the window
	// in which you can call ExtractImports.
	imports := ExtractImports(language, tree, code)
	fmt.Println(imports)
}

func ExtractImports(language *ts.Language, tree *ts.Tree, code []byte) []string {
	// Python example — adjust the pattern per language's grammar
	queryStr := `
	(import_statement
	name: (dotted_name) @module)

	(import_statement
	name: (aliased_import
		name: (dotted_name) @module))

	(import_from_statement
	module_name: (dotted_name) @module)

	(import_from_statement
	module_name: (relative_import) @module)
	`

	query, err := ts.NewQuery(language, queryStr)
	if err != nil {
		panic(err)
	}
	defer query.Close()

	cursor := ts.NewQueryCursor()
	defer cursor.Close()

	matches := cursor.Matches(query, tree.RootNode(), code)

	fmt.Println(tree.RootNode().ToSexp())

	var imports []string
	for {
		match := matches.Next()
		if match == nil {
			break
		}
		for _, capture := range match.Captures {
			imports = append(imports, capture.Node.Utf8Text(code))
		}
	}

	return imports
}
