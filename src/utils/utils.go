package utils

import "os"

func WithWorkingDirectory(path string) func() {
	currentWd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	err = os.Chdir(path)
	if err != nil {
		panic(err)
	}
	return func() {
		if err := os.Chdir(currentWd); err != nil {
			panic(err)
		}
	}
}