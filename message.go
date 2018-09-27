package main

import (
	"encoding/json"
	"fmt"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
)

// Application represents the content of an applicaton
type Application struct {
	Path         string    `json:"path"`
	Source       string    `json:"source"`
	CompResp     *CompResp `json:"compResp"`
	ScreenWidth  int       `json:"screenWidth"`
	ScreenHeight int       `json:"screenHeight"`
}

// SourceCode used from browser to backend
type SourceCode struct {
	Path   string `json:"path"`
	Source string `json:"source"`
}

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "init":
		var path string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &path); err != nil {
				payload = err.Error()
				return
			}
		}
		payload, err = initBackend(path)
		if err != nil {
			payload = err.Error()
		}
		return
	case "load":
		var path string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &path); err != nil {
				payload = err.Error()
				return
			}
		}
		payload, err = load(path)
		if err != nil {
			payload = err.Error()
		}
		return
	case "save":
		return nil, fmt.Errorf("Save function not implemented yet")
	case "run":
		// Unmarshal payload
		var source SourceCode

		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &source); err != nil {
				payload = err.Error()
				return
			}
		}
		payload, err = run(source)
		if err != nil {
			payload = err.Error()
		}
	}
	return
}
