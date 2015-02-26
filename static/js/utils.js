var app = app || {};

(function() {
    'use strict';

    app.Utils = {
        request: function(method, url, data, callback) {
            var req = new XMLHttpRequest();
            req.open(method, url, true);
            if (data !== null) {
                req.send(JSON.stringify(data));
            } else {
                req.send(null);
            }
            req.onreadystatechange = function() {
                if (typeof callback === 'function') {
                    callback(req);
                }
            }
        },
    }
})();
