package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	widthField  = "screenWidth"
	heightField = "screenHeight"

	defaultWidth  = 320
	defaultHeight = 240
)

// CompErr - compiler errors
type CompErr struct {
	Row     int64  `json:"row"`
	Column  int64  `json:"col"`
	Text    string `json:"text"`
	ErrType string `json:"type"`
}

// CompResp - compiler response
type CompResp struct {
	Raw    string    `json:"raw"`
	Errors []CompErr `json:"errors"`
}

// run - compiles and executes latest code
func run(sourceCode SourceCode) (a Application, err error) {

	// save source to file
	dir, err := ioutil.TempDir("", "example")
	if err != nil {
		err = fmt.Errorf("Failed to create temporary dir - %s", err)
		return
	}

	defer os.RemoveAll(dir) // clean up

	tmpCartFile := filepath.Join(dir, "cart.go")
	if err = ioutil.WriteFile(tmpCartFile, []byte(sourceCode.Source), 0666); err != nil {
		err = fmt.Errorf("Failed to write cart source to to temporary dir - %s", err)
		return
	}

	tmpMainFile := filepath.Join(dir, "gen-main.go")
	if err = ioutil.WriteFile(tmpMainFile, []byte(genMainSrc), 0666); err != nil {
		err = fmt.Errorf("Failed to write gen-main source to to temporary dir - %s", err)
		return
	}


	// compile with GopherJS
	outFile := filepath.Join(dir, "cart.js")

	cmd := getBuildCmd(tmpCartFile, tmpMainFile, outFile)
	var out []byte
	out, err = cmd.CombinedOutput()
	if err != nil {
		if out == nil {
			err = fmt.Errorf("Failed to call gopherjs - %s", err)
			return
		}
		raw := string(out)
		compResp := &CompResp{
			Raw: raw,
		}
		// decode compiler error
		compResp.Errors = getCompErrs(raw)
		a.CompResp = compResp
		err = nil
		return
	}

	// copy compile code back
	var src, dst *os.File
	src, err = os.Open(outFile)
	if err != nil {
		fmt.Printf("Failed to open cart js file - %s\n", err)
		err = fmt.Errorf("Failed to open cart js file - %s", err)
		return
	}
	defer src.Close()

	destFilename := filepath.Join(sourceCode.Path, "Local Storage", "cart.js")

	dst, err = os.Create(destFilename)
	if err != nil {
		fmt.Printf("Failed to create target cart js file - %s\n", err)
		err = fmt.Errorf("Failed to create target cart js file - %s", err)
		return
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		fmt.Printf("Failed to copy compiled cart js to target file - %s\n", err)
		err = fmt.Errorf("Failed to copy compiled cart js to target file - %s", err)
		return
	}

	a.ScreenWidth, a.ScreenHeight = getScreenDimensions(sourceCode.Source)

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

// getScreenDimensions - inspects source code for declaration of screen dimensions to pass to javascript
func getScreenDimensions(source string) (int, int) {

	var width int
	var height int

	buf := bytes.NewBuffer([]byte(source))
	r := bufio.NewReader(buf)

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		if strings.Contains(line, widthField) && width == 0 {
			width = getIntValue(line, widthField)
		}
		if strings.Contains(line, heightField) && height == 0 {
			height = getIntValue(line, heightField)
		}

		if width != 0 && height != 0 {
			break
		}

	}

	// if sizes are not found, use sensible defaults
	if width == 0 {
		width = defaultWidth
	}
	if height == 0 {
		height = defaultHeight
	}

	return width, height
}

// getIntValue - gets a integer value of a field base on a line in the source code
func getIntValue(line, fieldname string) int {

	// This code looks for declarations or assignments
	// eg. screenWidth := 320
	// or. screenWidth = 320

	trimmed := strings.TrimSpace(line)
	trimmed = strings.Replace(trimmed, "	", "", -1) // remove tabs
	replaced := strings.Replace(trimmed, "var "+fieldname, "", 1)
	replaced = strings.Replace(replaced, "const "+fieldname, "", 1)
	replaced = strings.Replace(replaced, fieldname, "", 1)
	replaced = strings.Replace(replaced, ":=", "", 1)
	replaced = strings.Replace(replaced, "=", "", 1)
	noComments := strings.Split(replaced, "//")
	intStr := strings.Split(noComments[0], "/*")

	// We "should" just have an integer value as a string now
	finalStr := strings.TrimSpace(intStr[0])

	i, err := strconv.ParseInt(finalStr, 10, 64)
	if err != nil {
		return 0
	}

	return int(i)
}
