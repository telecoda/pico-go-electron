const {app} = require('electron').remote
let index = {

    // init - called during application start, invokes backend initialisation
    init: function() {
        console.log("Set runCart false\n")
        localStorage.setItem("runCart", false);
        // Init
        asticode.loader.init();
        // capture electron app "userData" variable 
        // this is used as a location to host the compiled JS file.
        userPath = app.getPath("userData");

        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {
            // Listen
            index.listen();
            app.focus();
        
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
                    app.quit()
                    return
                }
            })
        })
    },
    // about menu clicked
    aboutMenu: function(message) {
        dialog.showMessageBox({"title": "About","message": message});
    },
    // new menu clicked
    newMenu: function(payload) {
        document.title = "<untitled>";
        editor.session.setValue(payload)
    },
    // open menu clicked
    openMenu: function(message) {
        filenames = dialog.showOpenDialog({"title": "Select sourcefile","message": message});
        if (typeof filenames !== "undefined" && filenames.length >0)  {
            filename = filenames[0];
            if (filename.endsWith(".go")) {
                this.load(filename);
            } else {
                this.loadSprites(filename);
            }
        }
    },
        // reload source changed
    reload: function(path) {
        // check if source in editor has been changed since it was loaded
        console.log("Reloading: " + path);
        that = this;
        cbFunc = function(buttonIndex){ 
            if (buttonIndex == 0) {
                that.load(path);
            }
        }
        if (loadedSource != editor.session.getValue()) {
            // display alert to ask for permission
            dialog.showMessageBox({"title": "Source changed","message": "The sourcecode has been changed in this editor and in an external editor, do you wish to load the latest changes and loose anything you have done here?", "type": "question", "buttons": [ "OK","cancel"], "defaultId": 0 }, cbFunc );
        } else {
            // just load it
            this.load(path);
        }
        
    },
    // run menu clicked
    runMenu: function(message) {
        this.run();
    },
    // save menu clicked
    saveMenu: function(message) {
        if (typeof filename !== "undefined") {
            this.save(filename);
        } else {
            this.saveAsMenu();
        }
    },
    // saveAs menu clicked
    saveAsMenu: function() {
        filename = dialog.showSaveDialog({"title": "Select file to save as","filters": [{"name":"go files","extensions":["go"]}]});
        if (typeof filename !== "undefined") {
            this.save(filename);
        }
    },
    // load - call backend to load source from path and init editor
    load: function(filename) {
        // Create message
        let message = {"name": "load"};
        if (typeof filename !== "undefined") {
            payload = {
                "path": filename,
            }
            message.payload = payload
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

            document.title = message.payload.path;
            editor.session.setMode("ace/mode/golang");
            editor.session.setValue(message.payload.source)
            // save loaded source
            loadedSource = message.payload.source;

            // switch to code tab
            document.getElementById("codeTab").click();

        })
    },

    loadSprites: function(filename) {
            // Create message
            let message = {"name": "loadSprites"};
            if (typeof filename !== "undefined") {
                payload = {
                    "path": filename,
                }
                message.payload = payload
            }
            // Send message
            asticode.loader.show();
            astilectron.sendMessage(message, function(message) {
                // Init
                asticode.loader.hide();
    
                // Check error
                if (message.name === "error") {
                    dialog.showErrorBox("Load Sprites Error",message.payload);
                    return
                }    
                // save loaded sprite data in local storage for sprite editor
                global.localStorage.setItem("pico-go-sprite-data",message.payload.spriteData)
                // set reload flag - so editor fetches latest sprite data
                global.localStorage.setItem("pico-go-sprite-data-reload",true)

                // switch to sprites tab
                document.getElementById("spriteEdTab").click();
            })
        },
        // run - call to backend to compile and run current sourcecode
    run: function() {

        payload = {
            "path": userPath,
            "source":editor.session.getValue()
        }
        // Create message
        let message = {"name": "run",
            "payload": payload
        };

        let spriteData = localStorage.getItem("sprite-data");

        console.log("SpriteData: " + spriteData);

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
                // highlight first error line
                editor.moveCursorToPosition(errs[0].row-1,0);
                editor.moveCursorTo(errs[0].row,0);
                editor.scrollToLine(errs[0].row-1);
                editor.session.setAnnotations(annotations);
                document.getElementById("compErrors").innerHTML =message.payload.compResp.raw;
                dialog.showErrorBox("Compile Error",errorMessage);
                return
            }
            
            // if no errors - code compiled successfully
            // switch to game tab
            screenWidth = message.payload.screenWidth;
            screenHeight = message.payload.screenHeight;
            localStorage.setItem("runCart", true);
            document.getElementById("gameFrame").contentWindow.location.reload();
            document.getElementById("gameTab").click();
            // refresh js
        })
    },

    // save - saves sourcecode to filename
    save: function(filename) {
        // Create message
        let message = {"name": "save"};
        if (typeof filename !== "undefined") {
            payload = {
                "path": filename,
                "source":editor.session.getValue()
            }
            message.payload = payload
        }
        // Send message
        asticode.loader.show();
        astilectron.sendMessage(message, function(message) {
            // Init
            asticode.loader.hide();

            // Check error
            if (message.name === "error") {
                dialog.showErrorBox("Save Error",message.payload);
                return
            }
            // update loadedSource to track for code changes
            loadedSource = editor.session.getValue()
            document.title = message.payload.path;
        })
    },

    // meno option listener
    listen: function() {
        astilectron.onMessage(function(message) {
            switch (message.name) {
                case "about":
                    index.aboutMenu(message.payload);
                    return {payload: "about clicked!"};
                    break;
                case "new":
                    index.newMenu(message.payload);
                    return {payload: "new clicked!"};
                    break;
                case "open":
                    index.openMenu(message.payload);
                    return {payload: "open clicked!"};
                    break;
                case "reload":
                    index.reload(message.payload);
                    return {payload: "reload source changed"};
                    break;
                case "run":
                    index.runMenu(message.payload);
                    return {payload: "run clicked!"};
                    break;
                case "save":
                    index.saveMenu(message.payload);
                    return {payload: "save clicked!"};
                    break;
                case "saveAs":
                    index.saveAsMenu(message.payload);
                    return {payload: "saveAs clicked!"};
                    break;
            }
        });
    }
};