let index = {
    about: function(message) {
        dialog.showMessageBox({"title": "About","message": message});
    },
    init: function() {
        // Init
        asticode.loader.init();
        
        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {
            // Listen
            index.listen();

            // Load initial code
            index.load();
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

            console.log(message.payload)
            document.getElementById("path").innerHTML = "path: " +message.payload.path;
            editor.session.setValue(message.payload.source)
        })
    },
    run: function() {
        // Create message
        let message = {"name": "run",
            "payload": editor.session.getValue()
        };

        // send sourcecode to backend for compilation

        // Send message
        asticode.loader.show();
        astilectron.sendMessage(message, function(message) {
            // Init
            asticode.loader.hide();
            editor.session.clearAnnotations();
            // Check error
            if (message.name === "error") {
                // convert response to annotations on sourcecode
                annotations = [];
                errorMessage = "";
                if (message.payload.CompErrs != undefined && message.payload.CompErrs.length > 0) {
                    errs = message.payload.CompErrs
                    for (var i = 0; i < errs.length; i++) {
                        annotations.push(errs[i]);
                        errorMessage += errs[i].text + "\n"
                    }
                    editor.session.setAnnotations(annotations);
                    dialog.showErrorBox("Compile Error",errorMessage);
                    return
                }
                
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