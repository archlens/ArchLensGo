package input

import (
	"encoding/json"
	"os"
)

type View struct {
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
} 

type Github struct {
	Url string `json:"url"`
	Branch string `json:"branch"`
}

type Input struct {
	Name string `json:"name"`
	Github Github `json:"github"`
	RootFolder string `json:"rootFolder"`
	Views map[string]View `json:"views"`
	SaveLocation string `json:"saveLocation"`
}

func Load(path string) (*Input, error) {
	var input *Input

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