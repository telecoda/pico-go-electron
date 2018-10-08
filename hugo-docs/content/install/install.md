+++
title = "install"
draft = false
date = "2018-10-06T16:26:07+01:00"

+++
# pico-go-electron Installation

These are the instructions to get the project up and running.  Please contact me or raise a PR if you have issues building and running pico-go on your chosen platform. (Currently it has been tested on OSX and Windows).

## Step 1: installation

Run the following commands:

    $ go get -u github.com/telecoda/pico-go-electron/...
    $ rm $GOPATH/src/github.com/telecoda/pico-go-electron/bind.go

on Windows:

    $ rm %GOPATH%/src/github.com/telecoda/pico-go-electron/bind.go

## Step 2: install the prerequisites

Run the following command:

    $ go get -u github.com/asticode/go-astilectron-bundler/...
    $ go get -u github.com/gopherjs/gopherjs
    $ go get -u github.com/gopherjs/gopherwasm/js
    $ go get -u github.com/hajimehoshi/ebiten
    
And don't forget to add `$GOPATH/bin` to your `$PATH`.
    
## Step 3: bundle the app for your current environment

Run the following commands:

    $ cd $GOPATH/src/github.com/telecoda/pico-go-electron
    $ astilectron-bundler -v
    
on Windows:

    $ cd %GOPATH%/src/github.com/telecoda/pico-go-electron
    $ astilectron-bundler.exe -v

## Step 4: test the app

The result is in the `output/<your os>-<your arch>` folder and is waiting for you to test it!

## Step 5: bundle the app for more environments

To bundle the app for more environments, add an `environments` key to the bundler configuration (`bundler.json`):

```json
"environments": [
  {"arch": "amd64", "os": "linux"},
  {"arch": "386", "os": "windows"}
]
```

and repeat **step 3**.

# Developing pico-go itself

The instructions above are for creating a fully fledged native Electron app for your operating system. However if you wish to tinker with `pico-go-electron` yourself and just run a simple go binary there is a far easier way to do things.

Providing you already have all the necesary prerequistes installed when you are in the project root directory type:

    go build

To run the app

    ./pico-go-electron -d

or on Windows

    pico-go-electron.exe -d

The `-d` flag enables debug mode which allows you to start the chrome developer tools if you are planning to tinker with the javascript.  There will be a debug option on the applications main menu.