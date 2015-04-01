var app = app || {};

(function() {
    'use strict';

    app.Source = function(data, model) {
        app.Entity.call(this);
        this.data = data;
        this.model = model;
    }

    app.Source.prototype = new app.Entity();

    app.Source.prototype.instance = function() {
        return "source";
    }

})();
