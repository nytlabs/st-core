var app = app || {};

// TODO: move emit() out of Edge.render();

(function() {
    'use strict';

    var connections = {};
    var selected = {};

    function Edge(data) {
        this.data = data;
        this.position = {
            x: 0,
            y: 0
        };
        this.dirtyPicking = false;
        this.visible = false;
        this.pickColor = app.PickingStore.getColor(this);
        this.canvas = document.createElement('canvas');
        this.pickCanvas = document.createElement('canvas');
        this.geometry();
        this.render();
    }

    Edge.prototype = Object.create(app.Emitter.prototype);
    Edge.prototype.constructor = Edge;

    Edge.prototype.geometry = function() {

        this.dirtyPicking = true;
        var from = app.NodeStore.getVisibleParent(this.idFrom);
        var to = app.NodeStore.getVisibleParent(this.idTo);

        // if either of the routes don't exist any more, simply return
        if (!(this.routeIdFrom in from.routeGeometry)) {
            this.visible = false
            return
        }
        if (!(this.routeIdTo in to.routeGeometry)) {
            this.visible = false
            return
        }

        // TODO: organize -- this is so terribly ugly
        var fromV = app.NodeStore.getNode(app.NodeStore.getRoot()).data.children.indexOf(from.visibleParent) + 1;
        var toV = app.NodeStore.getNode(app.NodeStore.getRoot()).data.children.indexOf(to.visibleParent) + 1;
        var fromH = false
        if (from instanceof app.Group) {
            fromH = from.data.hiddenRoutes.indexOf(this.routeIdFrom) != -1
        }
        var toH = false
        if (to instanceof app.Group) {
            toH = to.data.hiddenRoutes.indexOf(this.routeIdTo) != -1
        }

        this.visible = (!!fromV && !!toV) && (from.data.id != to.data.id) && (!fromH && !toH)

        // buffer accounts for bends in the bezier that may extend outside the
        // bounds of a non-buffered box.
        var buffer = 10;
        var yFrom = from.routeGeometry[this.routeIdFrom].y;
        var xFrom = from.routeGeometry[this.routeIdFrom].x + from.nodeGeometry.routeRadius;

        var yTo = to.routeGeometry[this.routeIdTo].y;
        var xTo = to.routeGeometry[this.routeIdTo].x - to.nodeGeometry.routeRadius + 1;

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

    Edge.prototype.render = function() {
        var ctx = this.canvas.getContext('2d');
        ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        if (!this.visible) {
            this.emit();
            return
        }

        var c = this.curve;

        // TODO: gradient coloring of connections is a test and should be evaluated for performance!
        var fromColor = app.Constants.TypeColors[app.RouteStore.getRoute(this.routeIdFrom).data.type];
        var toColor = app.Constants.TypeColors[app.RouteStore.getRoute(this.routeIdTo).data.type];

        var gradient = ctx.createLinearGradient(c[0], c[1], c[6], c[7]);
        gradient.addColorStop("0", fromColor);
        gradient.addColorStop("1.0", toColor);

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

    Edge.prototype.renderPicking = function() {
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

    function Connection(data) {
        this.routeIdFrom = data.from.id + '_' + data.from.route + '_output';
        this.routeIdTo = data.to.id + '_' + data.to.route + '_input';
        this.idFrom = data.from.id;
        this.idTo = data.to.id;
        Edge.call(this, data);
    }

    Connection.prototype = Object.create(Edge.prototype);
    Connection.prototype.constructor = Connection;

    function Link(data) {
        this.routeIdFrom = 'source_' + data.source.id + '_0_output';
        this.routeIdTo = 'source_' + data.block.id + '_0_input';
        this.idFrom = data.source.id;
        this.idTo = data.block.id;
        Edge.call(this, data);
    }

    Link.prototype = Object.create(Edge.prototype);
    Link.prototype.constructor = Link;

    function EdgeStore() {}
    EdgeStore.prototype = Object.create(app.Emitter.prototype);
    EdgeStore.constructor = EdgeStore;

    EdgeStore.prototype.getEdge = function(id) {
        return connections[id];
    }

    EdgeStore.prototype.getEdges = function() {
        return Object.keys(connections);
    }

    EdgeStore.prototype.getSelected = function() {
        return selected;
    }

    var rs = new EdgeStore();

    function createLink(link) {
        connections[link.id] = new Link(link);

        app.Dispatcher.dispatch({
            action: app.Actions.APP_ADD_NODE_CONNECTION,
            fromId: link.source.id,
            toId: link.block.id,
            id: link.id,
        })
    }

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

    function deleteEdge(id) {
        if (connections.hasOwnProperty(id) === false) {
            console.warn('could not delete connections: ', id, ' does not exist');
            return
        }

        app.Dispatcher.dispatch({
            action: app.Actions.APP_DELETE_NODE_CONNECTION,
            fromId: connections[id].idFrom,
            toId: connections[id].idTo,
            id: id,
        });

        delete connections[id]
    }

    function renderEdges(ids) {
        ids.forEach(function(id) {
            connections[id].geometry();
            connections[id].render()
        })
    }

    function translateEdges(ids, dx, dy) {
        ids.forEach(function(id) {
            connections[id].position.x += dx;
            connections[id].position.y += dy;
        })
    }

    function requestEdge(pickedRoutes) {
        var isLink = pickedRoutes.reduce(function(prev, cur) {
            return !!prev.source && !!cur.source;
        });

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

        if (isLink === true) {
            app.Utils.request(
                'POST',
                'links', {
                    'source': {
                        'id': from.blockId,
                    },
                    'block': {
                        'id': to.blockId
                    }
                },
                null);
        } else {
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
    }

    app.Dispatcher.register(function(event) {
        switch (event.action) {
            case app.Actions.APP_RENDER:
                if (!connections.hasOwnProperty(event.id)) return;
                connections[event.id].render();
                rs.emit();
                break;
            case app.Actions.APP_REQUEST_CONNECTION:
                requestEdge(event.routes);
                break;
            case app.Actions.WS_CONNECTION_CREATE:
                createConnection(event.data);
                rs.emit();
                break;
            case app.Actions.WS_LINK_CREATE:
                createLink(event.data);
                rs.emit();
                break;
            case app.Actions.WS_LINK_DELETE:
            case app.Actions.WS_CONNECTION_DELETE:
                deleteEdge(event.id);
                rs.emit();
                break;
            case app.Actions.APP_RENDER_CONNECTIONS:
                renderEdges(event.ids);
                rs.emit();
                break;
            case app.Actions.APP_TRANSLATE_CONNECTIONS:
                translateEdges(event.translate, event.dx, event.dy);
                renderEdges(event.ids);
                rs.emit();
                break;
            case app.Actions.APP_RENDER_CONNECTION_PICKING:
                connections[event.id].renderPicking();
                break;
        }
    })

    app.Edge = Edge;
    app.EdgeStore = rs;
    app.Connection = Connection;
    app.Link = Link;
}())
