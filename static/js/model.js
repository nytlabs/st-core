var app = app || {};

// TODO:
// create a standard model API that the rest of the components can use
// this standard API should use WS to communicate back to server.
// 
// the inform() situation is a mess. we should limit the number of places that
// it makes an appearance is, as it causes horrendous snake calls that refresh
// the UI too many times.

(function() {
    'use strict';

    var dm = new app.Utils.DebounceManager();

    app.CoreModel = function() {
        this.entities = {};
        this.list = [];
        this.groups = [];
        this.edges = [];
        this.blockLibrary = [];
        this.sourceLibrary = [];

        this.onChanges = [];

        this.focusedGroup = null;
        this.focusedNodes = []; // nodes apart of the focused group
        this.focusedEdges = []; // nedges that are apart of the focused group

        var ws = new WebSocket("ws://localhost:7071/updates");

        ws.onmessage = function(m) {
            this.update(JSON.parse(m.data));
        }.bind(this)

        ws.onopen = function() {
            ws.send('list');
        }

        app.Utils.request(
            "GET",
            "blocks/library",
            null,
            function(req) {
                this.blockLibrary = JSON.parse(req.response);
            }.bind(this)
        )

        app.Utils.request(
            "GET",
            "sources/library",
            null,
            function(req) {
                this.sourceLibrary = JSON.parse(req.response);
            }.bind(this)
        )
    }

    app.CoreModel.prototype.subscribe = function(onChange) {
        this.onChanges.push(onChange);
    }

    app.CoreModel.prototype.inform = function() {
        this.onChanges.forEach(function(cb) {
            cb();
        });
    }

    app.Entity = function() {
        this.isDragging = false;
    }

    app.Entity.prototype.setPosition = function(p) {
        this.data.position.x = p.x;
        this.data.position.y = p.y;

        // this function refreshes all connection geometry in view
        // it may be better to have a specific call for just connections that
        // are touching this particular entity.
        //this.model.refreshFocusedEdgeGeometry();
        this.model.inform()
        dm.push(this.id, function() {
            app.Utils.request(
                "PUT",
                this.instance() + "s/" + this.data.id + "/position", // would be nice to change API to not have the "S" in it!
                p,
                null
            );
        }.bind(this), 50)
    }

    app.Entity.prototype.setDragging = function(e) {
        this.isDragging = e;
    }


    app.Block = function(data, model) {
        app.Entity.call(this);

        this.routes = [];
        this.geometry = [];
        this.data = data;
        this.model = model;

        this.buildRoutes();
        this.buildGeometry();
    }

    app.Block.prototype = new app.Entity();

    app.Block.prototype.instance = function() {
        return "block";
    }


    app.Block.prototype.buildRoutes = function() {
        this.routes = this.data.inputs.map(function(input, index) {
            return {
                'direction': 'input',
                'index': index,
                'displayIndex': index
            };
        })

        this.routes = this.routes.concat(this.data.outputs.map(function(output, index) {
            return {
                'direction': 'output',
                'index': index,
                'displayIndex': index
            }
        }))

        this.routes = this.routes.map(function(r, index) {
            r.id = this.data.id;
            r.connections = [];
            r.data = this.data[r.direction + 's'][r.index];
            r.routesIndex = index;
            return r
        }.bind(this));

    }

    app.Block.prototype.buildGeometry = function() {
        var textMeasures,
            maxWidth = {
                input: 0,
                output: 0
            },
            routeHeight = 0,
            routeRadius = 5,
            routeGeometry = [],
            num = {
                input: 0,
                output: 0
            };

        textMeasures = this.routes.map(function(r) {
            var measure = app.Utils.measureText(r.data.name, 'route_label');

            if (measure.width > maxWidth[r.direction]) {
                maxWidth[r.direction] = measure.width;
            }
            if (measure.height > routeHeight) {
                routeHeight = measure.height;
            }
            num[r.direction]++;

            return measure
        });

        this.geometry = {
            'width': maxWidth.input + maxWidth.output,
            'height': Math.max(num.input, num.output) * routeHeight,
            'routeHeight': routeHeight,
            'routeRadius': routeRadius,
        }
    }

    /* app.Block.prototype.refresh = function() {

         this.routes = this.data.inputs.map(function(i){
             return {"geometry": app.Utils.measureText(i.name, "route_label");
         }.bind(this))
         
         this.inputs = this.data.inputs.map(function(i) {
             return app.Utils.measureText(i.name, "route_label");
         }.bind(this))
         this.outputs = this.data.outputs.map(function(o) {
             return app.Utils.measureText(o.name, "route_label");
         }.bind(this));

         var inputMaxWidth = [{
             width: 0
         }].concat(this.inputs).reduce(function(p, v) {
             return (p.width > v.width ? p : v);
         })

         var outputMaxWidth = [{
             width: 0
         }].concat(this.outputs).reduce(function(p, v) {
             return (p.width > v.width ? p : v);
         })

         var inputMaxHeight = [{
             height: 0
         }].concat(this.inputs).reduce(function(p, v) {
             return (p.height > v.height ? p : v);
         })

         var outputMaxHeight = [{
             height: 0
         }].concat(this.outputs).reduce(function(p, v) {
             return (p.height > v.height ? p : v);
         })

         this.width = inputMaxWidth.width + outputMaxWidth.width;
         this.routeHeight = Math.max(inputMaxHeight.height, outputMaxHeight.height);
         this.height = Math.max(this.inputs.length, this.outputs.length) * this.routeHeight;

         this.routeRadius = 5;

         this.inputs.forEach(function(e, i) {
             e.routeY = (i + 1) * this.routeHeight;
             e.routeX = 0;
             e.routeCircleX = -this.routeRadius * .5;
             e.routeCircleY = -this.routeHeight * .5;
             e.routeAlign = "start"; // this should be deleted & derived from routeDirection
             e.routeIndex = i;
             e.routeDirection = 'input';
             e.routeRadius = this.routeRadius;
             e.data = this.data.inputs[i];
             e.connections = [];
             e.blockId = this.data.id;
         }.bind(this));

         this.outputs.forEach(function(e, i) {
             e.routeY = (i + 1) * this.routeHeight;
             e.routeX = this.width;
             e.routeCircleX = this.routeRadius * .5;
             e.routeCircleY = -this.routeHeight * .5;
             e.routeDirection = 'output';
             e.routeIndex = i;
             e.routeAlign = "end";
             e.routeRadius = this.routeRadius;
             e.data = this.data.outputs[i];
             e.connections = [];
             e.blockId = this.data.id;
         }.bind(this));
     }*/

    /*app.Group = function(data, model, children) {
        app.Block.call(this, data, model)
        this.children = children
    }*/

    app.Group = function(data, model) {
        app.Entity.call(this);
        //        app.Block.call(this);
        this.data = data;
        this.model = model;

        // translation coords for each group workspace.
        // not synced with server.
        this.translateX = 0;
        this.translateY = 0;
        this.refresh();
    }

    app.Group.prototype = new app.Entity();

    app.Group.prototype.instance = function() {
        return "group";
    }

    app.Group.prototype.setTranslation = function(x, y) {
        this.translateX = x;
        this.translateY = y;
        this.model.inform();
    }

    app.Group.prototype.refresh = function() {
        console.log("iiii neeed to refresssshhhhhhhhh")
        this.data.outputs = [];
        this.data.inputs = [];
        this.inputs = [];
        this.outputs = [];
        /* for (var i in this.data.children) {
             this.data.outputs = this.data.outputs.concat(this.model.entities[this.data.children[i]].data.outputs)
             this.data.inputs = this.data.inputs.concat(this.model.entities[this.data.children[i]].data.inputs)
         }

         this.data.type = 'group'

         this.inputs = this.data.inputs.map(function(i) {
             return app.Utils.measureText(i.name, "route_label");
         }.bind(this))
         this.outputs = this.data.outputs.map(function(o) {
             return app.Utils.measureText(o.name, "route_label");
         }.bind(this));

         var inputMaxWidth = [{
             width: 0
         }].concat(this.inputs).reduce(function(p, v) {
             return (p.width > v.width ? p : v);
         })

         var outputMaxWidth = [{
             width: 0
         }].concat(this.outputs).reduce(function(p, v) {
             return (p.width > v.width ? p : v);
         })

         var inputMaxHeight = [{
             height: 0
         }].concat(this.inputs).reduce(function(p, v) {
             return (p.height > v.height ? p : v);
         })

         var outputMaxHeight = [{
             height: 0
         }].concat(this.outputs).reduce(function(p, v) {
             return (p.height > v.height ? p : v);
         })

         this.width = inputMaxWidth.width + outputMaxWidth.width;
         this.routeHeight = Math.max(inputMaxHeight.height, outputMaxHeight.height);
         this.height = Math.max(this.inputs.length, this.outputs.length) * this.routeHeight;

         this.routeRadius = 5;

         this.inputs.forEach(function(e, i) {
             e.routeY = (i + 1) * this.routeHeight;
             e.routeX = 0;
             e.routeCircleX = -this.routeRadius * .5;
             e.routeCircleY = -this.routeHeight * .5;
             e.routeAlign = "start"; // this should be deleted & derived from routeDirection
             e.routeIndex = i;
             e.routeDirection = 'input';
             e.routeRadius = this.routeRadius;
             e.data = this.data.inputs[i];
             e.connections = this.data.inputs[i].connections;
             console.log(this.data.inputs[i])
         }.bind(this));

         this.outputs.forEach(function(e, i) {
             e.routeY = (i + 1) * this.routeHeight;
             e.routeX = this.width;
             e.routeCircleX = this.routeRadius * .5;
             e.routeCircleY = -this.routeHeight * .5;
             e.routeDirection = 'output';
             e.routeIndex = i;
             e.routeAlign = "end";
             e.routeRadius = this.routeRadius;
             e.data = this.data.outputs[i];
             e.connections = this.data.outputs[i].connections;
         }.bind(this));*/
    }

    // when a group changes, this swaps out references in focusedNodes and focusedEdges
    app.Group.prototype.refreshFocusedGroup = function() {
        var model = this.model;
        var id = this.data.id;

        model.focusedNodes = model.entities[id].data.children.map(function(id) {
            return this.entities[id];
        }.bind(model))

        model.focusedEdges = model.edges.filter(function(e) {
            switch (e.instance()) {
                case 'connection':
                    if (this.entities[id].data.children.indexOf(e.data.to.id) !== -1) {
                        return true;
                    }
                    break;
                case 'link':
                    if (this.entities[id].data.children.indexOf(e.data.block.id) !== -1) {
                        return true;
                    }
                    break;
            }
            return false;
        }.bind(model));

        //model.refreshFocusedEdgeGeometry();
        model.inform()
    }

    // setFocusedGroup sets takes a group id and prepares that group to be 
    // viewed. It changes the model's current group in focus, in addition to
    // preparing focusedNodes and focusedEdges.
    app.Group.prototype.setFocusedGroup = function() {
        this.model.focusedGroup = this;
        this.refreshFocusedGroup();
        this.model.inform();
    }

    /*app.Block = function(data, model) {
        app.Entity.call(this);
        this.data = data;
        this.model = model;
        this.refreshGeometry();
    }

    app.Block.prototype = new app.Entity();

    app.Block.prototype.instance = function() {
        return "block";
    }

    app.Block.prototype.refreshGeometry = function() {
        this.inputs = this.data.inputs.map(function(i) {
            return app.Utils.measureText(i.name, "route_label");
        }.bind(this))
        this.outputs = this.data.outputs.map(function(o) {
            return app.Utils.measureText(o.name, "route_label");
        }.bind(this));

        var inputMaxWidth = [{
            width: 0
        }].concat(this.inputs).reduce(function(p, v) {
            return (p.width > v.width ? p : v);
        })

        var outputMaxWidth = [{
            width: 0
        }].concat(this.outputs).reduce(function(p, v) {
            return (p.width > v.width ? p : v);
        })

        var inputMaxHeight = [{
            height: 0
        }].concat(this.inputs).reduce(function(p, v) {
            return (p.height > v.height ? p : v);
        })

        var outputMaxHeight = [{
            height: 0
        }].concat(this.outputs).reduce(function(p, v) {
            return (p.height > v.height ? p : v);
        })

        this.width = inputMaxWidth.width + outputMaxWidth.width;
        this.routeHeight = Math.max(inputMaxHeight.height, outputMaxHeight.height);
        this.height = Math.max(this.inputs.length, this.outputs.length) * this.routeHeight;

        this.routeRadius = 5;

        this.inputs.forEach(function(e, i) {
            e.routeY = (i + 1) * this.routeHeight;
            e.routeX = 0;
            e.routeCircleX = -this.routeRadius * .5;
            e.routeCircleY = -this.routeHeight * .5;
            e.routeAlign = "start"; // this should be deleted & derived from routeDirection
            e.routeIndex = i;
            e.routeDirection = 'input';
            e.routeRadius = this.routeRadius;
            e.data = this.data.inputs[i];
            e.connections = [];
            e.blockId = this.data.id;
        }.bind(this));

        this.outputs.forEach(function(e, i) {
            e.routeY = (i + 1) * this.routeHeight;
            e.routeX = this.width;
            e.routeCircleX = this.routeRadius * .5;
            e.routeCircleY = -this.routeHeight * .5;
            e.routeDirection = 'output';
            e.routeIndex = i;
            e.routeAlign = "end";
            e.routeRadius = this.routeRadius;
            e.data = this.data.outputs[i];
            e.connections = [];
            e.blockId = this.data.id;
        }.bind(this));
    }*/

    app.Source = function(data, model) {
        app.Entity.call(this);
        this.data = data;
        this.model = model;
    }

    app.Source.prototype = new app.Entity();

    app.Source.prototype.instance = function() {
        return "source";
    }

    app.Connection = function(data, model) {
        this.data = data;
        this.model = model;
        this.from = {
            node: model.entities[data.from.id],
            route: model.entities[data.from.id].routes.filter(function(r) {
                return (r.index === data.from.route) && (r.direction === 'output');
            })[0]
        }

        this.to = {
            node: model.entities[data.to.id],
            route: model.entities[data.to.id].routes.filter(function(r) {
                return (r.index === data.to.route) && (r.direction === 'input');
            })[0]
        }

        this.attach();
    }


    // attach() and detach() adds/removes a reference to this connection the route on the block entity.
    app.Connection.prototype.attach = function() {
        this.model.entities[this.data.from.id].routes.filter(function(r) {
            return (r.index === this.data.from.route) && (r.direction === 'output')
        }.bind(this))[0].connections.push(this);

        this.model.entities[this.data.to.id].routes.filter(function(r) {
            return (r.index === this.data.to.route) && (r.direction === 'input')
        }.bind(this))[0].connections.push(this);
    }

    app.Connection.prototype.detach = function() {
        var fromConnections = this.model.entities[this.data.from.id].routes.filter(function(r) {
            return (r.index === this.data.from.route) && (r.direction === 'output')
        }.bind(this))[0].connections;

        var toConnections = this.model.entities[this.data.to.id].routes.filter(function(r) {
            return (r.index === this.data.to.route) && (r.direction === 'input')
        }.bind(this))[0].connections;

        fromConnections.splice(fromConnections.indexOf(this), 1);
        toConnections.splice(toConnections.indexOf(this), 1);
    }

    //app.Connection.prototype = new app.Entity();

    app.Connection.prototype.instance = function() {
        return "connection";
    }

    /*    app.Connection.prototype.refreshGeometry = function() {
            var from = this.model.entities[this.data.from.id];
            var to = this.model.entities[this.data.to.id];

            var x1 = from.data.position.x + from.outputs[this.data.from.route].routeX + from.outputs[this.data.from.route].routeCircleX;
            var y1 = from.data.position.y + from.outputs[this.data.from.route].routeY + from.outputs[this.data.from.route].routeCircleY;
            var cx1 = x1 + 50.0;
            var cy1 = y1;
            var x2 = to.data.position.x + to.inputs[this.data.to.route].routeX + to.inputs[this.data.to.route].routeCircleX;
            var y2 = to.data.position.y + to.inputs[this.data.to.route].routeY + to.inputs[this.data.to.route].routeCircleY;
            var cx2 = x2 - 50.0;
            var cy2 = y2;

            this.from = {
                x: x1,
                y: y1
            };

            this.to = {
                x: x2,
                y: y2
            };

            this.routeRadius = 3;

            this.path = ['M', x1, ' ', y1, ' C ', cx1, ' ', cy1, ' ', cx2, ' ', cy2, ' ', x2, ' ', y2].join('');
        }*/

    app.Link = function(data, model) {
        //app.Entity.call(this);
        this.data = data;
        this.model = model;
    }

    //app.Link.prototype = new app.Entity();

    app.Link.prototype.instance = function() {
        return "link";
    }

    app.Link.prototype.refreshGeometry = function() {
        //TODO
    }

    var nodes = {
        'block': app.Block,
        'source': app.Source,
        'group': app.Group,
        'connection': app.Connection,
        'link': app.Link
    }

    // this takes an id and puts it at the very top of the list
    app.CoreModel.prototype.select = function(id) {
        this.list.push(this.list.splice(this.list.indexOf(this.entities[id]), 1)[0]);
        this.inform();
    }

    app.CoreModel.prototype.addChild = function(groupId, id) {
        this.entities[groupId].data.children.push(id);
        this.entities[groupId].refresh();

        if (groupId === this.focusedGroup.data.id) this.entities[groupId].refreshFocusedGroup();
        this.inform();
    }

    app.CoreModel.prototype.removeChild = function(groupId, id) {
        this.entities[groupId].data.children.splice(this.entities[groupId].data.children.indexOf(id), 1);
        this.entities[groupId].refresh();

        if (groupId === this.focusedGroup.data.id) this.entities[groupId].refreshFocusedGroup();
        this.inform();
    }

    /*app.CoreModel.prototype.refreshFocusedEdgeGeometry = function() {
        this.focusedEdges.forEach(function(e) {
            e.refreshGeometry();
        })
        this.inform();
    }*/

    app.CoreModel.prototype.update = function(m) {
        switch (m.action) {
            case 'update':
                if (m.type === 'block' ||
                    m.type === 'group' ||
                    m.type === 'source') {
                    for (var key in m.data[m.type]) {
                        if (key !== 'id') {
                            this.entities[m.data[m.type].id].data[key] = m.data[m.type][key];

                            // TODO: sort out model updates
                            // this stops the feedback loop from a client making a request
                            // that ends up updating the client for the same change to the model.
                            // TWO separate models that represent the same thing is 
                            // an anti-pattern, HOWEVER in this circumstance we are doing this on 
                            // purpose -- we want the client to have immediate feedback from dragging
                            // a node, and we want to broadcast this to the rest of the clients
                            // at a throttle rate. This means we have to create a way to reconcile
                            // the messages coming from the server with the client side node that
                            // is being dragged.
                            //
                            // The following updates all node geometry for nodes that are NOT 
                            // currently being dragged in this client state. 
                            //if (!this.entities[m.data[m.type].id].isDragging) {
                            //   console.log(m)
                            //this.refreshFocusedEdgeGeometry();
                            //   return;
                            // }
                        }
                    }
                } else if (m.type === 'route') {
                    this.entities[m.data.id].data.inputs[m.data.route].value = m.data.value
                }
                break;
            case 'create':
                // create seperate action for child.
                if (m.type === "child") {
                    this.addChild(m.data.group.id, m.data.child.id);
                    return;
                }

                // we put a reference to model in each entitiy so that we can
                // propagate inform();
                var n = new nodes[m.type](m.data[m.type], this);
                this.entities[m.data[m.type].id] = n;
                this.list.push(this.entities[m.data[m.type].id]);

                // if we currently don't have a focus group, we wait for the 
                // first group to arrive from the server. We then set that
                // group as our currently focused group.
                if (m.type === "group") {
                    this.groups.push(n);
                    if (this.focusedGroup === null) {
                        n.setFocusedGroup();
                        return;
                    }
                }


                if (m.type === "connection" || m.type === "link") {
                    this.edges.push(n);
                }

                // if we have a focused group we need to have a way to update the 
                // conections that are currently on display. 
                if (this.focusedGroup != null) {
                    this.focusedGroup.refreshFocusedGroup();
                    return;
                }

                break;
            case 'delete':
                if (m.type === "child") {
                    this.removeChild(m.data.group.id, m.data.child.id); // this child nonsense is a mess
                    return
                }

                if (m.type === "connection") {
                    this.entities[m.data[m.type].id].detach();
                }

                var i = this.list.indexOf(this.entities[m.data[m.type].id]);
                this.list.splice(i, 1);

                if (m.type === "group") {
                    var i = this.groups.indexOf(this.entities[m.data[m.type].id]);
                    this.groups.splice(i, 1);
                }

                if (m.type === "connection" || m.type == "link") {
                    var i = this.edges.indexOf(this.entities[m.data[m.type].id]);
                    this.edges.splice(i, 1);
                }

                delete this.entities[m.data[m.type].id];

                if (this.focusedGroup != null) {
                    this.focusedGroup.refreshFocusedGroup();
                }
                break;
        }

        this.inform();
    }
})();
