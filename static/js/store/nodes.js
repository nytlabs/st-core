var app = app || {};

// TODO: remove inputs/outputs ascending
// TODO: node emit event for route removal.

(function() {
    // canonical store for all node objects
    var nodes = {};

    var root = null;

    var tree = null;
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

    function Node(data) {
        this.canvas = document.createElement('canvas');
        this.pickCanvas = document.createElement('canvas');
        this.pickColor = app.PickingStore.getColor(this);
        // calculate node width
        // potentially make a util so that this can be shared with Group.
        this.visibleParent = null;
        this.parent = null;
        this.data = {};
        this.routes = [];
        this.routeGeometry = {};
        this.connections = [];
        this.update(data);
        // when the state of the node changes, we need to know what status
        // was set last so that we can clear it. 
        this.lastRouteStatus = null;
        //this.crank = new Crank();
    }

    Node.prototype = Object.create(app.Emitter.prototype);
    Node.prototype.constructor = Node;


    Node.prototype.addRoute = function(id) {
        this.routes.push(id);
        var route = app.RouteStore.getRoute(id);
        // if this route is an input we need to listen for updates
        if (route.direction === 'input') {
            route.addListener(this.renderAndEmit.bind(this));
        }
    }

    Node.prototype.removeRoute = function(id) {
        this.routes.splice(this.routes.indexOf(id), 1);
        var route = app.RouteStore.getRoute(id);
        if (route.direction === 'input') {
            route.removeListener(this.renderAndEmit.bind(this));
        }
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

    Node.prototype.geometry = function() {
        var routeHeight = 15;
        var routeRadius = Math.floor(routeHeight / 2.0);

        var padding = {
            top: 0,
            bottom: 7,
            middle: 10,
            side: 3
        }

        var widths = {
            input: 0,
            output: 0
        }

        var counts = {
            input: 0,
            output: 0
        }

        // this is terrible.
        // if the node is a source, don't show the type
        // if the node is a block, default to the type
        // if the node is a source, group, or block, allow label to override type.
        this.label = ''
        if (!(this instanceof Group) && !(this instanceof Source)) {
            this.label = this.data.type;
        }
        this.label = !!this.data.label ? this.data.label : this.label;
        if (this.label.length) {
            padding.top += 20;
        }

        var labelWidth = canvasMeasureText(this.label, 'Bold 14px helvetica').width;

        this.routes.forEach(function(id) {
            if (this instanceof Group) {
                if (this.data.hiddenRoutes.indexOf(id) != -1) {
                    return
                }
            }
            var route = app.RouteStore.getRoute(id);
            var width = canvasMeasureText(route.data.name, '14px helvetica').width;
            if (width > widths[route.direction]) {
                widths[route.direction] = width;
            }
            counts[route.direction] += 1;
        }.bind(this));

        var routeWidth = widths.input + widths.output + padding.middle + routeHeight + padding.side * 2;
        var maxWidth = Math.max(routeWidth, labelWidth + padding.side * 2);

        // the following is derived data for use with UI
        this.nodeGeometry = {
            width: Math.floor(maxWidth),
            height: Math.floor(Math.max(counts.input, counts.output) * routeHeight + padding.top + padding.bottom),
            routeRadius: routeRadius,
            routeHeight: routeHeight,
        }

        counts = {
            input: 0,
            output: 0
        }

        this.routes.forEach(function(id) {
            if (this instanceof Group) {
                if (this.data.hiddenRoutes.indexOf(id) != -1) {
                    return
                }
            }
            var route = app.RouteStore.getRoute(id);
            var xOffset = route.direction === 'input' ? padding.side : this.nodeGeometry.width - padding.side;
            var x = xOffset + routeRadius;
            var y = ((counts[route.direction]) + .5) * routeHeight + padding.top;

            // this is for connections;
            this.routeGeometry[id] = {
                id: id,
                index: counts[route.direction]++,
                x: x,
                y: y,
                direction: route.direction,
            };
        }.bind(this))

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
        var fillStyle = 'rgba(230,230,230,1)';
        if (this instanceof Group) fillStyle = 'rgba(210,230,255,1)';
        if (this instanceof Source) fillStyle = 'rgba(255,230,210,1)';
        ctx.fillStyle = fillStyle;
        ctx.fillRect(this.nodeGeometry.routeRadius, 0, this.nodeGeometry.width, this.nodeGeometry.height);
        ctx.lineWidth = app.SelectionStore.isSelected(this) ? 2 : 1;
        ctx.strokeStyle = app.SelectionStore.isSelected(this) ? 'rgba(0,0,255,1)' : 'rgba(150,150,150,1)';
        ctx.strokeRect(this.nodeGeometry.routeRadius, 0, this.nodeGeometry.width, this.nodeGeometry.height);
        ctx.font = 'Bold 14px helvetica';
        ctx.textAlign = 'center';
        ctx.fillStyle = 'black';
        ctx.fillText(this.label, this.nodeGeometry.routeRadius + this.nodeGeometry.width * .5, 16);

        // now to do the picking buffer.
        var pctx = this.pickCanvas.getContext('2d');
        pctx.translate(.5, .5);
        pctx.clearRect(0, 0, this.pickCanvas.width, this.pickCanvas.height);
        pctx.fillStyle = this.pickColor;
        pctx.fillRect(this.nodeGeometry.routeRadius, 0, this.nodeGeometry.width, this.nodeGeometry.height);

        ctx.font = '14px helvetica';
        this.routes.forEach(function(id, i) {
            if (this instanceof Group) {
                if (this.data.hiddenRoutes.indexOf(id) != -1) {
                    return
                }
            }
            var route = app.RouteStore.getRoute(id);
            var x = this.routeGeometry[id].x;
            var y = this.routeGeometry[id].y;

            ctx.beginPath();
            ctx.arc(x, y, this.nodeGeometry.routeRadius, 0, 2 * Math.PI, false);
            ctx.fillStyle = app.Constants.TypeColors[route.data.type];
            ctx.fill();
            ctx.lineWidth = 1;
            ctx.strokeStyle = 'black';
            ctx.stroke();

            ctx.textAlign = route.direction === 'input' ? 'left' : 'right';
            ctx.fillStyle = 'black';
            ctx.fillText(route.data.name,
                x + (route.direction === 'input' ? 10 : -10),
                y + this.nodeGeometry.routeRadius);

            if (route.direction === 'input' && route.data.value !== null) {
                ctx.beginPath();
                ctx.fillStlye = 'rgba(100,100,100,1)';
                ctx.arc(x, y, 4, 0, 2 * Math.PI, false);
                ctx.fill();
            }

            pctx.beginPath();
            pctx.arc(x, y, this.nodeGeometry.routeRadius, 0, 2 * Math.PI, false);
            pctx.fillStyle = route.pickColor;
            pctx.fill();
        }.bind(this));
    }

    Node.prototype.update = function(data) {
        for (var key in data) {
            this.data[key] = data[key];
        }

        this.geometry();
        this.render();

        if (!!data.position) {
            this.position = data.position;
        }

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
        this.data.children = [];
        this.translation = {
            x: 0,
            y: 0,
        }

        this.tree = {
            id: this.data.id,
            children: []
        };
    }

    Group.prototype = Object.create(Node.prototype);
    Group.constructor = Group;

    // hideRoute and showRoute are intended to idempotent
    Group.prototype.hideRoute = function(routeId) {
        // only hide routes that aren't already hidden
        if (this.data.hiddenRoutes.indexOf(routeId) == -1) {
            this.data.hiddenRoutes.push(routeId)
        }
    }

    Group.prototype.showRoute = function(routeId) {
        // only show a route that has been hidden
        if (this.data.hiddenRoutes.indexOf(routeId) != -1) {
            this.data.hiddenRoutes.splice(this.data.hiddenRoutes.indexOf(routeId), 1); //TODO use sets
        }
    }

    function Source(data) {
        Node.call(this, data);
    }

    Source.prototype = Object.create(Node.prototype);
    Source.constructor = Source;

    Source.prototype.updateParams = function(params) {
        this.data.params[params.params] = params.value
        this.data.params.forEach(function(param) {
            if (param.name != params.param) return;
            param.value = params.value
        })
    }

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
        return root === null ? Object.keys(nodes) : nodes[root].data.children;
    }

    NodeCollection.prototype.setRoot = function(id) {
        setRoot(parseInt(id));
    }

    NodeCollection.prototype.getRoot = function() {
        return root;
    }

    NodeCollection.prototype.getTree = function() {
        return tree;
    }

    var rs = new NodeCollection();


    function getVisibleParent(id) {
        var node = nodes[id];
        while (node.parent !== null && nodes[root].data.children.indexOf(node.data.id) === -1) {
            node = nodes[node.parent];
        }
        return node.data.id;
    }

    function setRoot(id) {
        var oldConns = [];
        if (root !== null) {
            nodes[root].data.children.forEach(function(id) {
                oldConns = oldConns.concat(nodes[id].connections);
            })
        }

        root = id;
        nodes[root].data.children.forEach(function(id) {
            var woop = getVisibleParent(id);
            setVisibleParentDescending(id, woop);
            nodes[id].routes.forEach(function(route) {
                app.Dispatcher.dispatch({
                    action: app.Actions.APP_ROUTE_VISIBLE_PARENT,
                    id: route,
                    visibleParent: id
                })
            });
            app.Dispatcher.dispatch({
                action: app.Actions.APP_RENDER_CONNECTIONS,
                ids: nodes[id].connections,
            });
        });

        app.Dispatcher.dispatch({
            action: app.Actions.APP_RENDER_CONNECTIONS,
            ids: oldConns,
        });

        app.NodeStore.emit();
    }

    function addRouteAscending(id, routeId) {
        nodes[id].addRoute(routeId);
        if (nodes[id].parent !== null) {
            addRouteAscending(nodes[id].parent, routeId);
        }
    }

    function removeRouteAscending(id, routeId) {
        nodes[id].removeRoute(routeId);
        if (nodes[id].parent !== null) {
            removeRouteAscending(nodes[id].parent, routeId);
        }
    }

    function hideRouteAscending(id, routeId) {
        nodes[id].hideRoute(routeId)
        if (nodes[id].parent !== null) {
            hideRouteAscending(nodes[id].parent, routeId);
        }
    }

    function showRouteAscending(id, routeId) {
        nodes[id].showRoute(routeId)
        if (nodes[id].parent !== null) {
            showRouteAscending(nodes[id].parent, routeId);
        }
    }

    function setVisibleParentDescending(id, parent) {
        nodes[id].visibleParent = parent;
        if (nodes[id] instanceof Group) {
            nodes[id].data.children.forEach(function(childId) {
                setVisibleParentDescending(childId, parent);
            })
        }
    }

    function rebuildAscending(id) {
        nodes[id].geometry();
        nodes[id].render();
        if (nodes[id].parent !== null) {
            rebuildAscending(nodes[id].parent);
        }
    }

    function addChildToGroup(event) {
        nodes[event.id].data.children.push(event.child);
        nodes[event.child].parent = event.id;

        // add routes to all parent nodes
        nodes[event.child].routes.forEach(function(routeId) {
            addRouteAscending(event.id, routeId);
        });

        // inherit any hidden routes
        if (nodes[event.child] instanceof Group) {
            nodes[event.child].data.hiddenRoutes.forEach(function(routeId) {
                hideRouteAscending(event.id, routeId);
            });
        }

        nodes[event.child].connections.forEach(function(connId) {
            addConnectionAscending(event.child, connId);
        });

        rebuildAscending(event.id);

        // find the top-most visible node and store that id in all child nodes.
        var visibleParent = getVisibleParent(event.child);
        setVisibleParentDescending(event.child, visibleParent);

        nodes[visibleParent].routes.forEach(function(route) {
            app.Dispatcher.dispatch({
                action: app.Actions.APP_ROUTE_VISIBLE_PARENT,
                id: route,
                visibleParent: visibleParent,
            });
        });

        // only need to render the top-most visible node.
        nodes[visibleParent].render();

        app.Dispatcher.dispatch({
            action: app.Actions.APP_RENDER_CONNECTIONS,
            ids: nodes[visibleParent].connections,
        });

        if (nodes[event.child] instanceof Group) {
            nodes[event.id].tree.children.push(nodes[event.child].tree);
        }
    }

    function removeChildFromGroup(event) {
        nodes[event.child].routes.forEach(function(id) {
            removeRouteAscending(event.id, id);
        });

        nodes[event.child].connections.forEach(function(connId) {
            removeConnectionAscending(event.id, connId);
        });

        rebuildAscending(event.id);

        nodes[event.id].data.children.splice(nodes[event.id].data.children.indexOf(event.child), 1);

        //if our group is a child of the current root, then we need to render
        if (nodes[root].data.children.indexOf(event.id) !== -1) {
            nodes[event.id].render();
        }

        if (nodes[event.child] instanceof Group) {
            nodes[event.id].tree.children = nodes[event.id].tree.children.filter(function(leaf) {
                return leaf.id != event.child;
            })
        }
    }

    function updateGroupRouteVisibility(event) {
        if (event.data.isVisible) {
            nodes[event.id].showRoute(event.data.route.id)
            showRouteAscending(event.id, event.data.route.id)
        } else {
            nodes[event.id].hideRoute(event.data.route.id)
            hideRouteAscending(event.id, event.data.route.id)
        }
        rebuildAscending(event.id);
        nodes[event.id].geometry();
        nodes[event.id].render();
    }

    function createSource(node) {
        nodes[node.id] = new Source(node);

        var routes = [{
            direction: 'output',
            index: 0
        }]

        routes.forEach(function(e) {
            var id = 'source_' + node.id + '_' + e.index + '_' + e.direction;
            app.Dispatcher.dispatch({
                action: app.Actions.APP_ROUTE_CREATE,
                id: id,
                blockId: node.id,
                index: e.index,
                direction: e.direction,
                data: {
                    name: nodes[node.id].data.type,
                    type: 'any',
                    value: null
                },
                source: nodes[node.id].data.type
            });
            nodes[node.id].addRoute(id);
        })

        nodes[node.id].render();
    }

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
        if (node.id === 0) {
            setRoot(node.id);
            tree = nodes[node.id].tree;
        }
    }

    function createBlock(node) {
        if (nodes.hasOwnProperty(node.id) === true) {
            console.warn('could not create node:', node.id, ' already exists');
            return
        }

        nodes[node.id] = new Node(node);

        var inputs = node.inputs.map(function(input, i) {
            return {
                direction: 'input',
                index: i
            }
        })

        var outputs = node.outputs.map(function(output, i) {
            return {
                direction: 'output',
                index: i
            }
        })

        var routes = inputs.concat(outputs);

        routes.forEach(function(e) {
            var id = node.id + '_' + e.index + '_' + e.direction;
            app.Dispatcher.dispatch({
                action: app.Actions.APP_ROUTE_CREATE,
                id: id,
                blockId: node.id,
                index: e.index,
                direction: e.direction,
                data: node[e.direction + 's'][e.index]
            });
            nodes[node.id].addRoute(id);
        })

        // if this block is associated with a source
        if (nodes[node.id].data.source !== null) {
            var id = 'source_' + node.id + '_0_input';
            app.Dispatcher.dispatch({
                action: app.Actions.APP_ROUTE_CREATE,
                id: id,
                blockId: node.id,
                index: 0,
                direction: 'input',
                data: {
                    name: nodes[node.id].data.source,
                    type: 'any',
                    value: null
                },
                source: nodes[node.id].data.source
            })
            nodes[node.id].addRoute(id);
        }

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

        nodes[id].routes.forEach(function(route) {
            app.Dispatcher.dispatch({
                action: app.Actions.APP_ROUTE_DELETE,
                id: route
            })
        })

        delete nodes[id]
    }

    function updateNode(node) {
        if (nodes.hasOwnProperty(node.id) === false) {
            console.warn('could not update node: ', node.id, ' does not exist');
            return
        }
        node[node.id] = node;
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
        if (nodes[id] instanceof Source) return 'source';
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

    // after a collection of nodes has been grouped, request that their new 
    // position is relative to 0,0
    // TODO: this should and could be handled by the API. 
    function centerChildren(c) {
        var minX = +Infinity;
        var maxX = -Infinity;
        var minY = +Infinity;
        var maxY = -Infinity;

        c.forEach(function(id) {
            minX = Math.min(nodes[id].data.position.x, minX);
            minY = Math.min(nodes[id].data.position.y, minY);
            maxX = Math.max(nodes[id].data.position.x, maxX);
            maxY = Math.max(nodes[id].data.position.y, maxY);
        })

        var offX = Math.floor((minX + maxX) * .5);
        var offY = Math.floor((minY + maxY) * .5);

        c.forEach(function(id) {
            app.Utils.request(
                'PUT',
                nodeType(id) + 's/' + id + '/position', {
                    x: nodes[id].position.x - offX,
                    y: nodes[id].position.y - offY
                },
                null
            )
        });

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
            },
            function(e) {
                centerChildren(JSON.parse(e.response).children);
            }
        )
    }

    function requestNodeLabel(event) {
        app.Utils.request(
            'PUT',
            nodeType(event.id) + 's/' + event.id + '/label',
            event.label,
            null
        );
    }

    function requestSourceParams(event) {
        app.Utils.request(
            'PUT',
            nodeType(event.id) + 's/' + event.id + '/params', [{
                name: event.name,
                value: event.value
            }],
            null
        );
    }

    function requestGroupImport(event) {
        var pattern = null;
        try {
            pattern = JSON.parse(event.pattern);
        } catch (e) {
            console.warn("could not import pattern", e);
            return;
        }

        app.Utils.request(
            'POST',
            'groups/' + root + '/import',
            pattern,
            function(e) {
                var response = null;
                try {
                    response = JSON.parse(e.response);
                } catch (e) {
                    console.warn("error importing: ", response);
                }

                if (response != null) {
                    app.Dispatcher.dispatch({
                        action: app.Actions.APP_SELECT_ALL,
                        ids: response.map(function(id) {
                            return nodes.hasOwnProperty(id) ? nodes[id] : app.EdgeStore.getEdge(id);
                        }),
                    });
                }
            }
        );
    }

    function selectUnGroup() {
        var children = [];
        app.SelectionStore.getIdsByKind(Node).forEach(function(id) {
            if (nodes[id] instanceof Group) {
                children = children.concat(nodes[id].data.children);
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
            // make the ungrouped nodes center on where the parent group was
            app.Utils.request(
                'PUT',
                nodeType(children[i]) + 's/' + children[i] + '/position', {
                    x: nodes[children[i]].position.x + nodes[nodes[children[i]].parent].data.position.x,
                    y: nodes[children[i]].position.y + nodes[nodes[children[i]].parent].data.position.y
                },
                null
            )

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
            case app.Actions.WS_SOURCE_CREATE:
                createSource(event.data);
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
            case app.Actions.WS_GROUPROUTE_UPDATE:
                updateGroupRouteVisibility(event);
                nodes[event.id].emit();
                rs.emit();
                break;
            case app.Actions.APP_REQUEST_NODE_MOVE:
                finishMove();
                break;
            case app.Actions.APP_REQUEST_NODE_LABEL:
                requestNodeLabel(event);
                break;
            case app.Actions.APP_REQUEST_SOURCE_PARAMS:
                requestSourceParams(event);
                break;
            case app.Actions.APP_REQUEST_GROUP_IMPORT:
                requestGroupImport(event);
                break;
            case app.Actions.WS_BLOCK_CREATE:
                createBlock(event.data);
                rs.emit();
                break;
            case app.Actions.WS_SOURCE_DELETE:
            case app.Actions.WS_GROUP_DELETE:
            case app.Actions.WS_BLOCK_DELETE:
                deleteNode(event.id);
                rs.emit();
                break;
            case app.Actions.APP_SELECT_MOVE:
                selectMove(event.dx, event.dy);
                rs.emit();
                break;
            case app.Actions.WS_SOURCE_UPDATE:
            case app.Actions.WS_GROUP_UPDATE:
            case app.Actions.WS_BLOCK_UPDATE:
                if (!nodes.hasOwnProperty(event.id)) return;
                nodes[event.id].update(event.data);
                nodes[event.id].emit();
                rs.emit();
                break;
            case app.Actions.WS_SOURCE_UPDATE_PARAMS:
                nodes[event.id].updateParams(event.value);
                nodes[event.id].emit();
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
            case app.Actions.APP_RENDER:
                if (!nodes.hasOwnProperty(event.id)) return;
                nodes[event.id].render();
                rs.emit();
                break;
        }
    })
    app.Source = Source;
    app.Group = Group;
    app.Node = Node;
    app.NodeStore = rs;
    app.NodeSelection = selection;
}())
