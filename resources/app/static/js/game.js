let gameFuncs = {
    init: function() {
        if (typeof parent.userPath != 'undefined') {
            var script_element = document.body.getElementsByTagName("script");
            var cart_src = parent.userPath + "/Local Storage/cart.js";
            if (script_element[0].src =="") {
                script_element[0].innerText="";
                script_element[0].src = cart_src;
            }
        } 
    }
};



