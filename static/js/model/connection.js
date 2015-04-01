var app = app || {};

/* Connection Model
 *
 * {
 *      data: <core/ConnectionLedger,
 *      model: app.CoreModel (we may be able to get rid of this at this point)
 *      from: {
 *              node: <app.Entity>      // source node
 *              route: <route>          // source node's route
 *      },
 *      to: {
 *              node: <app.Entity>      // target node
 *              route: <route>          // target node's route
 *      }
 * }
 *
 * this.to/this.from change based on what group is focussed.
 */

(function() {
    'use strict';
    app.Connection = function(data, model) {
        this.data = data;
        this.model = model;
        this.from = {
            node: model.entities[data.from.id],
            route: model.entities[data.from.id].routes.filter(function(r) {
                return (r.index === data.from.route) && (r.direction === 'output');
            })[0]
        }

        this.to = {
            node: model.entities[data.to.id],
            route: model.entities[data.to.id].routes.filter(function(r) {
                return (r.index === data.to.route) && (r.direction === 'input');
            })[0]
        }

        this.attach();
    }

    app.Connection.prototype.setNodes = function(fromNode, fromRoute, toNode, toRoute) {
        this.from.node = fromNode
        this.from.route = fromRoute
        this.to.node = toNode
        this.to.route = toRoute
    }

    // attach() and detach() adds/removes a reference to this connection the route on the block entity.
    app.Connection.prototype.attach = function() {
        this.model.entities[this.data.from.id].routes.filter(function(r) {
            return (r.index === this.data.from.route) && (r.direction === 'output')
        }.bind(this))[0].connections.push(this);

        this.model.entities[this.data.to.id].routes.filter(function(r) {
            return (r.index === this.data.to.route) && (r.direction === 'input')
        }.bind(this))[0].connections.push(this);
    }

    app.Connection.prototype.detach = function() {
        var fromConnections = this.model.entities[this.data.from.id].routes.filter(function(r) {
            return (r.index === this.data.from.route) && (r.direction === 'output')
        }.bind(this))[0].connections;

        var toConnections = this.model.entities[this.data.to.id].routes.filter(function(r) {
            return (r.index === this.data.to.route) && (r.direction === 'input')
        }.bind(this))[0].connections;

        fromConnections.splice(fromConnections.indexOf(this), 1);
        toConnections.splice(toConnections.indexOf(this), 1);
    }

    app.Connection.prototype.instance = function() {
        return "connection";
    }
})();
