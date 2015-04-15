var app = app || {};

(function() {
    'use strict';
    app.Link = function(data, model) {
        //app.Entity.call(this);
        this.data = data;
        this.model = model;

        this.from = {
            node: model.entities[data.source.id],
            route: model.entities[data.source.id].routes.filter(function(r) {
                return r.source != 'null'
            })[0]
        }

        this.to = {
            node: model.entities[data.block.id],
            route: model.entities[data.block.id].routes.filter(function(r) {
                return r.source != 'null'
            })[0]
        }

    }

    //app.Link.prototype = new app.Entity();

    app.Link.prototype.instance = function() {
        return 'link';
    }

    app.Link.prototype.setNodes = function(fromNode, fromRoute, toNode, toRoute) {
        this.from.node = fromNode
        this.from.route = fromRoute
        this.to.node = toNode
        this.to.route = toRoute
    }

})();
