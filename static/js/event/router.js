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

    var blocks = {}

    function router(event) {
        if (event.type == 'block' && event.action == 'create') {
            event.data.block.inputs.forEach(function(route, index) {
                app.Dispatcher.dispatch({
                    action: app.Actions.APP_ROUTE_CREATE,
                    id: event.data.block.id + '_' + index + '_input',
                    data: route
                })
            })

            event.data.block.outputs.forEach(function(route, index) {
                app.Dispatcher.dispatch({
                    action: app.Actions.APP_ROUTE_CREATE,
                    id: event.data.block.id + '_' + index + '_output',
                    data: route
                })
            })
        }

        if (event.type == 'block' && event.action == 'info') {
            if (event.data.type === 'receive' || event.data.type === 'broadcast') {
                var s = event.data.type === 'receive' ? 'input' : 'output';
                var id = event.data.id + '_' + event.data.data + '_' + s;
                blocks[event.data.id] = id;

                app.Dispatcher.dispatch({
                    action: app.Actions.APP_ROUTE_UPDATE_STATUS,
                    id: event.data.id + '_' + event.data.data + '_' + s,
                    blocked: true,
                })
            } else {
                if (!blocks.hasOwnProperty(event.data.id)) return;
                app.Dispatcher.dispatch({
                    action: app.Actions.APP_ROUTE_UPDATE_STATUS,
                    id: blocks[event.data.id],
                    blocked: false,
                })
            }
        }


        //app.Dispatcher.dispatch({
        //    action: sanitizeEvent[event.type + '_' + event.action],
        //    data: event.data
        //})
    }

    app.Router = router;
})();
