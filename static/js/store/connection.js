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
        // TODO: instead of blocks, this should somehow find the top-most 
        // visible geometry that the route is apart of (for groups);
        var from = app.NodeStore.getNode(this.data.from.id);
        var to = app.NodeStore.getNode(this.data.to.id);
        // buffer accounts for bends in the bezier that may extend outside the
        // bounds of a non-buffered box.
        var buffer = 10;

        var routeIndexFrom = from.outputs.map(function(r) {
            return r.id
        }).indexOf(this.routeIdFrom);

        var routeIndexTo = to.inputs.map(function(r) {
            return r.id
        }).indexOf(this.routeIdTo);

        var yFrom = from.geometry.routeHeight * (routeIndexFrom + 1) - (from.geometry.routeRadius);
        var xFrom = from.geometry.routeRadius * 2 + from.geometry.width - 1;

        var yTo = to.geometry.routeHeight * (routeIndexTo + 1) - (to.geometry.routeRadius);
        var xTo = 1;

        // origin point for bounding box outside of connection
        this.position = {
            x: Math.min(from.position.x, to.position.x) - buffer,
            y: Math.min(from.position.y, to.position.y) - buffer,
        }

        // remove any translation from the connection as we are doing the
        // translation on the bounding box.
        xFrom += from.position.x - this.position.x
        yFrom += from.position.y - this.position.y
        xTo += to.position.x - this.position.x
        yTo += to.position.y - this.position.y

        this.curve = [xFrom, yFrom, xFrom + 50, yFrom, xTo - 50, yTo, xTo, yTo];

        // set the bounding box on the lower right, encapsulate the farthest
        // +x, +y block inside the bounding box. 
        var xMax = Math.max(from.position.x, to.position.x);
        var yMax = Math.max(from.position.y, to.position.y);

        var bWidth = from.position.x > to.position.x ? from.canvas.width : to.canvas.width;
        var bHeight = from.position.y > to.position.y ? from.canvas.height : to.canvas.height;

        this.canvas.width = xMax - this.position.x + bWidth + buffer;
        this.canvas.height = yMax - this.position.y + bHeight + buffer;
    }

    Connection.prototype.render = function() {
        this.geometry();

        var ctx = this.canvas.getContext('2d');
        var c = this.curve;

        // TODO: gradient coloring of connections is a test and should be evaluated for performance!
        var fromColor = app.Constants.TypeColors[app.RouteStore.getRoute(this.routeIdFrom).data.type];
        var toColor = app.Constants.TypeColors[app.RouteStore.getRoute(this.routeIdTo).data.type];

        var gradient = ctx.createLinearGradient(c[0], c[1], c[6], c[7]);
        gradient.addColorStop("0", fromColor);
        gradient.addColorStop("1.0", toColor);

        ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        ctx.setLineDash([]);
        ctx.lineWidth = 4.0

        ctx.strokeStyle = 'black';
        var conn = new Path2D();
        conn.moveTo(c[0], c[1]);
        conn.bezierCurveTo(c[2], c[3], c[4], c[5], c[6], c[7]);
        ctx.stroke(conn);

        ctx.strokeStyle = gradient;
        ctx.lineWidth = 2.0;
        ctx.stroke(conn);

        /*var path = new Path2D();
        path.arc(c[0], c[1], 5, 0, Math.PI * 2, true);
        path.arc(c[6], c[7], 5, 0, Math.PI * 2, true);
        ctx.fill(path);*/

        this.emit();
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

        app.Dispatcher.dispatch({
            action: app.Actions.APP_ADD_NODE_CONNECTION,
            fromId: connection.from.id,
            toId: connection.to.id,
            id: connection.id,
        })
    }

    function deleteConnection(id) {
        if (connections.hasOwnProperty(id) === false) {
            console.warn('could not delete connections: ', id, ' does not exist');
            return
        }

        app.Dispatcher.dispatch({
            action: app.Actions.APP_DELETE_NODE_CONNECTION,
            fromId: connections[id].data.from.id,
            toId: connections[id].data.to.id,
            id: id,
        });

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
            case app.Actions.WS_CONNECTION_CREATE:
                createConnection(event.data);
                rs.emit();
                break;
            case app.Actions.WS_CONNECTION_DELETE:
                deleteConnection(event.id);
                rs.emit();
                break;
            case app.Actions.APP_RENDER_CONNECTIONS:
                renderConnections(event.ids);
                rs.emit();
                break;
            case app.Actions.APP_TRANSLATE_CONNECTIONS:
                translateConnections(event.translate, event.dx, event.dy);
                renderConnections(event.ids);
                rs.emit();
                break;
        }
    })

    app.ConnectionStore = rs;
}())
