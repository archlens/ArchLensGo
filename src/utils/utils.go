package utils

import "os"

func WithWorkingDirectory(path string) func() error {
	currentWd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	os.Chdir(path)
	return func() error {
		return os.Chdir(currentWd)
	}
}