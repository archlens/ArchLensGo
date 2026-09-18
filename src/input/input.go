package input

import (
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
	"slices"
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

func GetFiles(view *View) ([]string, error) {
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
