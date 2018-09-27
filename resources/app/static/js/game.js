let gameFuncs = {
    init: function() {
        if (typeof parent.userPath != 'undefined') {
            /*
                What's this code doing?
                =======================
                When we compile the .go code through GopherJS we create a .js file and copy it
                to the "userData" +"/Local Storage/cart.js" location.

                This is so the code generation remains cross platform.

                Within the game.html page is an empty <script> tag.
                The code below updates the tag with the location of the .js file which
                forces the js to be loaded and run.
            */
            var script_element = document.body.getElementsByTagName("script");
            var cart_src = parent.userPath + "/Local Storage/cart.js";
            if (script_element[0].src =="") {
                script_element[0].innerText="";
                script_element[0].src = cart_src;
            }
        }
    }
};



