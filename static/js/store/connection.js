var app = app || {};

// TODO: move emit() out of Connection.render();

(function() {
    'use strict';

    var connections = {};
    var selected = {};

    function Connection(data) {
        this.data = data;
        this.dirtyPicking = false;
        this.routeIdFrom = this.data.from.id + '_' + this.data.from.route + '_output';
        this.routeIdTo = this.data.to.id + '_' + this.data.to.route + '_input';

        //app.Dispatcher.dispatch({
        //    action: app.Actions.APP_ROUTE_UPDATE_CONNECTED,
        //    connId: data.id,
        //    ids: [this.routeIdFrom, this.routeIdTo],
        //});

        this.pickColor = app.PickingStore.getColor(this);
        this.canvas = document.createElement('canvas');
        this.pickCanvas = document.createElement('canvas');
        this.geometry();
        this.render();
    }

    Connection.prototype = Object.create(app.Emitter.prototype);
    Connection.prototype.constructor = Connection;

    Connection.prototype.geometry = function() {
        this.dirtyPicking = true;
        // TODO: instead of blocks, this should somehow find the top-most 
        // visible geometry that the route is apart of (for groups);
        var from = app.NodeStore.getVisibleParent(this.data.from.id);
        var to = app.NodeStore.getVisibleParent(this.data.to.id);
        // buffer accounts for bends in the bezier that may extend outside the
        // bounds of a non-buffered box.
        var buffer = 10;

        var routeIndexFrom = from.outputsGeometry.map(function(r) {
            return r.id
        }).indexOf(this.routeIdFrom);

        var routeIndexTo = to.inputsGeometry.map(function(r) {
            return r.id
        }).indexOf(this.routeIdTo);

        var yFrom = from.nodeGeometry.routeHeight * (routeIndexFrom + 1) - (from.nodeGeometry.routeRadius);
        var xFrom = from.nodeGeometry.routeRadius * 2 + from.nodeGeometry.width - 1;

        var yTo = to.nodeGeometry.routeHeight * (routeIndexTo + 1) - (to.nodeGeometry.routeRadius);
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

        var bWidth = from.position.x + from.canvas.width >
            to.position.x + to.canvas.width ?
            from.canvas.width : to.canvas.width;
        var bHeight = from.position.y + from.canvas.height >
            to.position.y + to.canvas.height ?
            from.canvas.height : to.canvas.height;

        this.canvas.width = Math.floor(.5 + xMax - this.position.x + bWidth + buffer);
        this.canvas.height = Math.floor(.5 + yMax - this.position.y + bHeight + buffer);
        this.pickCanvas.width = this.canvas.width;
        this.pickCanvas.height = this.canvas.height;
    }

    Connection.prototype.render = function() {
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

        // debug bounding boxes
        //ctx.fillStyle = "rgba(255,0,0,.1)";
        //ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
        ctx.strokeStyle = app.SelectionStore.isSelected(this) ? 'blue' : 'black';
        var conn = new Path2D();
        conn.moveTo(c[0], c[1]);
        conn.bezierCurveTo(c[2], c[3], c[4], c[5], c[6], c[7]);
        ctx.lineWidth = app.SelectionStore.isSelected(this) ? 6.0 : 4.0;
        ctx.stroke(conn);

        ctx.strokeStyle = gradient;
        ctx.lineWidth = 2.0;
        ctx.stroke(conn);
        this.emit();

    }

    Connection.prototype.renderPicking = function() {
        this.dirtyPicking = false;
        var ctx = this.canvas.getContext('2d');
        var pctx = this.pickCanvas.getContext('2d');
        var imgData = ctx.getImageData(0, 0, this.pickCanvas.width, this.pickCanvas.height);
        var pixels = imgData.data;
        var rgb = this.pickColor.replace('rgb(', '').replace(')', '').split(',').map(function(str) {
            return parseInt(str)
        })

        for (var i = 0; i < pixels.length; i += 4) {
            if (pixels[i + 3] != 0) {
                pixels[i] = rgb[0];
                pixels[i + 1] = rgb[1];
                pixels[i + 2] = rgb[2];
                pixels[i + 3] = 255;
            }
        }
        pctx.putImageData(imgData, 0, 0);
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

    ConnectionStore.prototype.getSelected = function() {
        return selected;
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
            connections[id].geometry();
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
            case app.Actions.APP_RENDER:
                if (!connections.hasOwnProperty(event.id)) return;
                connections[event.id].render();
                rs.emit();
                break;
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
            case app.Actions.APP_RENDER_CONNECTION_PICKING:
                connections[event.id].renderPicking();
                break;
        }
    })

    app.ConnectionStore = rs;
    app.Connection = Connection;
}())
