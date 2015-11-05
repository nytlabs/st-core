var app = app || {};

(function() {
    var routes = {};

    function Route(data, blockId, direction, index, id, source) {
        this.id = id;
        this.blocked = false;
        this.active = data.hasOwnProperty('value') && data.value != null;
        this.data = data;
        this.blockId = blockId;
        this.visibleParent = blockId;
        this.index = index;
        this.direction = direction;
        this.pickColor = app.PickingStore.getColor(this);
        this.source = !!source ? source : null;
    }

    Route.prototype = Object.create(app.Emitter.prototype);
    Route.prototype.constructor = Route;

    // we've received an update for the value of the route
    Route.prototype.updateData = function(data) {
        this.data.value = data;
        this.active = data !== null;
    }

    // we've received an update for the status of the route
    Route.prototype.update = function(data) {
        this.blocked = data;
    }

    function RouteStore() {}
    RouteStore.prototype = Object.create(app.Emitter.prototype);
    RouteStore.constructor = RouteStore;

    RouteStore.prototype.getRoute = function(id) {
        return routes[id];
    }

    var rs = new RouteStore();

    function createRoute(route) {
        if (routes.hasOwnProperty(route.id) === true) {
            console.warn('could not create route:', route.id, ' already exists');
            return
        }
        routes[route.id] = new Route(route.data, route.blockId, route.direction, route.index, route.id, route.source);
    }

    function deleteRoute(id) {
        if (routes.hasOwnProperty(id) === false) {
            console.warn('could not delete route: ', id, ' does not exist');
            return
        }
        app.PickingStore.removeColor(routes[id].pickColor);
        delete routes[id]
    }

    function requestRouteUpdate(event) {
        var route = routes[event.id];
        app.Utils.request(
            'PUT',
            '/blocks/' + route.blockId + '/routes/' + route.index,
            event.value, {},
            null)
    }

    function setVisibleParent(id, visibleParent) {
        routes[id].visibleParent = visibleParent;
    }

    app.Dispatcher.register(function(event) {
        switch (event.action) {
            case app.Actions.APP_ROUTE_CREATE:
                createRoute(event);
                rs.emit();
                break;
            case app.Actions.APP_ROUTE_DELETE:
                deleteRoute(event.id);
                rs.emit();
                break;
            case app.Actions.APP_ROUTE_UPDATE_POSITION:
                break;
            case app.Actions.APP_ROUTE_UPDATE:
                routes[event.id].updateData(event.value);
                routes[event.id].emit();
                break;
            case app.Actions.APP_REQUEST_ROUTE_UPDATE:
                requestRouteUpdate(event);
                break;
            case app.Actions.APP_ROUTE_VISIBLE_PARENT:
                setVisibleParent(event.id, event.visibleParent);
                break;
            case app.Actions.APP_ROUTE_UPDATE_STATUS:
                // TODO: it's possible for the server to emit statuses of blocks
                // that have not yet been added to the client. this is due to 
                // statuses being emit before the 'list' websocket command is 
                // complete. this should be addressed in the server.
                if (!routes.hasOwnProperty(event.id)) return;
                routes[event.id].update(event.blocked);
                routes[event.id].emit();
                break;
                //case app.Actions.APP_ROUTE_UPDATE_CONNECTED:
                //    updateConnected(event);
                //    break;
        }
    })

    app.RouteStore = rs;
    app.Route = Route;
}())
