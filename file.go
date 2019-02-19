package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
)

var watcher *fsnotify.Watcher
var cancel chan bool

// load - loads sourcecode from a specific path
func load(path string) (a Application, err error) {
	
	if strings.HasSuffix(path, ".go") {
		return loadGo(path)
	}
	
	if strings.HasSuffix(path, ".gif") {
		return loadSprites(path)
	}

	err = fmt.Errorf("Failed to open file: %s. File MUST be a .go or .gif file", path)

	return


}


// load - loads sprite data from a specific path
func loadSprites(path string) (a Application, err error) {

	f, err := os.Open(path)
	if err != nil {
		err = fmt.Errorf("Failed to open file: %s", err)
		return
	}

	spriteData, err := ioutil.ReadAll(f)
	if err != nil {
		err = fmt.Errorf("Failed to read file: %s", err)
		return
	}

	

	// Init Application
	a = Application{
		Path:   path,
		SpriteData: base64.StdEncoding.EncodeToString(spriteData),
	}

	return
}

// loadGo - loads sourcecode from a specific path
func loadGo(path string) (a Application, err error) {

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

	if watcher != nil {
		// close old watcher
		cancel <- true
		err = watcher.Close()
		if err != nil {
			err = fmt.Errorf("Failed to close previous watcher %s - %s", path, err)
			return
		}
	}

	// creates a new file watcher
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return
	}
	cancel = make(chan bool)
	if err = watcher.Add(path); err != nil {
		return
	}

	go func(reloadPath string, cancel chan bool) {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				//fmt.Printf("EVENT! %#v\n", event)
				if event.Op == fsnotify.Write {
					// reload
					if err := bootstrap.SendMessage(w, "reload", reloadPath, func(m *bootstrap.MessageIn) {
						// Unmarshal payload
						var s string
						if m != nil {
							if err := json.Unmarshal(m.Payload, &s); err != nil {
								astilog.Error(errors.Wrap(err, "unmarshaling payload failed"))
								return
							}
						}
					}); err != nil {
						astilog.Error(errors.Wrap(err, "sending reload event failed"))
					}
				}

				// watch for errors
			case err := <-watcher.Errors:
				if err != nil {
					fmt.Printf("ERROR: %s watching for file changes", err)
				}

			case _ = <-cancel:
				return
			}
		}
	}(path, cancel)

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
