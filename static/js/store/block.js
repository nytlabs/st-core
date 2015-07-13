var app = app || {};

(function() {

    function Block(data) {

        // TODO: drop the whole "inputs" and "outputs" part of the schema, put
        // distinction inside the map as a field. 
        var inputs = data.inputs.map(function(input, i) {
            return data.id + '_' + i + '_input';
        });

        var outputs = data.outputs.map(function(output, i) {
            return data.id + '_' + i + '_output';
        });

        // create the routes in the route store.
        inputs.map(function(id, i) {
            app.Dispatcher.dispatch({
                action: app.Actions.APP_ROUTE_CREATE,
                id: id,
                data: data.inputs[i]
            });
        });

        outputs.map(function(id, i) {
            app.Dispatcher.dispatch({
                action: app.Actions.APP_ROUTE_CREATE,
                id: id,
                data: data.outputs[i]
            });
        });

        // holds the ID of the last blocked route.
        this.lastBlockedRoute = null;

        /*        if (event.type == 'block' && event.action == 'create') {
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
                }*/


        this.data = data;
    }

    Block.prototype = Object.create(app.Emitter.prototype);
    Block.constructor = Block;

    Block.prototype.update = function(data) {
        for (var key in data) {
            this.data[key] = data[key];
        }
    }

    Block.prototype.updateStatus = function(status) {

    }

    var blocks = {};

    function BlockCollection() {}
    BlockCollection.prototype = Object.create(app.Emitter.prototype);
    BlockCollection.constructor = BlockCollection;

    BlockCollection.prototype.getBlock = function(id) {
        return blocks[id];
    }

    var rs = new BlockCollection();

    function createBlock(block) {
        if (blocks.hasOwnProperty(block.id) === true) {
            console.warn('could not create block:', block.id, ' already exists');
            return
        }
        blocks[block.id] = new Block(block);
    }

    function deleteBlock(id) {
        if (blocks.hasOwnProperty(id) === false) {
            console.warn('could not delete block: ', id, ' does not exist');
            return
        }
        delete blocks[id]
    }

    function updateBlock(block) {
        if (blocks.hasOwnProperty(block.id) === false) {
            console.warn('could not update block: ', block.id, ' does not exist');
            return
        }
        block[block.id] = block;
    }

    app.Dispatcher.register(function(event) {
        switch (event.action) {
            case app.Actions.WS_BLOCK_CREATE:
                console.log(event.action);
                createBlock(event.data);
                rs.emit();
                break;
            case app.Actions.WS_BLOCK_DELETE:
                console.log(event.action);
                deleteBlock(action.id);
                rs.emit();
                break;
            case app.Actions.WS_BLOCK_UPDATE:
                blocks[event.id].update(event.data);
                blocks[event.id].emit();
                break;
            case app.Actions.WS_BLOCK_UPDATE_STATUS:
                blocks[event.id].updateStatus(event);
                break;
        }
    })

    app.BlockCollection = rs;
}())
