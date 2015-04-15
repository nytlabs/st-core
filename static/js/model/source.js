var app = app || {};

(function() {
    'use strict';

    app.Source = function(data, model) {
        app.Entity.call(this);
        this.data = data;
        this.model = model;

        this.routes = [];
        this.geometry = [];

        this.buildRoutes();
        this.buildGeometry();
    }

    app.Source.prototype = Object.create(app.Block.prototype, {});

    app.Source.prototype.instance = function() {
        return 'source';
    }

    app.Source.prototype.buildRoutes = function() {
        this.routes = [{
            'direction': 'output',
            'index': 0,
            'displayIndex': 0,
            'id': this.data.id,
            'connections': [],
            'data': {
                'name': this.data.type,
                'value': null,
                'type': 'any',
            },
            'routesIndex': 0,
            'parentNode': this,
            'source': this.data.type,
        }];
    }
})();
