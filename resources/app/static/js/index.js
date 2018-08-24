let index = {
    about: function(html) {
        let c = document.createElement("div");
        c.innerHTML = html;
        asticode.modaler.setContent(c);
        asticode.modaler.show();
    },
    addFolder(name, path) {
        let div = document.createElement("div");
        div.className = "dir";
        div.onclick = function() { index.explore(path) };
        div.innerHTML = `<i class="fa fa-folder"></i><span>` + name + `</span>`;
        document.getElementById("dirs").appendChild(div)
    },
    init: function() {
        // Init
        asticode.loader.init();
        asticode.modaler.init();
        asticode.notifier.init();

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
                asticode.notifier.error(message.payload);
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
                asticode.notifier.error("Compilation error(s)");

                // convert response to annotations on sourcecode
                annotations = [];
                if (message.payload.CompErrs != undefined && message.payload.CompErrs.length > 0) {
                    errs = message.payload.CompErrs
                    for (var i = 0; i < errs.length; i++) {
                        annotations.push(errs[i]);
                    }
                    editor.session.setAnnotations(annotations);
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
                case "check.out.menu":
                    asticode.notifier.info(message.payload);
                    break;
            }
        });
    }
};