package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// initBackend resources
func initBackend(path string) (a Application, err error) {

	/*
		This function is responsible for making sure the backend environment is correctly configured

		- Check supporting libraries / apps are installed
		- Copy demo sourcecode into working directory
	*/

	// Init Application
	a = Application{
		Path: path,
	}

	fmt.Printf("Setting up resources at: %s\n", path)

	// check prerequisites are installed

	// check for Go itself
	goPath := getGoPath()
	if goPath == "" {
		err = fmt.Errorf("Unable to find GOPATH environment variable, please set one")
		return
	}

	// check for GopherJS
	cmd := getVersionCmd()
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("Unable to find `gopherjs` command\n\nERROR: %s\n\nPlease install using `go get -u github.com/gopherjs/gopherjs`", err)
		return
	}

	// check for Ebiten sourcecode
	ebitenPath := filepath.Join(goPath, "src", ebitenRepo)
	_, err = os.Stat(ebitenPath)
	if err != nil {
		err = fmt.Errorf("Unable to find `ebiten` source\n\nERROR: %s\n\nPlease install using `go install %s ", err, ebitenRepo)
		return
	}

	// create dirs if they don't exist
	defaultCodePath := filepath.Join(path, defaultCodeDir)
	err = os.MkdirAll(defaultCodePath, os.FileMode(0755))
	if err != nil {
		err = fmt.Errorf("Failed to create go source dir: %s", err)
		return
	}

	err = os.MkdirAll(filepath.Join(path, defaultCompileDir), os.FileMode(0755))
	if err != nil {
		err = fmt.Errorf("Failed to create javascript dir: %s", err)
		return
	}

	// check for default source
	fullSourcePath := filepath.Join(defaultCodePath, defaultSourceFile)
	if _, err = os.Stat(fullSourcePath); err != nil {
		// file doesn't exist so create it
		var dst *os.File
		dst, err = os.Create(fullSourcePath)
		if err != nil {
			fmt.Printf("Failed to create source file - %s\n", err)
			err = fmt.Errorf("Failed to create source file - %s", err)
			return
		}
		defer dst.Close()
		buf := bytes.NewBuffer([]byte(demoSrc))
		_, err = io.Copy(dst, buf)
		if err != nil {
			fmt.Printf("Failed to copy compiled cart js to target file - %s\n", err)
			err = fmt.Errorf("Failed to copy compiled cart js to target file - %s", err)
			return
		}
	}
	return
}
