var app = app || {};

/* Block Model
 *
 * {
 *      isDragging: <boolean>,              // provided by app.Entity
 *      parentNode: <app.Entity>,           // provided by app.Entity
 *      data: <server/types/blockLedger>    // exact server state
 *      routes: [                           // flattened version of .data.inputs/.data.outputs
 *          {                               // contains state specifically for the interface
 *
 *              direction: 'input', 'output',   // is route from .data.input or.data.output
 *              index: <int>                    // from .data.input/.data.output route index (for state)
 *              displayIndex: <int>             // from .data.input/.data.output route index (for display)
 *              id: <int>                       // from .data.id
 *              connections: [<app.Connection>] // list of of connections attached to this route
 *              data:{core/input}               // data.inputs[N]/data.outputs[N]
 *              routesIndex: <int>              // the index of this element in routes:[]
 *              parentNode: <app.Entity>        // the Group or Block this route is attachd to
 *          }, ...
 *      ],
 *      geometry:{                              // display data derived from routes{}
 *          width: <number>                     // the block width calculated from route widths
 *          height: <number>                    // the block height calculated from route height
 *          routeRadius: <number>               // the radius of the route, a constant
 *          routeHeight: <number>               // the max height of a route label
 *      }
 * }
 *
 */

(function() {
    'use strict';

    app.Block = function(data, model) {
        app.Entity.call(this);

        this.routes = [];
        this.geometry = [];
        this.data = data;
        console.log(data)
        this.model = model;

        this.buildRoutes();
        this.buildGeometry();
    }

    app.Block.prototype = new app.Entity();

    app.Block.prototype.instance = function() {
        return "block";
    }

    // transmogrify routes from .data to .routes
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
        }));

        this.routes = this.routes.map(function(r, index) {
            r.id = this.data.id;
            r.connections = [];
            r.data = this.data[r.direction + 's'][r.index];
            r.routesIndex = index;
            r.parentNode = this;
            r.source = null;
            return r
        }.bind(this));


        if (this.data.source != null) {
            this.routes.push({
                'direction': 'input',
                'index': this.data.inputs.length,
                'displayIndex': this.data.inputs.length,
                'id': this.data.id,
                'connections': [],
                'data': {
                    'name': this.data.source,
                    'value': null,
                    'type': 'any',
                },
                'routesIndex': this.routes.length,
                'parentNode': this,
                'source': this.data.source,
            })
        }
    }

    app.Block.prototype.buildGeometry = function() {
        var textMeasures,
            maxWidth = {
                input: 0,
                output: 0
            },
            routeHeight = 0,
            padding = {
              horizontal: 6,
              vertical: 6
            },
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
            'width': maxWidth.input + maxWidth.output + padding.horizontal,
            'height': Math.max(num.input, num.output) * routeHeight + padding.vertical,
            'routeHeight': routeHeight,
            'routeRadius': routeRadius,
        }
    }
})();
