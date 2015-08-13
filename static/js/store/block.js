var app = app || {};

(function() {
    // canonical store for all block objects
    var blocks = {};

    // ids for all selected blocks
    var selected = [];

    function createInputGeometry(inputs, geometry) {
        return inputs.map(function(id, i) {
            return {
                id: id,
                x: geometry.routeRadius,
                y: (i + .5) * geometry.routeHeight,
                direction: 'input'
            }
        });
    }

    function createOutputGeometry(outputs, geometry) {
        return outputs.map(function(id, i) {
            return {
                id: id,
                x: geometry.width + geometry.routeRadius,
                y: (+.5) * geometry.routeHeight,
                direction: 'output'
            }
        });
    }

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

    function Group(data) {}

    function Source(data) {}

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
                index: i,
                direction: 'input',
                data: data.inputs[i]
            });
        });

        outputs.map(function(id, i) {
            app.Dispatcher.dispatch({
                action: app.Actions.APP_ROUTE_CREATE,
                id: id,
                blockId: data.id,
                index: i,
                direction: 'output',
                data: data.outputs[i]
            });
        });

        function canvasMeasureText(text, style) {
            var canvas = document.createElement('canvas');
            var ctx = canvas.getContext('2d');
            ctx.font = style;
            return ctx.measureText(text);
        }

        // calculate block width
        // potentially make a util so that this can be shared with Group.
        var inputMeasures = inputs.map(function(r) {
            return canvasMeasureText(app.RouteStore.getRoute(r).data.name, '');
        });

        var outputMeasures = outputs.map(function(r) {
            return canvasMeasureText(app.RouteStore.getRoute(r).data.name, '');
        });

        var maxInputWidth = inputMeasures.length ? Math.max.apply(null, inputMeasures.map(function(im) {
            return im.width;
        })) : 0;

        var maxOutputWidth = outputMeasures.length ? Math.max.apply(null, outputMeasures.map(function(om) {
            return om.width;
        })) : 0;

        var routeHeight = 15;

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
        this.inputs = createInputGeometry(inputs, this.geometry);
        this.outputs = createOutputGeometry(outputs, this.geometry);
        this.connections = [];
        this.position = {
            x: data.position.x,
            y: data.position.y
        }

        // when the state of the block changes, we need to know what status
        // was set last so that we can clear it. 
        this.lastRouteStatus = null;
        this.crank = new Crank();
        this.data = data;

        this.canvas = document.createElement('canvas');
        this.canvas.width = this.geometry.width + (this.geometry.routeRadius * 2);
        this.canvas.height = this.geometry.height + (this.geometry.routeRadius * 2);


        this.render();
    }

    Block.prototype = Object.create(app.Emitter.prototype);
    Block.constructor = Block;

    Block.prototype.render = function() {
        var ctx = this.canvas.getContext('2d');
        ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        ctx.fillStyle = 'rgba(230,230,230,1)';
        ctx.fillRect(this.geometry.routeRadius, 0, this.geometry.width, this.geometry.height);
        ctx.lineWidth = selected.indexOf(this.data.id) !== -1 ? 2 : 1;
        ctx.strokeStyle = selected.indexOf(this.data.id) !== -1 ? 'rgba(0,0,255,1)' : 'rgba(0,0,0,1)';
        ctx.strokeRect(this.geometry.routeRadius, 0, this.geometry.width, this.geometry.height);

        // move to constants
        var types = {
            'number': 'rgba(170, 255, 0, 1)',
            'object': 'rgba(255, 170, 0, 1)',
            'array': 'rgba(255, 0, 170, 1)',
            'string': 'rgba(170, 0, 255, 1)',
            'boolean': 'rgba(0, 170, 255, 1)',
            'writer': 'rgba(0, 255, 170, 1)',
            'any': 'rgba(255, 255, 255, 1)',
            'error': 'rgba(255, 0, 0, 1)'
        }

        function renderRoute(routeGeometry, route, geometry) {
            ctx.beginPath();
            ctx.arc(routeGeometry.x, routeGeometry.y, geometry.routeRadius, 0, 2 * Math.PI, false);
            ctx.fillStyle = types[route.data.type];
            ctx.fill();
            ctx.lineWidth = 1;
            ctx.strokeStyle = 'black';
            ctx.stroke();
        };

        this.inputs.forEach(function(routeGeometry) {
            var route = app.RouteStore.getRoute(routeGeometry.id);
            renderRoute(routeGeometry, route, this.geometry);
        }.bind(this))

        this.outputs.forEach(function(routeGeometry, i) {
            var route = app.RouteStore.getRoute(routeGeometry.id);
            renderRoute(routeGeometry, route, this.geometry);
        }.bind(this))

        this.emit();
    }

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


    function BlockCollection() {}
    BlockCollection.prototype = Object.create(app.Emitter.prototype);
    BlockCollection.constructor = BlockCollection;

    BlockCollection.prototype.getBlock = function(id) {
        return blocks[id];
    }

    BlockCollection.prototype.getBlocks = function() {
        return Object.keys(blocks);
    }

    BlockCollection.prototype.getSelected = function() {
        return selected;
    }

    // TODO: make it so that this only works for visible blocks
    BlockCollection.prototype.pickBlock = function(x, y) {
        var picked = [];
        for (var id in blocks) {
            if (app.Utils.pointInRect(
                blocks[id].position.x,
                blocks[id].position.y,
                blocks[id].canvas.width,
                blocks[id].canvas.height,
                x,
                y
            )) {
                picked.push(parseInt(id));
            }
        }
        return picked;
    }

    BlockCollection.prototype.pickRoute = function(id, x, y) {
        var block = blocks[id];
        x -= block.position.x;
        y -= block.position.y;

        var picked = block.inputs.filter(function(route) {
            return block.geometry.routeRadius > app.Utils.distance(route.x, route.y, x, y);
        })

        if (picked.length > 0) {
            return picked[0];
        }

        picked = block.outputs.filter(function(route) {
            return block.geometry.routeRadius > app.Utils.distance(route.x, route.y, x, y);
        })

        if (picked.length > 0) {
            return picked[0];
        }

        return null
    }

    BlockCollection.prototype.pickArea = function(x, y, w, h) {
        // should be optimized to only try to select nodes that are
        // 1) in the current visible group
        // 2) ideally within the visible workspace
        var picked = [];

        for (var id in blocks) {
            // center of block 
            var blockX = blocks[id].position.x + blocks[id].geometry.routeRadius +
                (.5 * blocks[id].geometry.width);
            var blockY = blocks[id].position.y + blocks[id].geometry.routeRadius +
                (.5 * blocks[id].geometry.height);
            if (app.Utils.pointInRect(x, y, w, h, blockX, blockY)) {
                picked.push(parseInt(id));
            }
        }
        return picked;
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
        //blocks[id].emit();
    }


    function selectToggle(ids) {
        ids.forEach(function(id) {
            if (selected.indexOf(id) === -1) {
                selected.push(id);
            } else {
                selected = selected.slice().filter(function(i) {
                    return i != id;
                });
            }
            blocks[id].render();
        })
    }

    function deselectAll() {
        var toRender = selected.slice();
        selected = [];
        toRender.forEach(function(id) {
            blocks[id].render();
        });
    }

    function selectMove(dx, dy) {
        // an object containing the set of connections that are effected by 
        // this block move.
        var connections = {};
        selected.forEach(function(id) {
            blocks[id].position.x += dx;
            blocks[id].position.y += dy;
            blocks[id].connections.forEach(function(id) {
                connections[id] = connections.hasOwnProperty(id) ? connections[id] + 1 : 1;
            })
        });

        // when a block moves we need to tell our connectionstore which 
        // connections need to be either translated or re-rendered.
        // yucky message
        app.Dispatcher.dispatch({
            action: app.Actions.APP_RENDER_CONNECTIONS,
            // if only end of a connection is being moved, then we need to 
            // re-render the whole connection
            ids: Object.keys(connections).filter(function(id) {
                return connections[id] === 1
            }),
            // if this connection is referecned more than once then we don't
            // need to be re-rendered, simply translated
            translate: Object.keys(connections).filter(function(id) {
                return connections[id] != 1
            }),
            dx: dx,
            dy: dy
        })
    }

    function updateConnections(event) {
        blocks[event.data.from.id].connections.push(event.data.id);
        blocks[event.data.to.id].connections.push(event.data.id);
    }

    app.Dispatcher.register(function(event) {
        switch (event.action) {
            case app.Actions.WS_BLOCK_CREATE:
                createBlock(event.data);
                rs.emit();
                break;
            case app.Actions.WS_BLOCK_DELETE:
                //console.log(event.action);
                //deleteBlock(action.id);
                rs.emit();
                break;
            case app.Actions.APP_MOVE: // this is deprecated
                if (!blocks.hasOwnProperty(event.id)) return;
                moveBlock(event.id, event.dx, event.dy);
                break;
            case app.Actions.APP_SELECT_MOVE:
                selectMove(event.dx, event.dy);
                rs.emit();
                break;
            case app.Actions.APP_SELECT:
                if (!blocks.hasOwnProperty(event.id)) return;
                deselectAll();
                selectToggle([event.id]);
                rs.emit();
                break;
            case app.Actions.APP_SELECT_TOGGLE:
                selectToggle(event.ids);
                rs.emit();
                break;
            case app.Actions.APP_DESELECT_ALL:
                deselectAll();
                rs.emit();
                break;
            case app.Actions.WS_BLOCK_UPDATE:
                if (!blocks.hasOwnProperty(event.id)) return;
                blocks[event.id].update(event.data);
                blocks[event.id].emit();
                break;
            case app.Actions.WS_BLOCK_UPDATE_STATUS:
                blocks[event.id].updateStatus(event);
                break;
            case app.Actions.WS_CONNECTION_CREATE:
                updateConnections(event);
                //rs.emit();
                break;
        }
    })

    app.BlockStore = rs;
}())
