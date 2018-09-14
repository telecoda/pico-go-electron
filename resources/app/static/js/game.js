let gameFuncs = {
    init: function() {
        if (typeof parent.userPath != 'undefined') {
            var cart_script = document.createElement('script');
            cart_script.setAttribute('src',parent.userPath + "/Local Storage/cart.js");
            document.body.appendChild(cart_script);    
            console.log("path here:" +parent.userPath);
        } 
        console.log("game.init()");       
    }
};



