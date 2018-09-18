const {app} = require('electron').remote
let index = {
    about: function(message) {
        dialog.showMessageBox({"title": "About","message": message});
    },
    init: function() {
        // Init
        asticode.loader.init();
        // capture electron app "userData" variable 
        // this is used as a location to host the compiled JS file.
        userPath = app.getPath("userData");

        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {
            // Listen
            index.listen();
        
            // call backend init
            // Create message
            let message = {"name": "init"};
            if (typeof userPath !== "undefined") {
                message.payload = userPath
            }
            // Send message
            asticode.loader.show();
            astilectron.sendMessage(message, function(message) {
                // Init
                asticode.loader.hide();

                // Check error
                if (message.name === "error") {
                    dialog.showErrorBox("Init Error",message.payload);
                    return
                }
            })

            
            // load initial source code
            index.load(userPath);
        })
    },
    load: function(path) {
        // Create message
        let message = {"name": "load"};
        if (typeof path !== "undefined") {
            message.payload = path
        }
        // Send message
        asticode.loader.show();
        astilectron.sendMessage(message, function(message) {
            // Init
            asticode.loader.hide();

            // Check error
            if (message.name === "error") {
                dialog.showErrorBox("Load Error",message.payload);
                return
            }

            document.getElementById("path").innerHTML = "path: " +message.payload.path;
            editor.session.setValue(message.payload.source)
        })
    },
    run: function() {

        payload = {
            "path": userPath,
            "source":editor.session.getValue()
        }
        // Create message
        let message = {"name": "run",
            "payload": payload
        };

        // send sourcecode to backend for compilation
        asticode.loader.show();
        astilectron.sendMessage(message, function(message) {
            // Init
            asticode.loader.hide();
            editor.session.clearAnnotations();
            document.getElementById("compErrors").innerHTML = ""; // Clear errors
            // Check error
            if (message.name === "error") {
                dialog.showErrorBox("Load Error",message.payload);
                return
            }

            // convert response to annotations on sourcecode
            annotations = [];
            errorMessage = "";
            if (message.payload.compResp != undefined && message.payload.compResp.errors != undefined && message.payload.compResp.errors.length > 0) {
                errs = message.payload.compResp.errors
                for (var i = 0; i < errs.length; i++) {
                    annotations.push(errs[i]);
                    errorMessage += errs[i].text + "\n"
                }
                editor.session.setAnnotations(annotations);
                document.getElementById("compErrors").innerHTML =message.payload.compResp.raw;
                dialog.showErrorBox("Compile Error",errorMessage);
                return
            }
            
            // if no errors - code compiled successfully
            // switch to game tab
            document.getElementById("gameFrame").contentWindow.location.reload();
            document.getElementById("gameTab").click();
            // refresh js
        })
    },
    listen: function() {
        astilectron.onMessage(function(message) {
            switch (message.name) {
                case "about":
                    index.about(message.payload);
                    return {payload: "payload"};
                    break;
            }
        });
    }
};