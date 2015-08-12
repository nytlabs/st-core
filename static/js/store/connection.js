var app = app || {};

(function() {
    'use strict';

    var connections = {};

    function Connection(data) {
        this.data = data;

        this.routeIdFrom = this.data.from.id + '_' + this.data.from.route + '_output';
        this.routeIdTo = this.data.to.id + '_' + this.data.to.route + '_input';

        //app.Dispatcher.dispatch({
        //    action: app.Actions.APP_ROUTE_UPDATE_CONNECTED,
        //    connId: data.id,
        //    ids: [this.routeIdFrom, this.routeIdTo],
        //});

        this.canvas = document.createElement('canvas');
        this.render();
    }

    Connection.prototype = Object.create(app.Emitter.prototype);
    Connection.constructor = Connection;

    Connection.prototype.geometry = function() {
        // TODO: instead of blocks, this should somehow find the top-most visible geometry that
        // the route is apart of (for groups);
        var from = app.BlockStore.getBlock(this.data.from.id);
        var to = app.BlockStore.getBlock(this.data.to.id);

        var routeIndexFrom = from.outputs.map(function(r) {
            return r.id
        }).indexOf(this.routeIdFrom);

        var routeIndexTo = to.inputs.map(function(r) {
            return r.id
        }).indexOf(this.routeIdTo);

        var yFrom = from.geometry.routeHeight * (routeIndexFrom + 1) - (from.geometry.routeRadius * .5);
        var xFrom = from.geometry.routeRadius * .5 + from.geometry.width;

        var yTo = to.geometry.routeHeight * (routeIndexTo + 1) - (to.geometry.routeRadius * .5);
        var xTo = to.geometry.routeRadius * -.5 + 0;

        this.position = {
            x: Math.min(xFrom, xTo - 50),
            y: Math.min(yFrom, yTo),
        }

        xFrom += from.position.x - this.position.x;
        yFrom += from.position.y - this.position.y;
        xTo += to.position.x - this.position.x;
        yTo += to.position.y - this.position.y;

        this.curve = [xFrom, yFrom, xFrom + 50, yFrom, xTo - 50, yTo, xTo, yTo];

        var xMax = Math.max(xFrom + 50, xTo);
        var yMax = Math.max(yFrom, yTo);

        this.canvas.width = xMax - this.position.x + 100;
        this.canvas.height = yMax - this.position.y + 100;
    }

    Connection.prototype.render = function() {
        this.geometry();

        var ctx = this.canvas.getContext('2d');
        var c = this.curve;

        ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);

        ctx.beginPath()
        ctx.moveTo(c[0], c[1]);
        ctx.setLineDash([]);
        ctx.lineWidth = 2.0
        ctx.bezierCurveTo(c[2], c[3], c[4], c[5], c[6], c[7]);
        ctx.stroke();

        this.emit();

        //app.Dispatcher.dispatch({
        //    action: app.Actions.APP_CONNECTION_UPDATE,
        //})
    }

    function ConnectionStore() {}
    ConnectionStore.prototype = Object.create(app.Emitter.prototype);
    ConnectionStore.constructor = ConnectionStore;

    ConnectionStore.prototype.getConnection = function(id) {
        return connections[id];
    }

    ConnectionStore.prototype.getConnections = function() {
        return Object.keys(connections);
    }

    var rs = new ConnectionStore();

    function createConnection(connection) {
        if (connections.hasOwnProperty(connection.id) === true) {
            console.warn('could not create connection:', connection.id, ' already exists');
            return
        }
        connections[connection.id] = new Connection(connection);
    }

    function deleteConnection(id) {
        if (connections.hasOwnProperty(id) === false) {
            console.warn('could not delete connections: ', id, ' does not exist');
            return
        }
        delete connections[id]
    }

    function renderConnections(ids) {
        ids.forEach(function(id) {
            connections[id].render()
        })
    }

    function translateConnections(ids, dx, dy) {
        ids.forEach(function(id) {
            connections[id].position.x += dx;
            connections[id].position.y += dy;
        })
    }

    function requestConnection(pickedRoutes) {
        var routes = pickedRoutes.map(function(route) {
            return app.RouteStore.getRoute(route.id);
        });

        var from = routes.filter(function(route) {
            return route.direction === 'output';
        });

        var to = routes.filter(function(route) {
            return route.direction === 'input';
        })

        if (from.length === 0 || to.length === 0) return;

        from = from[0];
        to = to[0];

        app.Utils.request(
            'POST',
            'connections', {
                'from': {
                    'id': from.blockId,
                    'route': from.index,
                },
                'to': {
                    'id': to.blockId,
                    'route': to.index,
                }
            },
            null)
    }

    app.Dispatcher.register(function(event) {
        switch (event.action) {
            case app.Actions.APP_REQUEST_CONNECTION:
                requestConnection(event.routes);
                break;
                //case app.Actions.APP_CONNECTION_UPDATE:
                //     rs.emit();
                //    break;
            case app.Actions.WS_CONNECTION_CREATE:
                createConnection(event.data);
                rs.emit();
                break;
            case app.Actions.WS_CONNECTION_DELETE:
                console.log(event.action);
                deleteConnection(event.id);
                rs.emit();
                break;
            case app.Actions.APP_RENDER_CONNECTIONS:
                renderConnections(event.ids);
                translateConnections(event.translate, event.dx, event.dy);
                rs.emit();
                break;
        }
    })

    app.ConnectionStore = rs;
}())
