var app = app || {};

(function() {
    'use strict';
    app.Link = function(data, model) {
        //app.Entity.call(this);
        this.data = data;
        this.model = model;
    }

    //app.Link.prototype = new app.Entity();

    app.Link.prototype.instance = function() {
        return "link";
    }

    app.Link.prototype.refreshGeometry = function() {
        //TODO
    }
})();
