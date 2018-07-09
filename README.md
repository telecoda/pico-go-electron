# pico-go-electron

This code is based upon the [go-astilectron-demo](https://github.com/asticode/go-astilectron-demo) app that uses the [bootstrap](https://github.com/asticode/go-astilectron-bootstrap) and the [bundler](https://github.com/asticode/go-astilectron-bundler).

It's also the subject of this [blog post](https://medium.com/@social_57971/how-to-add-a-gui-to-your-golang-app-in-5-easy-steps-c25c99d4d8e0).

![screenshot](screenshot.png)

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