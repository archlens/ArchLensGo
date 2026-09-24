package utils

import "os"

func WithWorkingDirectory(path string) func() error {
	currentWd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	err = os.Chdir(path)
	if err != nil {
		panic(err)
	}
	return func() error {
		return os.Chdir(currentWd)
	}
}