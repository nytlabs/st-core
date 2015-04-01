var app = app || {};

/* Group Model
 *
 * Extended from Block
 * {
 *      isDragging: <boolean>,              //
 *      parentNode: <app.Entity>,           //
 *      data: <server/types/blockLedger>    //
 *      routes: [                           //
 *          {                               //
 *
 *              direction: 'input', 'output',   //
 *              index: <int>                    //
 *              displayIndex: <int>             // Nth input or output for display
 *              id: <int>                       //
 *              connections: [<app.Connection>] // inherited from parent route
 *              data:{core/input}               //
 *              routesIndex: <int>              //
 *              parentNode: <app.Entity>        //
 *          }, ...
 *      ],
 *      geometry:{                              //
 *          width: <number>                     //
 *          height: <number>                    //
 *          routeRadius: <number>               //
 *          routeHeight: <number>               //
 *      }
 * }
 *
 */
(function() {
    'use strict';

    app.Group = function(data, model) {
        app.Entity.call(this);
        //        app.Block.call(this);
        this.data = data;
        this.model = model;

        this.routes = [];
        this.geometry = [];

        // translation coords for each group workspace.
        // not synced with server.
        this.translateX = 0;
        this.translateY = 0;

        this.refresh();
    }

    app.Group.prototype = Object.create(app.Block.prototype, {});

    app.Group.prototype.refresh = function() {
        this.buildRoutes();
        this.buildGeometry();
    }

    app.Group.prototype.instance = function() {
        return "group";
    }

    app.Group.prototype.buildRoutes = function() {
        var displayIndex = {
            'input': 0,
            'output': 0
        }

        this.routes = [];

        this.data.children.forEach(function(child) {
            this.model.entities[child].routes.forEach(function(r, index) {
                this.routes.push({
                    id: r.id,
                    connections: r.connections,
                    data: r.data,
                    routesIndex: index,
                    direction: r.direction,
                    index: r.index,
                    displayIndex: displayIndex[r.direction]++,
                    parentNode: this,
                })
            }.bind(this))
        }.bind(this))

        // propagate changes to parent groups
        if (this.parentNode !== null) this.parentNode.refresh();
    }

    app.Group.prototype.setTranslation = function(x, y) {
        this.translateX = x;
        this.translateY = y;
        this.model.inform();
    }


    app.Group.prototype.setFocusedGroup = function() {
        this.model.setFocusedGroup(this);
    }
})();
