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
		// var wd string
		// wd, err = os.Getwd()
		// if err != nil {
		// 	err = fmt.Errorf("Failed to get current working dir: %s", err)
		// 	return
		// }
		//path = filepath.Join(wd,defaultCodeDir,defaultSourceFile)
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
	//  Read source
	// var files []os.FileInfo

	// if files, err = ioutil.ReadDir(path); err != nil {
	// 	return Application{}, err
	// }

	// Init Application
	a = Application{
		Source:   string(src),
		Path:     path,
		CompErrs: make([]CompErr, 0),
	}

	return
}
