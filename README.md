# pico-go-electron

The aim of this project is to produce a standalone (offline) gamedev tool using Go.  It is for developing simple games / demos like the [pico8 console](https://www.lexaloffle.com/pico-8.php) but instead of coding in Lua you use Go.

The idea originally started in my [pico-go](https://github.com/telecoda/pico-go) repo but I found the external dependencies too complex/limiting for getting developers quickly productive on multiple platforms.

This tool wraps the [Ebiten](https://hajimehoshi.github.io/ebiten/) game engine which supports compiling to JS using [GopherJS](https://github.com/gopherjs/gopherjs) into an Electron app.

This code is based upon the [go-astilectron-demo](https://github.com/asticode/go-astilectron-demo) app that uses the [bootstrap](https://github.com/asticode/go-astilectron-bootstrap) and the [bundler](https://github.com/asticode/go-astilectron-bundler).

Watch this space for future developments.


# Step 1: installation

Run the following commands:

    $ go get -u github.com/telecoda/pico-go-electron/...
    $ rm $GOPATH/src/github.com/telecoda/pico-go-electron/bind.go

# Step 2: install the electron bundler

Run the following command:

    $ go get -u github.com/asticode/go-astilectron-bundler/...
    
And don't forget to add `$GOPATH/bin` to your `$PATH`.
    
# Step 3: bundle the app for your current environment

Run the following commands:

    $ cd $GOPATH/src/github.com/telecoda/pico-go-electron
    $ astilectron-bundler -v
    
# Step 4: test the app

The result is in the `output/<your os>-<your arch>` folder and is waiting for you to test it!

# Step 5: bundle the app for more environments

To bundle the app for more environments, add an `environments` key to the bundler configuration (`bundler.json`):

```json
"environments": [
  {"arch": "amd64", "os": "linux"},
  {"arch": "386", "os": "windows"}
]
```

and repeat **step 3**.
