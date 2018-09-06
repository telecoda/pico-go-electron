package main

import (
	"encoding/json"
	"fmt"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
)



// Application represents the content of an applicaton
type Application struct {
	Path     string `json:"path"`
	Source   string `json:"source"`
	CompResp *CompResp `json:"compResp"`
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
	}
	return
}
