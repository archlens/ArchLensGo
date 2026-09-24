package parsers

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"sync"

	"github.com/archlens/ArchLens/utils"
)

type ASTNode struct {
	Type     string    `json:"type"`
	Value    string    `json:"value,omitempty"`
	Children []ASTNode `json:"children,omitempty"`
	Line     int       `json:"line,omitempty"`
}

type ASTResult struct {
	File string
	AST  *ASTNode
	Err  error
}

// GetASTs concurrently parses each file in `files` (relative to rootDir)
// and returns the collected results. Errors are captured per-file rather
// than aborting the whole batch.
func GetASTs(files []string, rootDir string) []ASTResult {
	var wg sync.WaitGroup
	resultsCh := make(chan ASTResult, len(files))
	// Just to make sure we don't spawn too many threads
	sem := make(chan struct{}, runtime.NumCPU())

	for _, file := range files {
		wg.Add(1)
		sem <- struct{}{}
		go func(f string) {
			defer wg.Done()
			defer func() { <-sem }()

			ast, err := GetAST(f, rootDir)
			resultsCh <- ASTResult{File: f, AST: ast, Err: err}
		}(file)
	}

	wg.Wait()
	close(resultsCh)

	results := make([]ASTResult, 0, len(files))
	for r := range resultsCh {
		results = append(results, r)
	}
	return results
}

func GetAST(filePath string, rootDir string) (*ASTNode, error) {
	restore := utils.WithWorkingDirectory(rootDir)
	defer func() {
		_ = restore()
	}()
	// Might want to make this the default path and add a k-v in archlens.json for ast_parser path in case people wanna put it weird places
	cmd := exec.Command("python3", rootDir+"/arch.py", filePath)

	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("python failed: %s", exitErr.Stderr)
		}
		return nil, err
	}

	var root ASTNode
	if err := json.Unmarshal(out, &root); err != nil {
		return nil, fmt.Errorf("bad json from python: %w", err)
	}
	return &root, nil
}
