var app = app || {};

(function() {
    'use strict';

    app.Utils = {};

    app.Utils.distance = function(x1, y1, x2, y2) {
        return Math.sqrt((x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2));
    }

    app.Utils.request = function(method, url, data, callback) {
        var req = new XMLHttpRequest();
        req.open(method, url, true);
        req.send(JSON.stringify(data));
        req.onreadystatechange = function() {
            if (typeof callback === 'function' && this.readyState == 4) {
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

    app.Utils.escape = function(value) {
        return value.replace(/([.*+?^${}()|\[\]\/\\])/g, "\\$1");
    };

    /* apps.Utils.pointInRect checks to see if (px, py) is in rect defined by
     * x,y,w,h.
     */
    app.Utils.pointInRect = function(x, y, w, h, px, py) {
        return (px >= x) && (px < x + w) && (py >= y) && (py < y + h);
    }

})();
