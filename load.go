package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// loads sourcecode from a specific path
func load(path string) (a Application, err error) {

	fmt.Printf("TEMP: path: %s\n", path)
	// If doesn't end with a filename
	// look in default location
	if !strings.HasSuffix(path,".go") {
		// this must be userData path..
		path = filepath.Join(path,defaultCodeDir, defaultSourceFile)
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
		Source: string(src),
		Path:   path,
	}

	return
}
