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
        'WS_GROUPROUTE_UPDATE',
        'WS_SOURCE_CREATE',
        'WS_SOURCE_UPDATE',
        'WS_SOURCE_UPDATE_PARAMS',
        'WS_SOURCE_DELETE',
        'WS_LINK_CREATE',
        'WS_LINK_DELETE',
        'WS_CONNECTION_CREATE',
        'WS_CONNECTION_DELETE',

        // for UI
        'APP_SET_ROOT',
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
        'APP_SELECT_ALL',
        'APP_DESELECT_ALL',
        'APP_REQUEST_CONNECTION',
        'APP_CONNECTION_UPDATE',
        'APP_RENDER_CONNECTIONS',
        'APP_REQUEST_NODE_MOVE',
        'APP_REQUEST_NODE_LABEL',
        'APP_REQUEST_SOURCE_PARAMS',
        'APP_REQUEST_GROUP_IMPORT',
        'APP_TRANSLATE_CONNECTIONS',
        'APP_DESELECT',

        'APP_ADD_NODE_CONNECTION',
        'APP_DELETE_NODE_CONNECTION',
        'APP_REQUEST_ROUTE_UPDATE',
        'APP_DELETE_SELECTION',
        'APP_GROUP_SELECTION',
        'APP_UNGROUP_SELECTION',
        'APP_ROUTE_VISIBLE_PARENT',
        'APP_RENDER',
        'APP_RENDER_CONNECTION_PICKING'
    ];


    app.Actions = {};

    actions.forEach(function(a) {
        app.Actions[a] = a;
    })
})();
