package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// load - loads sourcecode from a specific path
func load(path string) (a Application, err error) {

	// If doesn't end with a filename
	// look in default location
	if !strings.HasSuffix(path, ".go") {
		// this must be userData path..
		path = filepath.Join(path, defaultCodeDir, defaultSourceFile)
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

// save - saves sourcecode to path
func save(source SourceCode) (a Application, err error) {

	// If doesn't end with a filename
	if !strings.HasSuffix(source.Path, ".go") {
		err = fmt.Errorf("Path %s is not a valid filename, MUST end with .go extension", source.Path)
	}

	// write sourcecode to file
	f, err := os.Create(source.Path)
	if err != nil {
		err = fmt.Errorf("Failed to create file - %s - %s\n", source.Path, err)
		return
	}
	defer f.Close()

	_, err = f.Write([]byte(source.Source))
	if err != nil {
		err = fmt.Errorf("Failed to write to file - %s - %s\n", source.Path, err)
		return
	}

	a.Path = source.Path

	return
}
