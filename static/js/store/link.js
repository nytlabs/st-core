var app = app || {};

(function() {
    var routes = {};

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
        routes[route.id] = route;
    }

    function deleteRoute(id) {
        if (routes.hasOwnProperty(id) === false) {
            console.warn('could not delete route: ', id, ' does not exist');
            return
        }
        delete routes[id]
    }

    function updateRoute(route) {
        if (routes.hasOwnProperty(route.id) === false) {
            console.warn('could not update route: ', route.id, ' does not exist');
            return
        }
        route[route.id] = route;
    }

    app.Dispatcher.register(function(event) {
        switch (event.action) {
            case app.Actions.APP_ROUTE_CREATE:
                console.log(event.action);
                createRoute(action.data);
                rs.emit();
                break;
            case app.Actions.APP_ROUTE_DELETE:
                console.log(event.action);
                deleteRoute(action.id);
                rs.emit();
                break;
            case app.Actions.APP_ROUTE_UPDATE_POSITION:
                console.log(event.action);
                rs.emit();
                break;
            case app.Actions.APP_ROUTE_UPDATE_STATUS:
                console.log(event.action);
                rs.emit();
                break;
            case app.Actions.APP_ROUTE_UPDATE_CONNECTED:
                console.log(event.action);
                rs.emit();
                break;
        }
    })

    app.RouteStore = rs;
}())
