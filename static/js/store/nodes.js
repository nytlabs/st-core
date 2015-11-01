var app = app || {};

// TODO: remove inputs/outputs ascending
// TODO: node emit event for route removal.

(function() {
    // canonical store for all node objects
    var nodes = {};

    var root = null;

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
        this.pickCanvas = document.createElement('canvas');
        this.pickColor = app.PickingStore.getColor(this);
        // calculate node width
        // potentially make a util so that this can be shared with Group.
        this.visibleParent = null;
        this.parent = null;
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
    Node.prototype.constructor = Node;

    Node.prototype.addInput = function(id) {
        this.inputs.push(id);
        var route = app.RouteStore.getRoute(id);
        route.addListener(this.renderAndEmit.bind(this));
    }

    Node.prototype.addOutput = function(id) {
        this.outputs.push(id);
    }

    Node.prototype.removeInput = function(id) {
        this.inputs.splice(this.inputs.indexOf(id), 1);
        var route = app.RouteStore.getRoute(id);
        route.removeListener(this.renderAndEmit.bind(this));
    }

    Node.prototype.renderAndEmit = function() {
        // really don't like this -- somewhat confusing handling of events
        //
        // this function is fired when a route has emitted an update event.
        // this means we need to re-render that route, because somehow its
        // state has changed (had a value set, etc). 
        //
        // we bind this function (renderAndEmit) to that event, and whenever
        // that happens we call a render on the block, and then we tell the
        // NodeStore to emit an update event. Since the canvasgraph component
        // is listening to NodeStore events, the render is fully propagated.
        //
        // TODO: an optimized version of this would check to see if the block is
        // visible. if it's not then don't do anything
        this.render();
        app.NodeStore.emit();
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
            width: Math.floor(maxInputWidth + maxOutputWidth + padding.horizontal + routeHeight),
            height: Math.floor(Math.max(this.inputs.length, this.outputs.length) * routeHeight + padding.vertical),
            routeRadius: Math.floor(routeHeight / 2.0),
            routeHeight: routeHeight,
        }

        this.inputsGeometry = createInputGeometry(this.inputs, this.nodeGeometry);
        this.outputsGeometry = createOutputGeometry(this.outputs, this.nodeGeometry);
        this.canvas.width = this.nodeGeometry.width + (this.nodeGeometry.routeRadius * 2) + 1; // magic number for buffer...
        this.canvas.height = this.nodeGeometry.height + 1; // 1 is magic buffer number :(

        this.pickCanvas.width = this.canvas.width;
        this.pickCanvas.height = this.canvas.height;
    }

    Node.prototype.render = function() {
        // TODO: move colors to constants
        this.geometry();

        var ctx = this.canvas.getContext('2d');
        // seriously? http://www.mobtowers.com/html5-canvas-crisp-lines-every-time/
        ctx.translate(.5, .5);
        ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        ctx.fillStyle = this instanceof Group ? 'rgba(210,230,255,1)' : 'rgba(230,230,230,1)';
        ctx.fillRect(this.nodeGeometry.routeRadius, 0, this.nodeGeometry.width, this.nodeGeometry.height);
        ctx.lineWidth = app.SelectionStore.isSelected(this) ? 2 : 1;
        ctx.strokeStyle = app.SelectionStore.isSelected(this) ? 'rgba(0,0,255,1)' : 'rgba(150,150,150,1)';
        ctx.strokeRect(this.nodeGeometry.routeRadius, 0, this.nodeGeometry.width, this.nodeGeometry.height);

        // now to do the picking buffer.
        var pctx = this.pickCanvas.getContext('2d');
        pctx.translate(.5, .5);
        pctx.clearRect(0, 0, this.pickCanvas.width, this.pickCanvas.height);
        pctx.fillStyle = this.pickColor;
        pctx.fillRect(this.nodeGeometry.routeRadius, 0, this.nodeGeometry.width, this.nodeGeometry.height);

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
                routeGeometry.y + geometry.routeRadius);

            if (route.direction === 'input' && route.data.value !== null) {
                ctx.beginPath();
                ctx.fillStlye = 'rgba(100,100,100,1)';
                ctx.arc(routeGeometry.x, routeGeometry.y, 4, 0, 2 * Math.PI, false);
                ctx.fill();
            }

            pctx.beginPath();
            pctx.arc(routeGeometry.x, routeGeometry.y, geometry.routeRadius, 0, 2 * Math.PI, false);
            pctx.fillStyle = route.pickColor;
            pctx.fill();
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

    var selection = new Selection();

    function NodeCollection() {}
    NodeCollection.prototype = Object.create(app.Emitter.prototype);
    NodeCollection.constructor = NodeCollection;

    NodeCollection.prototype.getNode = function(id) {
        return nodes[id];
    }

    // given a block id, gets the top-most group to be visible in the tree
    NodeCollection.prototype.getVisibleParent = function(id) {
        return nodes[nodes[id].visibleParent];
    }

    // getNodes returns all nodes that should be on-screen
    NodeCollection.prototype.getNodes = function() {
        // if for some reason we don't have a root set, return all nodes
        return root === null ? Object.keys(nodes) : nodes[root].children;
    }

    var rs = new NodeCollection();

    function createGroup(node) {
        if (nodes.hasOwnProperty(node.id) === true) {
            console.warn('could not create node:', node.id, ' already exists');
            return
        }

        nodes[node.id] = new Group(node);
        nodes[node.id].render();

        // set group 0 as our current parent when we recieve it
        // TODO: in the future, we may want multiple patterns with a 'null'
        // parent, thus making them 'root' groups the same way that group 0
        // is. Currently, all groups descend from a single root, and there 
        // isn't necessarily a reason for that.
        if (node.id === 0) setRoot(node.id);
    }

    function getVisibleParent(id) {
        var node = nodes[id];
        while (nodes[root].children.indexOf(node.data.id) === -1) {
            node = nodes[node.parent];
        }
        return node.data.id;
    }

    function setRoot(id) {
        root = id;

        Object.keys(nodes).filter(function(node) {
            return node !== "0"
        }).forEach(function(node) {
            // TODO: setvisibleparents for all children of new root. 
        })

        // we need to re-render right now!
        app.NodeStore.emit();
    }

    function addInputAscending(id, routeId) {
        nodes[id].addInput(routeId);
        if (nodes[id].parent !== null) {
            addInputAscending(nodes[id].parent, routeId);
        }
    }

    function addOutputAscending(id, routeId) {
        nodes[id].addOutput(routeId);
        if (nodes[id].parent !== null) {
            addOutputAscending(nodes[id].parent, routeId);
        }
    }

    function setVisibleParentDescending(id, parent) {
        nodes[id].visibleParent = parent;
        if (nodes[id] instanceof Group) {
            nodes[id].children.forEach(function(childId) {
                setVisibleParentDescending(childId, parent);
            })
        }
    }

    function addChildToGroup(event) {
        nodes[event.id].children.push(event.child);
        nodes[event.child].parent = event.id;

        // add routes to all parent nodes
        nodes[event.child].inputs.forEach(function(routeId) {
            addInputAscending(event.id, routeId);
        });

        nodes[event.child].outputs.forEach(function(routeId) {
            addOutputAscending(event.id, routeId);
        });

        nodes[event.child].connections.forEach(function(connId) {
            addConnectionAscending(event.child, connId);
        })

        // find the top-most visible node and store that id in all child nodes.
        var visibleParent = getVisibleParent(event.child);
        setVisibleParentDescending(event.child, visibleParent);


        nodes[visibleParent].inputs.forEach(function(route) {
            app.Dispatcher.dispatch({
                action: app.Actions.APP_ROUTE_VISIBLE_PARENT,
                id: route,
                visibleParent: visibleParent,
            })
        })

        nodes[visibleParent].outputs.forEach(function(route) {
            app.Dispatcher.dispatch({
                action: app.Actions.APP_ROUTE_VISIBLE_PARENT,
                id: route,
                visibleParent: visibleParent,
            })
        })

        // only need to render the top-most visible node.
        nodes[visibleParent].render();

        app.Dispatcher.dispatch({
            action: app.Actions.APP_RENDER_CONNECTIONS,
            ids: nodes[visibleParent].connections,
        });
    }

    function removeChildFromGroup(event) {
        nodes[event.child].inputs.forEach(function(id) {
            nodes[event.id].removeInput(id);
        })

        nodes[event.child].outputs.forEach(function(id) {
            nodes[event.id].removeOutput(id);
        });

        nodes[event.id].children.splice(nodes[event.id].children.indexOf(event.child), 1);

        //if our group is a child of the current root, then we need to render
        if (nodes[root].children.indexOf(event.id) !== -1) {
            nodes[event.id].render();
        }
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
        if (app.SelectionStore.isSelected(nodes[id]) !== -1) {
            app.Dispatcher.dispatch({
                action: app.Actions.APP_DESELECT,
                id: nodes[id],
            });
        }

        // remove the picking color from the store so that we can re-use it later
        app.PickingStore.removeColor(nodes[id].pickColor);

        // if this is a block then we need to remove its routes
        // if it is a group we do nothing.
        if (nodes[id].constructor === Node) {
            nodes[id].inputs.forEach(function(route) {
                app.Dispatcher.dispatch({
                    action: app.Actions.APP_ROUTE_DELETE,
                    id: route
                })
            })

            nodes[id].outputs.forEach(function(route) {
                app.Dispatcher.dispatch({
                    action: app.Actions.APP_ROUTE_DELETE,
                    id: route
                })
            })
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

    function deleteSelection() {
        // TODO: update this for when we add sources
        app.SelectionStore.getIdsByKind(Node).forEach(function(id) {
            var type;
            if (nodes[id] instanceof Node) type = 'blocks'
            if (nodes[id] instanceof Group) type = 'groups';

            app.Utils.request(
                'DELETE',
                type + '/' + id, {},
                null
            )
        });
    }

    function selectMove(dx, dy) {
        var connections = {};
        app.SelectionStore.getIdsByKind(Node).forEach(function(id) {
            nodes[id].position.x += dx;
            nodes[id].position.y += dy;
            nodes[id].connections.forEach(function(id) {
                connections[id] = connections.hasOwnProperty(id) ? connections[id] + 1 : 1;
            });
        });

        app.Dispatcher.dispatch({
            action: app.Actions.APP_RENDER_CONNECTIONS,
            ids: Object.keys(connections),
        });
    }

    function nodeType(id) {
        if (nodes[id] instanceof Group) return 'group';
        return 'block';
    }

    function finishMove() {
        app.SelectionStore.getIdsByKind(Node).forEach(function(id) {
            app.Utils.request(
                'PUT',
                nodeType(id) + 's/' + id + '/position', {
                    x: nodes[id].position.x,
                    y: nodes[id].position.y
                },
                null
            )
        })
    }

    function selectGroup() {
        var selected = app.SelectionStore.getIdsByKind(Node)
        if (selected.length === 0) return;

        var position = {
            x: 0,
            y: 0
        }

        selected.forEach(function(id) {
            position.x += nodes[id].position.x;
            position.y += nodes[id].position.y;
        })

        position.x /= selected.length;
        position.y /= selected.length;

        position.x = Math.round(position.x);
        position.y = Math.round(position.y);

        app.Utils.request(
            'POST',
            'groups', {
                parent: root,
                children: selected,
                position: position
            }, null)
    }

    function selectUnGroup() {
        var children = [];
        app.SelectionStore.getIdsByKind(Node).forEach(function(id) {
            if (nodes[id] instanceof Group) {
                children = children.concat(nodes[id].children);
            }
        })

        function jobsDone() {
            if (children.length === 0) {
                app.SelectionStore.getIdsByKind(Node).forEach(function(id) {
                    app.Utils.request(
                        'DELETE',
                        'groups/' + id, {}, function() {}
                    )
                })
            }
        }

        for (var i = 0; i < children.length; i++) {
            app.Utils.request(
                'PUT',
                'groups/' + root + '/children/' + children[i], null,
                function() {
                    children.splice(children.indexOf(children[i]), 1);
                    jobsDone();
                }
            )
        }
    }

    function addConnectionAscending(id, connectionId) {
        if (nodes[id].connections.indexOf(connectionId) === -1) {
            nodes[id].connections.push(connectionId);
        }

        if (nodes[id].parent !== null) {
            addConnectionAscending(nodes[id].parent, connectionId);
        }
    }

    function addConnection(event) {
        addConnectionAscending(event.fromId, event.id);
        addConnectionAscending(event.toId, event.id);
    }

    function removeConnectionAscending(id, connectionId) {
        if (nodes[id].connections.indexOf(connectionId) !== -1) {
            nodes[id].connections.splice(nodes[id].connections.indexOf(connectionId), 1);
        }
        if (nodes[id].parent !== null) {
            removeConnectionAscending(nodes[id].parent, connectionId);
        }
    }

    function deleteConnection(event) {
        removeConnectionAscending(event.fromId, event.id);
        removeConnectionAscending(event.toId, event.id);
    }

    app.Dispatcher.register(function(event) {
        switch (event.action) {
            case app.Actions.APP_GROUP_SELECTION:
                selectGroup();
                break;
            case app.Actions.APP_UNGROUP_SELECTION:
                selectUnGroup();
                break;
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
                break;
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
            case app.Actions.APP_SELECT_MOVE:
                selectMove(event.dx, event.dy);
                rs.emit();
                break;
            case app.Actions.WS_GROUP_UPDATE:
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
            case app.Actions.APP_RENDER:
                if (!nodes.hasOwnProperty(event.id)) return;
                nodes[event.id].render();
                rs.emit();
                break;
        }
    })
    app.Node = Node;
    app.NodeStore = rs;
    app.NodeSelection = selection;
}())
