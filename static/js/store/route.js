var app = app || {};

(function() {
    var routes = {};

    function Route(data) {
        this.data = data;
    }

    Route.prototype = Object.create(app.Emitter.prototype);
    Route.constructor = Route;

    Route.prototype.update = function(data) {
        this.blocked = data;
        this.emit();
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
        routes[route.id] = new Route(route);
    }

    function deleteRoute(id) {
        if (routes.hasOwnProperty(id) === false) {
            console.warn('could not delete route: ', id, ' does not exist');
            return
        }
        delete routes[id]
    }

    app.Dispatcher.register(function(event) {
        switch (event.action) {
            case app.Actions.APP_ROUTE_CREATE:
                console.log(event.action);
                createRoute(event);
                rs.emit();
                break;
            case app.Actions.APP_ROUTE_DELETE:
                console.log(event.action);
                deleteRoute(action.id);
                rs.emit();
                break;
            case app.Actions.APP_ROUTE_UPDATE_POSITION:
                console.log(event.action);
                break;
            case app.Actions.APP_ROUTE_UPDATE_STATUS:
                routes[event.id].update(event.blocked);
                break;
            case app.Actions.APP_ROUTE_UPDATE_CONNECTED:
                console.log(event.action);
                break;
        }
    })

    app.RouteStore = rs;
}())
