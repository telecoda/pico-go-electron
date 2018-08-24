package main

import (
	"encoding/json"
	"fmt"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
)

// CompErr - compiler errors
type CompErr struct {
	Row     int64  `json:"row"`
	Column  int64  `json:"col"`
	Text    string `json:"text"`
	ErrType string `json:"type"`
}

// Application represents the content of an applicaton
type Application struct {
	Path     string `json:"path"`
	Source   string `json:"source"`
	CompErrs []CompErr
}

const (
	defaultCodeDir    = "gosrc"
	defaultSourceFile = "main.go"
)

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "load":
		payload, err = load("")
		if err != nil {
			payload = err.Error()
		}
		return
	case "save":
		return nil, fmt.Errorf("Save function not implemented yet")
	case "run":
		// Unmarshal payload
		var source string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &source); err != nil {
				payload = err.Error()
				return
			}
		}
		return run(source)
		//payload, err = run(source)
		// if err != nil {
		// 	payload = err.Error()
		// }
		// return
	}
	return
}
