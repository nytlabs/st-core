var app = app || {};

(function() {

    // this is a bad thing
    var sanitizeEvent = {
        'block_create': app.Actions.WS_BLOCK_CREATE,
        'block_update': app.Actions.WS_BLOCK_UPDATE,
        'block_delete': app.Actions.WS_BLOCK_DELETE,
        'block_info': app.Actions.WS_BLOCK_UPDATE_STATUS,
        'group_create': app.Actions.WS_GROUP_CREATE,
        'group_update': app.Actions.WS_GROUP_UPDATE,
        'group_delete': app.Actions.WS_GROUP_DELETE,
        'child_create': app.Actions.WS_GROUP_ADD_CHILD,
        'child_delete': app.Actions.WS_GROUP_REMOVE_CHILD,
        'source_create': app.Actions.WS_SOUCE_CREATE,
        'source_update': app.Actions.WS_SOURCE_UPDATE,
        'param_update': app.Actions.WS_SOURCE_UPDATE_PARAMS,
        'source_delete': app.Actions.WS_SOURCE_DELETE,
        'link_create': app.Actions.WS_LINK_CREATE,
        'link_delete': app.Actions.WS_LINK_DELETE,
        'connection_create': app.Actions.WS_CONNECTION_CREATE,
        'connection_delete': app.Actions.WS_CONNECTION_DELETE,
    }

    function router(event) {
        var action = sanitizeEvent[event.type + '_' + event.action];
        switch (event.type) {
            case 'block':
                app.Dispatcher.dispatch({
                    action: action,
                    id: event.data.block.id,
                    data: event.data.block
                });
                break;
            case 'route':
                app.Dispatcher.dispatch({
                    action: app.Actions.APP_ROUTE_UPDATE,
                    id: event.data.id + '_' + event.data.route + '_input',
                    value: event.data.value
                });
                break;
            case 'connection':
                app.Dispatcher.dispatch({
                    action: action,
                    id: event.data.connection.id,
                    data: event.data.connection
                });
                break;
            case 'group':
                break;
            case 'source':
                break;
            case 'link':
                break;
        }
    }

    var ws = new WebSocket('ws://localhost:7071/updates');
    ws.onmessage = function(m) {
        app.Router(JSON.parse(m.data));
    }.bind(this)

    ws.onopen = function() {
        ws.send('list');
    }


    app.Router = router;
})();
