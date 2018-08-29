package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// run - compiles and executes latest code
func run(source string) (a Application, err error) {

	// save source to file
	dir, err := ioutil.TempDir("", "example")
	if err != nil {
		err = fmt.Errorf("Failed to create temporary dir - %s", err)
		return
	}

	defer os.RemoveAll(dir) // clean up

	tmpfn := filepath.Join(dir, "main.go")
	if err = ioutil.WriteFile(tmpfn, []byte(source), 0666); err != nil {
		err = fmt.Errorf("Failed to write source to to temporary dir - %s", err)
		return
	}

	fmt.Printf("Writing to file: %s - %s\n", dir, tmpfn)

	// compile with GopherJS
	cartName := filepath.Join(dir, "cart.js")
	// we use GOOS=linux to compile to JS even on windows...
	cmd := exec.Command(gopherJS, "build", tmpfn, "-o", cartName)
	cmd.Env = append(os.Environ(),"GOOS=linux")
	//cmd := exec.Command("GOOS=linux",gopherJS, "build", tmpfn, "-o", cartName)
	var out []byte
	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("TEMP: dir:%s-%s\n",dir, tmpfn)
		fmt.Printf("TEMP: command:%s-%s\n",gopherJS, string(out))
		// decode compiler error
		a.CompErrs = getCompErrs(string(out))
		err = fmt.Errorf("Failed to compile source using GopherJS - %s", err)
		return
	}

	// copy compile code back
	var src, dst *os.File
	src, err = os.Open(cartName)
	if err != nil {
		fmt.Printf("Failed to open cart js file - %s\n", err)
		err = fmt.Errorf("Failed to open cart js file - %s", err)
		return
	}
	defer src.Close()

	destFilename := "./resources/app/dynamic/js/cart.js"

	dst, err = os.Create(destFilename)
	if err != nil {
		fmt.Printf("Failed to create target cart js file - %s\n", err)
		err = fmt.Errorf("Failed to create target cart js file - %s", err)
		return
	}
	defer dst.Close()
	fmt.Printf("TEMP: copying %s to %s\n",cartName,destFilename)
	_, err = io.Copy(dst, src)
	if err != nil {
		fmt.Printf("Failed to copy compiled cart js to target file - %s\n", err)
		err = fmt.Errorf("Failed to copy compiled cart js to target file - %s", err)
		return
	}

	return
}

func getCompErrs(output string) []CompErr {

	/*
		eg.
		../../../../../../../var/folders/5s/pxq8rc1d6wx8d5f5vsbz5vth0000gn/T/example010888711/main.go:47:2: expected operand, found 'return'
	*/
	if output == "" {
		return nil
	}

	// split into separate lines

	lines := strings.Split(output, "\n")

	errs := make([]CompErr, len(lines))

	for i, line := range lines {
		// parse line for error details

		pos := strings.Index(line, defaultSourceFile+":")
		if pos != -1 {
			// get rest of error message
			errPart := line[pos:]

			compErr := CompErr{}
			parts := strings.Split(errPart, ":")
			// parts should contain 4 parts
			if len(parts) < 4 {
				errs[i] = compErr
				fmt.Printf("TEMP: parts! %#v\n", parts)
				continue
			}

			// get line no
			compErr.Row, _ = strconv.ParseInt(parts[1], 10, 64)
			compErr.Row -= 1 // change to 0 based value
			// get column
			compErr.Column, _ = strconv.ParseInt(parts[2], 10, 64)
			// get error text - concat remaining parts with colons
			compErr.Text = strings.Join(parts[3:], ":")
			compErr.Text = strings.TrimSpace(compErr.Text)

			compErr.ErrType = "error"

			errs[i] = compErr
		}

	}

	return errs
}
