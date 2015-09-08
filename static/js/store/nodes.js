var app = app || {};

(function() {
    // canonical store for all node objects
    var nodes = {};

    // ids for all selected nodes
    var selected = [];
    var groups = [];

    function createInputGeometry(inputs, geometry) {
        return inputs.map(function(id, i) {
            return {
                id: id,
                x: geometry.routeRadius,
                y: (i + .5) * geometry.routeHeight,
                direction: 'input',
                hasValue: false
            }
        });
    }

    function createOutputGeometry(outputs, geometry) {
        return outputs.map(function(id, i) {
            return {
                id: id,
                x: geometry.width + geometry.routeRadius,
                y: (i + .5) * geometry.routeHeight,
                direction: 'output',
            }
        });
    }

    /*
    TODO: implement crank
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
    */

    function canvasMeasureText(text, style) {
        var canvas = document.createElement('canvas');
        var ctx = canvas.getContext('2d');
        ctx.font = style;
        return ctx.measureText(text);
    }


    function Source(data) {}

    function Node(data) {

        this.canvas = document.createElement('canvas');
        // calculate node width
        // potentially make a util so that this can be shared with Group.
        this.data = {};
        this.inputs = [];
        this.outputs = [];
        this.connections = [];
        this.update(data);
        // when the state of the node changes, we need to know what status
        // was set last so that we can clear it. 
        this.lastRouteStatus = null;
        //this.crank = new Crank();

    }

    Node.prototype = Object.create(app.Emitter.prototype);
    Node.constructor = Node;

    Node.prototype.addInput = function(id) {
        this.inputs.push(id);
    }

    Node.prototype.addOutput = function(id) {
        this.outputs.push(id);
    }

    Node.prototype.removeInput = function(id) {
        this.inputs.splice(this.inputs.indexOf(id), 1);
    }

    Node.prototype.removeOutput = function(id) {
        this.outputs.splice(this.outputs.indexOf(id), 1);
    }

    Node.prototype.geometry = function() {
        var inputMeasures = this.inputs.map(function(r) {
            return canvasMeasureText(app.RouteStore.getRoute(r).data.name, '16px helvetica');
        });

        var outputMeasures = this.outputs.map(function(r) {
            return canvasMeasureText(app.RouteStore.getRoute(r).data.name, '16px helvetica');
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
        this.nodeGeometry = {
            width: maxInputWidth + maxOutputWidth + padding.horizontal + routeHeight,
            height: Math.max(this.inputs.length, this.outputs.length) * routeHeight + padding.vertical,
            routeRadius: Math.floor(routeHeight / 2.0),
            routeHeight: routeHeight,
        }

        this.inputsGeometry = createInputGeometry(this.inputs, this.nodeGeometry);
        this.outputsGeometry = createOutputGeometry(this.outputs, this.nodeGeometry);
        this.canvas.width = this.nodeGeometry.width + (this.nodeGeometry.routeRadius * 2);
        this.canvas.height = this.nodeGeometry.height;
    }

    Node.prototype.render = function() {
        this.geometry();

        var ctx = this.canvas.getContext('2d');
        ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        ctx.fillStyle = 'rgba(230,230,230,1)';
        ctx.fillRect(this.nodeGeometry.routeRadius, 0, this.nodeGeometry.width, this.nodeGeometry.height);
        ctx.lineWidth = selected.indexOf(this.data.id) !== -1 ? 2 : 1;
        ctx.strokeStyle = selected.indexOf(this.data.id) !== -1 ? 'rgba(0,0,255,1)' : 'rgba(0,0,0,1)';
        ctx.strokeRect(this.nodeGeometry.routeRadius, 0, this.nodeGeometry.width, this.nodeGeometry.height);

        function renderRoute(routeGeometry, route, geometry) {
            ctx.beginPath();
            ctx.arc(routeGeometry.x, routeGeometry.y, geometry.routeRadius, 0, 2 * Math.PI, false);
            ctx.fillStyle = app.Constants.TypeColors[route.data.type];
            ctx.fill();
            ctx.lineWidth = 1;
            ctx.strokeStyle = 'black';
            ctx.stroke();

            ctx.font = '16px helvetica';
            ctx.textAlign = route.direction === 'input' ? 'left' : 'right';
            ctx.fillStyle = 'black';
            ctx.fillText(route.data.name,
                routeGeometry.x + (route.direction === 'input' ? 1 : -1) * geometry.routeRadius,
                routeGeometry.y + geometry.routeRadius)
        };

        this.inputsGeometry.forEach(function(routeGeometry) {
            var route = app.RouteStore.getRoute(routeGeometry.id);
            renderRoute(routeGeometry, route, this.nodeGeometry);
        }.bind(this))

        this.outputsGeometry.forEach(function(routeGeometry, i) {
            var route = app.RouteStore.getRoute(routeGeometry.id);
            renderRoute(routeGeometry, route, this.nodeGeometry);
        }.bind(this))
    }

    Node.prototype.update = function(data) {
        for (var key in data) {
            this.data[key] = data[key];
        }
        this.position = data.position;

        // re-render this node's connections.
        // TODO: this can probably be ignored in the future in cases where the
        // node is not visible in the current top most group.
        app.Dispatcher.dispatch({
            action: app.Actions.APP_RENDER_CONNECTIONS,
            ids: this.connections,
        });
    }

    Node.prototype.updateStatus = function(event) {
        // good gravy what did this ever do?
        // what is a nodeed?
        // why is id undefiend?
        /*if (event.data.type === 'input' || event.data.type === 'output') {
            var id = event.data.id + '_' + event.data.data + '_' + event.data.type;
            this.lastRouteStatus = id;
            app.Dispatcher.dispatch({
                action: app.Actions.APP_ROUTE_UPDATE_STATUS,
                id: id,
                nodeed: true,
            })
        } else {
            app.Dispatcher.dispatch({
                action: app.Actions.APP_ROUTE_UPDATE_STATUS,
                id: this.lastRouteStatus,
                nodeed: false,
            })
        }*/

        //this.crank.update(event.data.type);
    }

    function Group(data) {
        Node.call(this, data);
        this.children = [];
        this.translation = {
            x: 0,
            y: 0,
        }
    }

    Group.prototype = Object.create(Node.prototype);
    Group.constructor = Group;

    function Selection() {}
    Selection.prototype = Object.create(app.Emitter.prototype);
    Selection.constructor = Selection;

    Selection.prototype.getSelected = function() {
        return selected
    }

    var selection = new Selection();

    function NodeCollection() {}
    NodeCollection.prototype = Object.create(app.Emitter.prototype);
    NodeCollection.constructor = NodeCollection;

    NodeCollection.prototype.getNode = function(id) {
        return nodes[id];
    }

    NodeCollection.prototype.getNodes = function() {
        return Object.keys(nodes);
    }

    NodeCollection.prototype.getSelected = function() {
        return selected;
    }

    // TODO: make it so that this only works for visible nodes
    NodeCollection.prototype.pickNode = function(x, y) {
        var picked = [];
        for (var id in nodes) {
            if (app.Utils.pointInRect(
                nodes[id].position.x,
                nodes[id].position.y,
                nodes[id].canvas.width,
                nodes[id].canvas.height,
                x,
                y
            )) {
                picked.push(parseInt(id));
            }
        }
        return picked;
    }

    NodeCollection.prototype.pickRoute = function(id, x, y) {
        var node = nodes[id];
        x -= node.position.x;
        y -= node.position.y;

        var picked = node.inputsGeometry.filter(function(route) {
            return node.nodeGeometry.routeRadius > app.Utils.distance(route.x, route.y, x, y);
        })

        if (picked.length > 0) {
            return picked[0];
        }

        picked = node.outputsGeometry.filter(function(route) {
            return node.nodeGeometry.routeRadius > app.Utils.distance(route.x, route.y, x, y);
        })

        if (picked.length > 0) {
            return picked[0];
        }

        return null
    }

    NodeCollection.prototype.pickArea = function(x, y, w, h) {
        // should be optimized to only try to select nodes that are
        // 1) in the current visible group
        // 2) ideally within the visible workspace
        var picked = [];

        for (var id in nodes) {
            // center of node 
            var nodeX = nodes[id].position.x + nodes[id].nodeGeometry.routeRadius +
                (.5 * nodes[id].nodeGeometry.width);
            var nodeY = nodes[id].position.y + nodes[id].nodeGeometry.routeRadius +
                (.5 * nodes[id].nodeGeometry.height);
            if (app.Utils.pointInRect(x, y, w, h, nodeX, nodeY)) {
                picked.push(parseInt(id));
            }
        }
        return picked;
    }

    var rs = new NodeCollection();

    function createGroup(node) {
        if (nodes.hasOwnProperty(node.id) === true) {
            console.warn('could not create node:', node.id, ' already exists');
            return
        }

        nodes[node.id] = new Group(node);
        nodes[node.id].render();
    }

    function addChildToGroup(event) {
        if (!nodes.hasOwnProperty(event.child)) {
            console.log(event);
        }
        nodes[event.child].inputs.forEach(function(id) {
            nodes[event.id].addInput(id);
        })
        nodes[event.child].outputs.forEach(function(id) {
            nodes[event.id].addOutput(id);
        })
        nodes[event.id].children.push(event.child);
        nodes[event.id].render();
    }

    function removeChildFromGroup(event) {
        nodes[event.child].inputs.forEach(function(id) {
            nodes[event.id].removeInput(id);
        })

        nodes[event.child].outputs.forEach(function(id) {
            nodes[event.id].removeOutput(id);
        });
        nodes[event.id].render();
    }

    function createBlock(node) {
        if (nodes.hasOwnProperty(node.id) === true) {
            console.warn('could not create node:', node.id, ' already exists');
            return
        }
        // TODO: drop the whole "inputs" and "outputs" part of the schema, put
        // distinction inside the map as a field. 
        var inputs = node.inputs.map(function(input, i) {
            return node.id + '_' + i + '_input';
        });

        var outputs = node.outputs.map(function(output, i) {
            return node.id + '_' + i + '_output';
        });

        // ask the RouteStore to create some routes.
        // TODO: consider using facebook's waitFor() in the future. in that case, 
        // we'd just make RouteStore consume the WS_BLOCK_CREATE message,
        // and have the RouteStore do the job of what is happening here.
        inputs.map(function(id, i) {
            app.Dispatcher.dispatch({
                action: app.Actions.APP_ROUTE_CREATE,
                id: id,
                blockId: node.id,
                index: i,
                direction: 'input',
                data: node.inputs[i]
            });
        });

        outputs.map(function(id, i) {
            app.Dispatcher.dispatch({
                action: app.Actions.APP_ROUTE_CREATE,
                id: id,
                blockId: node.id,
                index: i,
                direction: 'output',
                data: node.outputs[i]
            });
        });

        nodes[node.id] = new Node(node);

        inputs.forEach(function(id) {
            nodes[node.id].addInput(id);
        })

        outputs.forEach(function(id) {
            nodes[node.id].addOutput(id);
        })

        nodes[node.id].render();

    }

    function deleteNode(id) {
        if (nodes.hasOwnProperty(id) === false) {
            console.warn('could not delete node: ', id, ' does not exist');
            return
        }

        // if this id is currently selected, ensure that we remove it and fire
        // selection event
        if (selected.indexOf(id) !== -1) {
            deselect(id);
            selection.emit();
        }

        delete nodes[id]
    }

    function updateNode(node) {
        if (nodes.hasOwnProperty(node.id) === false) {
            console.warn('could not update node: ', node.id, ' does not exist');
            return
        }
        node[node.id] = node;
    }

    function moveNode(id, dx, dy) {
        nodes[id].position.x += dx;
        nodes[id].position.y += dy;
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
            nodes[id].render();
            nodes[id].emit();
        })
    }

    function deselect(id) {
        selected = selected.slice().filter(function(i) {
            return i != id;
        });
        nodes[id].render();
        nodes[id].emit();
    }

    function deselectAll() {
        var toRender = selected.slice();
        selected = [];
        toRender.forEach(function(id) {
            nodes[id].render();
            nodes[id].emit();
        });
    }

    function deleteSelection() {
        // TODO: update this for when we add sources
        selected.forEach(function(id) {
            app.Utils.request(
                'DELETE',
                'blocks/' + id, {},
                null
            )
        })

        // worried this may be async madness
        selected = [];
    }

    function selectMove(dx, dy) {
        // an object containing the set of connections that are effected by 
        // this node move.
        var connections = {};
        selected.forEach(function(id) {
            nodes[id].position.x += dx;
            nodes[id].position.y += dy;
            nodes[id].connections.forEach(function(id) {
                connections[id] = connections.hasOwnProperty(id) ? connections[id] + 1 : 1;
            })
        });

        // when a node moves we need to tell our connectionstore which 
        // connections need to be either translated or re-rendered.
        // yucky message
        if (Object.keys(connections).length == 0) return;
        app.Dispatcher.dispatch({
            action: app.Actions.APP_TRANSLATE_CONNECTIONS,
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

    function finishMove() {
        selected.forEach(function(id) {
            app.Utils.request(
                'PUT',
                'blocks/' + id + '/position', {
                    x: nodes[id].position.x,
                    y: nodes[id].position.y
                },
                null
            )
        })
    }

    function addConnection(event) {
        nodes[event.fromId].connections.push(event.id);
        nodes[event.toId].connections.push(event.id);
    }

    function deleteConnection(event) {
        nodes[event.fromId].connections = nodes[event.fromId].connections.filter(function(id) {
            return !(id == event.id)
        })

        nodes[event.toId].connections = nodes[event.toId].connections.filter(function(id) {
            return !(id == event.id)
        })
    }

    app.Dispatcher.register(function(event) {
        switch (event.action) {
            case app.Actions.WS_GROUP_CREATE:
                createGroup(event.data);
                rs.emit();
                break;
            case app.Actions.WS_GROUP_ADD_CHILD:
                addChildToGroup(event);
                rs.emit();
                break;
            case app.Actions.WS_GROUP_REMOVE_CHILD:
                removeChildFromGroup(event);
                rs.emit();
            case app.Actions.APP_REQUEST_NODE_MOVE:
                finishMove();
                break;
            case app.Actions.WS_BLOCK_CREATE:
                createBlock(event.data);
                rs.emit();
                break;
            case app.Actions.WS_BLOCK_DELETE:
                deleteNode(event.id);
                rs.emit();
                break;
            case app.Actions.APP_MOVE: // this is deprecated
                if (!nodes.hasOwnProperty(event.id)) return;
                moveNode(event.id, event.dx, event.dy);
                break;
            case app.Actions.APP_SELECT_MOVE:
                selectMove(event.dx, event.dy);
                rs.emit();
                break;
            case app.Actions.APP_SELECT:
                if (!nodes.hasOwnProperty(event.id)) return;
                deselectAll();
                selectToggle([event.id]);
                rs.emit();
                selection.emit();
                break;
            case app.Actions.APP_SELECT_TOGGLE:
                selectToggle(event.ids);
                rs.emit();
                selection.emit();
                break;
            case app.Actions.APP_DESELECT_ALL:
                deselectAll();
                rs.emit();
                selection.emit();
                break;
            case app.Actions.WS_BLOCK_UPDATE:
                if (!nodes.hasOwnProperty(event.id)) return;
                nodes[event.id].update(event.data);
                nodes[event.id].emit();
                rs.emit();
                break;
            case app.Actions.WS_BLOCK_UPDATE_STATUS:
                if (!nodes.hasOwnProperty(event.id)) return;
                nodes[event.id].updateStatus(event);
                break;
            case app.Actions.APP_ADD_NODE_CONNECTION:
                addConnection(event);
                break;
            case app.Actions.APP_DELETE_NODE_CONNECTION:
                deleteConnection(event);
                break;
            case app.Actions.APP_DELETE_SELECTION:
                deleteSelection();
                break;
        }
    })

    app.NodeStore = rs;
    app.NodeSelection = selection;
}())
