package input

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/archlens/ArchLens/utils"
)

type View struct {
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
}

type Github struct {
	Url    string `json:"url"`
	Branch string `json:"branch"`
}

type Input struct {
	Name         string          `json:"name"`
	Github       Github          `json:"github"`
	RootFolder   string          `json:"rootFolder"`
	Views        map[string]View `json:"views"`
	SaveLocation string          `json:"saveLocation"`
}

func Load(path string) (*Input, error) {
	input := &Input{}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, input)
	if err != nil {
		return nil, err
	}
	return input, nil
}

func matchFiles(patterns []string) (map[string]struct{}, error) {
	files := make(map[string]struct{})
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return files, err
		}
		for _, match := range matches {
			files[match] = struct{}{}
		}
	}
	return files, nil
}

func GetFiles(view *View, rootDir string) ([]string, error) {
	restore := utils.WithWorkingDirectory(rootDir)
	defer restore()

	include, err := matchFiles(view.Include)
	if err != nil {
		return []string{}, err
	}
	exclude, err := matchFiles(view.Exclude)
	if err != nil {
		return []string{}, err
	}
	
	for path, _ := range exclude {
		delete(include, path)
	}
	return slices.Collect(maps.Keys(include)), nil
}

func ReadFiles(files []string) ([][]byte, []error) {
	results := make([][]byte, len(files))
	errs := make([]error, len(files))

	// Using wait group as to synchronize
	var wg sync.WaitGroup
	wg.Add(len(files)) 

	for i, file := range files {
		go func(i int, file string) {
			defer wg.Done()
			data, err := ReadFile(file)
			results[i] = data
			errs[i] = err
		}(i, file)
	}

	wg.Wait()
	return results, errs
}

func ReadFile(file string) ([]byte, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", file, err)
	}
	return data, nil
}