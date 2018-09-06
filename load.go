package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// loads sourcecode from a specific path
func load(path string) (a Application, err error) {

	// If no path is provided, use the user's home dir
	if len(path) == 0 {
		path = filepath.Join(".", "resources/app/gosrc/main.go")
	}

	f, err := os.Open(path)
	if err != nil {
		err = fmt.Errorf("Failed to open file: %s", err)
		return
	}

	src, err := ioutil.ReadAll(f)
	if err != nil {
		err = fmt.Errorf("Failed to read file: %s", err)
		return
	}

	// Init Application
	a = Application{
		Source:   string(src),
		Path:     path,
	}

	return
}
