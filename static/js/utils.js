var app = app || {};

(function() {
    'use strict';

    app.Utils = {};

    app.Utils.request = function(method, url, data, callback) {
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
    }

    app.Utils.Debounce = function() {
        this.func = null;
        this.fire = null;
        this.last = null;
    }

    app.Utils.Debounce.prototype.push = function(e, duration) {
        if (this.last === null || this.last + duration < +new Date()) {
            this.last = +new Date();
            e();
            return;
        }
        this.func = e;
        if (this.fire != null) clearInterval(this.fire);
        this.fire = setTimeout(function() {
            this.func();
            this.last = +new Date()
        }.bind(this), duration);
    }

    app.Utils.DebounceManager = function() {
        this.entities = {};
    }

    app.Utils.DebounceManager.prototype.push = function(id, f, duration) {
        if (!this.entities.hasOwnProperty(id)) {
            this.entities[id] = new app.Utils.Debounce();
        }
        this.entities[id].push(f, duration)

    }

})();
