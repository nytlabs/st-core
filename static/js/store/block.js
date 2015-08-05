var app = app || {};

(function() {

    function Crank() {
        this.status = null;
    }

    Crank.prototype = Object.create(app.Emitter.prototype);
    Crank.constructor = Crank;

    Crank.prototype.update = function(s) {
        if (s != this.status) {
            this.status = s;
            this.emit();
        }
    }

    function Block(data) {

        // TODO: drop the whole "inputs" and "outputs" part of the schema, put
        // distinction inside the map as a field. 
        var inputs = data.inputs.map(function(input, i) {
            return data.id + '_' + i + '_input';
        });

        var outputs = data.outputs.map(function(output, i) {
            return data.id + '_' + i + '_output';
        });

        // ask the RouteStore to create some routes.
        // TODO: consider using facebook's waitFor() in the future. in that case, 
        // we'd just make RouteStore consume the WS_BLOCK_CREATE message,
        // and have the RouteStore do the job of what is happening here.
        inputs.map(function(id, i) {
            app.Dispatcher.dispatch({
                action: app.Actions.APP_ROUTE_CREATE,
                id: id,
                blockId: data.id,
                direction: 'input',
                data: data.inputs[i]
            });
        });

        outputs.map(function(id, i) {
            app.Dispatcher.dispatch({
                action: app.Actions.APP_ROUTE_CREATE,
                id: id,
                blockId: data.id,
                direction: 'output',
                data: data.outputs[i]
            });
        });

        // calculate block width
        // potentially make a util so that this can be shared with Group.
        var inputMeasures = inputs.map(function(r) {
            return app.Utils.measureText(app.RouteStore.getRoute(r).data.name, 'route_label');
        });

        var outputMeasures = outputs.map(function(r) {
            return app.Utils.measureText(app.RouteStore.getRoute(r).data.name, 'route_label');
        });

        var maxInputWidth = inputMeasures.length ? Math.max.apply(null, inputMeasures.map(function(im) {
            return im.width;
        })) : 0;

        var maxOutputWidth = outputMeasures.length ? Math.max.apply(null, outputMeasures.map(function(om) {
            return om.width;
        })) : 0;

        var maxInputHeight = inputMeasures.length ? Math.max.apply(null, inputMeasures.map(function(im) {
            return im.height;
        })) : 0;

        var maxOutputHeight = outputMeasures.length ? Math.max.apply(null, outputMeasures.map(function(om) {
            return om.height;
        })) : 0;

        var routeHeight = Math.max(maxInputHeight, maxOutputHeight);

        var padding = {
            vertical: 6,
            horizontal: 6
        }

        // the following is derived data for use with UI
        this.geometry = {
            width: maxInputWidth + maxOutputWidth + padding.horizontal,
            height: Math.max(inputs.length, outputs.length) * routeHeight + padding.vertical,
            routeRadius: Math.floor(routeHeight / 2.0),
            routeHeight: routeHeight,
        }

        this.inputs = inputs;
        this.outputs = outputs;
        this.position = {
            x: data.position.x,
            y: data.position.y
        }

        // when the state of the block changes, we need to know what status
        // was set last so that we can clear it. 
        this.lastRouteStatus = null;
        this.crank = new Crank();
        this.data = data;
    }

    Block.prototype = Object.create(app.Emitter.prototype);
    Block.constructor = Block;

    Block.prototype.update = function(data) {
        for (var key in data) {
            this.data[key] = data[key];
        }
    }

    Block.prototype.updateStatus = function(event) {
        if (event.data.type === 'input' || event.data.type === 'output') {
            var id = event.data.id + '_' + event.data.data + '_' + event.data.type;
            this.lastRouteStatus = id;
            app.Dispatcher.dispatch({
                action: app.Actions.APP_ROUTE_UPDATE_STATUS,
                id: id,
                blocked: true,
            })
        } else {
            app.Dispatcher.dispatch({
                action: app.Actions.APP_ROUTE_UPDATE_STATUS,
                id: this.lastRouteStatus,
                blocked: false,
            })
        }

        this.crank.update(event.data.type);
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

    function moveBlock(id, dx, dy) {
        blocks[id].position.x += dx;
        blocks[id].position.y += dy;
        blocks[id].emit();
    }

    app.Dispatcher.register(function(event) {

        switch (event.action) {
            case app.Actions.WS_BLOCK_CREATE:
                console.log(event.action);
                createBlock(event.data);
                rs.emit();
                break;
            case app.Actions.WS_BLOCK_DELETE:
                //console.log(event.action);
                //deleteBlock(action.id);
                rs.emit();
                break;
            case app.Actions.APP_MOVE:
                if (!blocks.hasOwnProperty(event.id)) return;
                moveBlock(event.id);
                break
            case app.Actions.WS_BLOCK_UPDATE:
                if (!blocks.hasOwnProperty(event.id)) {
                    console.warn('an action was sent to a non-existant block:', event);
                    return;
                }
                blocks[event.id].update(event.data);
                blocks[event.id].emit();
                break;
            case app.Actions.WS_BLOCK_UPDATE_STATUS:
                blocks[event.id].updateStatus(event);
                break;
        }
    })

    app.BlockStore = rs;
}())
