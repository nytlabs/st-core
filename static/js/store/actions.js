var app = app || {};

(function() {
    var actions = [
        // from websocket
        'WS_BLOCK_CREATE',
        'WS_BLOCK_UPDATE',
        'WS_BLOCK_DELETE',
        'WS_BLOCK_UPDATE_STATUS',
        'WS_GROUP_CREATE',
        'WS_GROUP_UPDATE',
        'WS_GROUP_DELETE',
        'WS_GROUP_ADD_CHILD',
        'WS_GROUP_REMOVE_CHILD',
        'WS_SOUCRE_CREATE',
        'WS_SOURCE_UPDATE',
        'WS_SOURCE_UPDATE_PARAMS',
        'WS_SOURCE_DELETE',
        'WS_LINK_CREATE',
        'WS_LINK_DELETE',
        'WS_CONNECTION_CREATE',
        'WS_CONNECTION_DELETE',

        // for UI
        'APP_ROUTE_CREATE',
        'APP_ROUTE_DELETE',
        'APP_ROUTE_UPDATE',
        'APP_ROUTE_UPDATE_POSITION',
        'APP_ROUTE_UPDATE_STATUS',
        'APP_ROUTE_UPDATE_CONNECTED',
        'APP_MOVE', // to remove
        'APP_SELECT_MOVE',
        'APP_SELECT',
        'APP_SELECT_TOGGLE',
        'APP_DESELECT_ALL',
    ]


    app.Actions = {};

    actions.forEach(function(a) {
        app.Actions[a] = a;
    })
})();
