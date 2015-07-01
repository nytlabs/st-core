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


    // when a group changes, this swaps out references in focusedNodes and focusedEdges
    function refreshFocusedGroup(g, model) {
        //        var model = this.model;
        //        var id = this.data.id;
        var id = g.data.id;
        model.focusedNodes = model.entities[id].data.children.map(function(id) {
            return this.entities[id];
        }.bind(model))

        model.focusedEdges = model.edges.filter(function(e) {
            switch (e.instance()) {
                case 'link':
                    // TODO: possibly combine 'link' and 'connection' cases
                    // into one. The only difference is that links don't care
                    // about route indices. 

                    var toRoute, fromRoute;

                    var to = model.focusedNodes.filter(function(n) {
                        return !!(n.routes.filter(function(r) {
                            if (r.id === e.data.source.id) {
                                toRoute = r
                                return true
                            }
                            return false
                        }).length)
                    })

                    var from = model.focusedNodes.filter(function(n) {
                        return !!(n.routes.filter(function(r) {
                            if (r.id === e.data.block.id) {
                                fromRoute = r
                                return true
                            }
                            return false
                        }).length)
                    })

                    if (!!to.length && !!from.length && to[0] !== from[0]) {
                        e.setNodes(to[0], toRoute, from[0], fromRoute);
                        return true
                    }


                    break;
                case 'connection':
                    var toRoute, fromRoute;

                    var to = model.focusedNodes.filter(function(n) {
                        return !!(n.routes.filter(function(r) {
                            if ((r.id === e.data.to.id) && (r.index === e.data.to.route) && r.direction === 'input') {
                                toRoute = r
                                return true
                            }
                            return false
                        }).length)
                    })

                    var from = model.focusedNodes.filter(function(n) {
                        return !!(n.routes.filter(function(r) {
                            if ((r.id === e.data.from.id) && (r.index === e.data.from.route) && r.direction === 'output') {
                                fromRoute = r
                                return true
                            }
                            return false
                        }).length)
                    })

                    if (!!to.length && !!from.length && to[0] !== from[0]) {
                        e.setNodes(from[0], fromRoute, to[0], toRoute);
                        return true
                    }

                    break;
            }
            return false;
        }.bind(model));

        //model.refreshFocusedEdgeGeometry();
        model.inform()
    }

    app.CoreModel = function() {
        this.entities = {};
        this.list = [];
        this.groups = [];
        this.edges = [];
        this.library = [];

        this.onChanges = [];

        this.focusedGroup = null;
        this.focusedNodes = []; // nodes apart of the focused group
        this.focusedEdges = []; // nedges that are apart of the focused group

        var ws = new WebSocket('ws://localhost:7071/updates');

        ws.onmessage = function(m) {
            this.update(JSON.parse(m.data));
        }.bind(this)

        ws.onopen = function() {
            ws.send('list');
        }

        app.Utils.request(
            'GET',
            'blocks/library',
            null,
            function(req) {
                this.library = this.library.concat(JSON.parse(req.response).map(function(entry) {
                    return {
                        type: entry.type,
                        source: entry.source,
                        nodeClass: 'blocks'
                    }
                }));
            }.bind(this)
        )

        app.Utils.request(
            'GET',
            'sources/library',
            null,
            function(req) {
                this.library = this.library.concat(JSON.parse(req.response).map(function(entry) {
                    return {
                        type: entry.type,
                        source: entry.source,
                        nodeClass: 'sources'
                    }
                }));
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

    // TODO: make addchild set parentNode and propagate 'refresh' upwards to parentNode.
    app.CoreModel.prototype.addChild = function(groupId, id) {
        this.entities[id].parentNode = this.entities[groupId]
        this.entities[groupId].data.children.push(id);
        this.entities[groupId].refresh();

        // if (this.entities[groupId].parentNode !== null) {
        // if this group has a parent, refresh that parent.
        //     this.entities[groupId].refresh();
        // }

        refreshFocusedGroup(this.focusedGroup, this);
        this.inform();
    }

    app.CoreModel.prototype.removeChild = function(groupId, id) {
        this.entities[groupId].data.children.splice(this.entities[groupId].data.children.indexOf(id), 1);
        this.entities[groupId].refresh();


        //if (this.entities[groupId].parentNode !== null) {
        // if this group has a parent, refresh that parent.
        //    this.entities[groupId].refresh();
        // }

        refreshFocusedGroup(this.focusedGroup, this);
        this.inform();
    }

    app.CoreModel.prototype.setFocusedGroup = function(group) {
        this.focusedGroup = group
        refreshFocusedGroup(this.focusedGroup, this);
        this.inform();
    }

    app.CoreModel.prototype.update = function(m) {
        switch (m.action) {
            case 'update':

                if (m.type === 'block' ||
                    m.type === 'group' ||
                    m.type === 'source') {
                    for (var key in m.data[m.type]) {
                        if (key !== 'id') {
                            this.entities[m.data[m.type].id].data[key] = m.data[m.type][key];
                        }
                    }
                } else if (m.type === 'route') {
                    this.entities[m.data.id].data.inputs[m.data.route].value = m.data.value
                } else if (m.type === 'param') {
                    var key = {} 
                    this.entities[m.data.id].data.params.map(function(kv, index, array){
                      key[kv.name]=index
                    })
                    this.entities[m.data.id].data.params[key[m.data.param]].value = m.data.value
                }
                break;
            case 'create':

                // create seperate action for child.
                if (m.type === 'child') {
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
                if (m.type === 'group') {
                    this.groups.push(n);
                    if (this.focusedGroup === null) {
                        n.setFocusedGroup();
                        refreshFocusedGroup(this.focusedGroup, this);
                        return;
                    }
                }


                if (m.type === 'connection' || m.type === 'link') {
                    this.edges.push(n);
                }

                // for source we need to populate its parameters
                if (m.type === 'source'){
                  this.entities[m.data.source.id].data.params = m.data.source.params
                }

                // if we have a focused group we need to have a way to update the 
                // conections that are currently on display. 
                if (this.focusedGroup != null) {
                    refreshFocusedGroup(this.focusedGroup, this);
                    //this.focusedGroup.refreshFocusedGroup();
                    return;
                }


                break;
            case 'delete':
                if (m.type === 'child') {
                    this.removeChild(m.data.group.id, m.data.child.id); // this child nonsense is a mess
                    return
                }

                if (m.type === 'connection') {
                    this.entities[m.data[m.type].id].detach();
                }

                var i = this.list.indexOf(this.entities[m.data[m.type].id]);
                this.list.splice(i, 1);

                if (m.type === 'group') {
                    var i = this.groups.indexOf(this.entities[m.data[m.type].id]);
                    this.groups.splice(i, 1);
                }

                if (m.type === 'connection' || m.type == 'link') {
                    var i = this.edges.indexOf(this.entities[m.data[m.type].id]);
                    this.edges.splice(i, 1);
                }

                delete this.entities[m.data[m.type].id];

                if (this.focusedGroup != null) {
                    refreshFocusedGroup(this.focusedGroup, this);
                    //this.focusedGroup.refreshFocusedGroup();
                }
                break;
        }

        this.inform();
    }
})();
