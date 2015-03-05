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

    /* apps.Utils.measureText provides bounding box information for text nodes
     * created with SVG. AFAIK, it's impossible to measure bounding boxes
     * without actually rendering an element. Even setting display: none sets
     * the measurements to 0 (which makes sense, I guess)
     *
     * This function re-uses an already instantiated (and shoddily hidden)
     * text node so that we can easily measure text without needing to worry
     * about rendering, checking, and re-rendering elements elsewhere.
     */
    var cachedText = document.createElement('text');
    cachedText.classList.add('measure_text_scratch');
    document.body.appendChild(cachedText);

    app.Utils.measureText = function(text, className) {
        var textNode = document.createTextNode(text);
        if (!!className) cachedText.classList.add(className);
        cachedText.appendChild(textNode);
        var length = cachedText.getBoundingClientRect()
        cachedText.removeChild(textNode);
        if (!!className) cachedText.classList.remove(className);
        return length
    }

    /* apps.Utils.pointInRect checks to see if (px, py) is in rect defined by
     * x,y,w,h.
     */
    app.Utils.pointInRect = function(x, y, w, h, px, py) {
        return (px >= x) && (px < x + w) && (py >= y) && (py < y + h);
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
