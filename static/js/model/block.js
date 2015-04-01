var app = app || {};

(function() {
    'use strict';
    console.log(app)
    app.Block = function(data, model) {
        app.Entity.call(this);

        this.routes = [];
        this.geometry = [];
        this.data = data;
        this.model = model;

        this.buildRoutes();
        this.buildGeometry();
    }

    app.Block.prototype = new app.Entity();

    app.Block.prototype.instance = function() {
        return "block";
    }

    app.Block.prototype.buildRoutes = function() {
        this.routes = this.data.inputs.map(function(input, index) {
            return {
                'direction': 'input',
                'index': index,
                'displayIndex': index
            };
        })

        this.routes = this.routes.concat(this.data.outputs.map(function(output, index) {
            return {
                'direction': 'output',
                'index': index,
                'displayIndex': index
            }
        }))

        this.routes = this.routes.map(function(r, index) {
            r.id = this.data.id;
            r.connections = [];
            r.data = this.data[r.direction + 's'][r.index];
            r.routesIndex = index;
            r.parentNode = this;
            return r
        }.bind(this));

    }

    app.Block.prototype.buildGeometry = function() {
        var textMeasures,
            maxWidth = {
                input: 0,
                output: 0
            },
            routeHeight = 0,
            routeRadius = 5,
            routeGeometry = [],
            num = {
                input: 0,
                output: 0
            };

        textMeasures = this.routes.map(function(r) {
            var measure = app.Utils.measureText(r.data.name, 'route_label');

            if (measure.width > maxWidth[r.direction]) {
                maxWidth[r.direction] = measure.width;
            }
            if (measure.height > routeHeight) {
                routeHeight = measure.height;
            }
            num[r.direction]++;

            return measure
        });

        this.geometry = {
            'width': maxWidth.input + maxWidth.output,
            'height': Math.max(num.input, num.output) * routeHeight,
            'routeHeight': routeHeight,
            'routeRadius': routeRadius,
        }
    }
})();
