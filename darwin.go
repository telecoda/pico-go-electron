// +build darwin

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	gopherJS = "gopherjs"
)

func getDestDir() (string,error) {
	wd,err := os.Executable()
	if err != nil {
		err = fmt.Errorf("Failed to get executable details: %s", err)
		return "",err
	}

	destDir := filepath.Join(wd, "../resources/app/dynamic/js")

	return destDir, nil
}