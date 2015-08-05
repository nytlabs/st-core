var app = app || {};

// TODO 
// This file desperately needs to be refactored. The portion of CoreApp that 
// is related to the stage, the background lines, and the selection of nodes
// can be put into its own component. 

(function() {
    'use strict';

    app.CoreApp = React.createClass({
        displayName: 'CoreApp',
        getInitialState: function() {
            return {
                width: null,
                height: null,
                keys: {
                    shift: false,
                },
                selected: [],
                connecting: null,
                library: {
                    x: null,
                    y: null,
                    enabled: false,
                },
                controlKey: false,
            }
        },
        handleSelectionChange: function(rectX, rectY, width, height) {
            var selected = [];

            // check to see which nodes are currently in the selection box
            selected = this.props.model.focusedNodes.filter(function(node) {
                if (!node.data.hasOwnProperty('position')) return false; // we may be able to get rid of this now.
                var position = node.data.position;
                return app.Utils.pointInRect(rectX, rectY, width, height,
                    position.x + this.props.model.focusedGroup.translateX,
                    position.y + this.props.model.focusedGroup.translateY);
            }.bind(this));

            // check to see which edges are in selection box
            selected = selected.concat(this.props.model.focusedEdges.filter(function(node) {
                if (!node.hasOwnProperty('from')) return false; // we may be able to get rid of this now
                var p1 = node.from.node.data.position;
                var p2 = node.to.node.data.position;
                return (app.Utils.pointInRect(rectX, rectY, width, height,
                        p1.x + this.props.model.focusedGroup.translateX,
                        p1.y + this.props.model.focusedGroup.translateY) &&
                    app.Utils.pointInRect(rectX, rectY, width, height,
                        p2.x + this.props.model.focusedGroup.translateX,
                        p2.y + this.props.model.focusedGroup.translateY));
            }.bind(this)));

            // allow appending of multiple shift-selects
            if (this.state.keys.shift === true) {
                this.state.selected.forEach(function(e) {
                    if (selected.indexOf(e) === -1) {
                        selected.push(e);
                    }
                })
            }

            // update the state of the selection box
            this.setState({
                selected: selected,
            })

        },
        documentKeyDown: function(e) {
            // only fire delete if we have the stage in focus
            if (e.keyCode === 8 && e.target === document.body) {
                e.preventDefault();
                e.stopPropagation();
                this.deleteSelection();
            }

            // only fire ctrl key state if we don't have anything in focus
            if (document.activeElement === document.body && (e.keyCode === 91 || e.keyCode === 17)) {
                this.setState({
                    controlKey: true
                })
            }

            if (e.shiftKey === true) {
                this.setState({
                    keys: {
                        shift: true
                    }
                })
            }

            if (e.keyCode === 71 && e.target === document.body) this.handleGroup()
            if (e.keyCode === 85 && e.target === document.body) this.handleUngroup()
        },
        documentKeyUp: function(e) {
            if (e.keyCode === 91 || e.keyCode === 17) {
                this.setState({
                    controlKey: false
                })

            }
            if (e.shiftKey === false) this.setState({
                keys: {
                    shift: false
                }
            })
        },
        deleteSelection: function() {
            this.state.selected.forEach(function(e) {
                app.Utils.request(
                    'DELETE',
                    e.instance() + 's/' + e.data.id, {},
                    null
                )
            })
        },
        createNode: function(b) {
            app.Utils.request(
                'POST',
                b.nodeClass, {
                    'type': b.type,
                    'parent': this.props.model.focusedGroup.data.id,
                    'position': {
                        'x': this.state.library.x - this.props.model.focusedGroup.translateX,
                        'y': this.state.library.y - this.props.model.focusedGroup.translateY
                    }
                },
                function(e) {
                    if (e.status !== 200) {
                        // error
                    } else {
                        this.setState({
                            library: {
                                enabled: false
                            }
                        })
                    }
                }.bind(this))

        },
        createConnection: function(from, to) {
            if (from.source != null && from.source == to.source) {
                var source = from.direction === 'output' ? from : to;
                var block = from.direction === 'output' ? to : from;
                app.Utils.request(
                    'POST',
                    'links', {
                        'source': {
                            'id': source.id,
                        },
                        'block': {
                            'id': block.id,
                        }
                    },
                    function(e) {
                        this.setState({
                            connecting: null
                        })
                    }.bind(this)
                )
            } else {
                app.Utils.request(
                    'POST',
                    'connections', {
                        'from': {
                            'id': from.direction === 'output' ? from.id : to.id,
                            'route': from.direction === 'output' ? from.index : to.index,
                        },
                        'to': {
                            'id': to.direction === 'input' ? to.id : from.id,
                            'route': to.direction === 'input' ? to.index : from.index,
                        }
                    },
                    function(e) {
                        this.setState({
                            connecting: null
                        })
                    }.bind(this)
                )
            }
        },
        componentWillMount: function() {
            window.addEventListener('keydown', this.documentKeyDown);
            window.addEventListener('keyup', this.documentKeyUp);
            window.addEventListener('resize', this.handleResize);
            window.addEventListener('mousedown', this.handleMouseDown);
            this.handleResize();
        },
        componentWillUnmount: function() {
            window.removeEventListener('keydown', this.documentKeyDown);
            window.removeEventListener('keyup', this.documentKeyUp);
            window.removeEventListener('resize', this.handleResize);
            window.removeEventListener('mousedown', this.handleMouseDown);
        },
        handleResize: function() {
            this.setState({
                width: document.body.clientWidth,
                height: document.body.clientHeight
            })
        },
        handleMouseDown: function(e) {
            /* TODO: ensure that we are clicking on something that is not a
             * route.
             *
             * If we click anywhere other than a route, cancel the connection
             * tool. This is kind of a hack because we're only checking the
             * tag name. we don't know _for sure_ that we've clicked somewhere
             * else other than a route. eg: clicking on the origin circle won't
             * cancel the tool.
             */
            if (e.target.tagName !== 'circle') {
                this.setState({
                    connecting: null
                })
            }
        },
        handleMouseUp: function(e) {
            this.setState({
                library: {
                    enabled: false,
                }
            })
        },
        setSelection: function(ids) {
            this.setState({
                selected: ids.map(function(id) {
                    return this.props.model.entities[id];
                }.bind(this))
            })
        },
        nodeSelect: function(id) {
            var node = this.props.model.entities[id];
            if (this.state.keys.shift === true) {
                if (this.state.selected.indexOf(node) === -1) {
                    this.setState({
                        selected: this.state.selected.concat([node])
                    })
                } else {
                    this.setState({
                        selected: this.state.selected.slice().filter(function(i) {
                            return i != node;
                        })
                    })
                }
            } else if (this.state.selected.indexOf(node) === -1) {
                this.setState({
                    selected: [node],
                })
            }
        },
        handleDoubleClick: function(e) {
            this.setState({
                library: {
                    enabled: true,
                    x: e.pageX,
                    y: e.pageY
                }
            })
        },
        handleRouteEvent: function(r) {
            if (this.state.connecting === null) {
                this.setState({
                    connecting: r
                })
            }
            if (this.state.connecting !== null) {
                this.createConnection(r, this.state.connecting)
            }

        },
        handleGroup: function() {
            var bounds = {
                x1: Number.POSITIVE_INFINITY,
                x2: Number.NEGATIVE_INFINITY,
                y1: Number.POSITIVE_INFINITY,
                y2: Number.NEGATIVE_INFINITY
            }

            var ids = this.state.selected.filter(function(e) {
                return (e instanceof app.Entity)
            }).map(function(e) {
                //                console.log(e)
                if (e.data.position.x < bounds.x1) bounds.x1 = e.data.position.x;
                if (e.data.position.y < bounds.y1) bounds.y1 = e.data.position.y;
                if (e.data.position.x > bounds.x2) bounds.x2 = e.data.position.x;
                if (e.data.position.y > bounds.y2) bounds.y2 = e.data.position.y;
                return e.data.id
            });

            // prevent us from grouping nothing
            if (ids.length === 0) return;

            var position = {
                x: (bounds.x2 - bounds.x1) * .5 + bounds.x1,
                y: (bounds.y2 - bounds.y1) * .5 + bounds.y1
            }

            app.Utils.request(
                'post',
                'groups', {
                    'parent': this.props.model.focusedGroup.data.id,
                    'children': ids,
                    'position': position
                },
                function(resp) {
                    this.setState({
                        selected: [this.props.model.entities[JSON.parse(resp.response).id]]
                    })
                }.bind(this)
            )
        },
        handleDrag: function(x, y) {
            this.state.selected.filter(function(e) {
                return (e instanceof app.Entity)
            }).forEach(function(n) {
                app.Dispatcher.dispatch({
                    action: app.Actions.APP_MOVE,
                    id: n.data.id, // the selection list should simply be a list of IDs in the future.
                    dx: x,
                    dy: y
                })
                n.setPosition({
                    x: n.data.position.x + x,
                    y: n.data.position.y + y
                });
            })
        },
        handleDragStop: function() {
            this.state.selected.filter(function(e) {
                return (e instanceof app.Entity)
            }).forEach(function(n) {
                n.postPosition();
            })
        },
        handleUngroup: function() {
            var groups = this.state.selected.filter(function(e) {
                return (e.instance() === 'group') && (e.parentNode !== null)
            })

            if (groups.length === 0) return;

            var wait = {
                'jobs': 0,
                'done': 0,
            }

            var selected = [];

            groups.forEach(function(group) {
                wait.jobs += group.data.children.length
            })

            // wait for all children to be moved before deleting the group
            Object.observe(wait, function(change) {
                if (wait.jobs === wait.done) {

                    // make children of ungrouped group all selected
                    this.setState({
                        selected: selected.map(function(id) {
                            return this.props.model.entities[id]
                        }.bind(this))
                    })

                    groups.forEach(function(group) {
                        app.Utils.request(
                            'delete',
                            'groups/' + group.data.id, {},
                            null
                        )
                    })
                }
            }.bind(this))

            groups.forEach(function(group) {
                group.data.children.forEach(function(childId) {
                    selected.push(childId)
                    app.Utils.request(
                        'put',
                        'groups/' + group.parentNode.data.id + '/children/' + childId, {},
                        function() {
                            wait.done++;
                        }
                    )
                })

            }.bind(this))
        },
        render: function() {
            var nodes = {
                'source': app.BlockComponent,
                'group': app.BlockComponent,
                'block': app.BlockComponent
            }

            var edges = {
                'link': app.ConnectionComponent,
                'connection': app.ConnectionComponent
            }

            var nodeElements = this.props.model.focusedNodes.map(function(c) {
                return React.createElement(app.DragContainer, {
                    key: c.data.id,
                    model: c,
                    x: c.data.position.x,
                    y: c.data.position.y,
                    nodeSelect: this.nodeSelect,
                    onDrag: this.handleDrag,
                    onDragStop: this.handleDragStop
                }, React.createElement(nodes[c.instance()], {
                    key: c.data.id,
                    id: c.data.id,
                    //model: c,
                    //onRouteEvent: this.handleRouteEvent,
                    selected: this.state.selected.indexOf(c) !== -1 ? true : false,
                }, null))
            }.bind(this));

            var edgeElements = this.props.model.focusedEdges.map(function(c) {
                return React.createElement(edges[c.instance()], {
                    key: c.data.id,
                    model: c,
                    nodeSelect: this.nodeSelect,
                    selected: this.state.selected.indexOf(c) !== -1 ? true : false,
                }, null)
            }.bind(this));

            var renderGroups = null;
            if (this.props.model.focusedGroup !== null) {
                renderGroups = React.createElement('g', {
                    transform: 'translate(' +
                        this.props.model.focusedGroup.translateX + ', ' +
                        this.props.model.focusedGroup.translateY + ')',
                    key: 'renderGroups'
                }, edgeElements.concat(nodeElements));
            }

            var background = [];

            if (this.props.model.focusedGroup !== null) {
                background.push(React.createElement(app.StageComponent, {
                    key: 'bg',
                    group: this.props.model.focusedGroup,
                    onSelectionChange: this.handleSelectionChange,
                    onDoubleClick: this.handleDoubleClick,
                    onMouseUp: this.handleMouseUp,
                    width: this.state.width,
                    height: this.state.height,
                }, null));
            }

            if (this.props.model.focusedGroup !== null) {
                var groupList = React.createElement(app.GroupSelectorComponent, {
                    focusedGroup: this.props.model.focusedGroup.data.id,
                    groups: this.props.model.groups,
                    key: 'group_list',
                }, null);

                if (this.state.connecting != null) {
                    background.push(React.createElement(app.ConnectionToolComponent, {
                        key: 'tool',
                        connecting: this.state.connecting,
                        node: this.props.model.entities[this.state.connecting.id],
                        translateX: this.props.model.focusedGroup.translateX,
                        translateY: this.props.model.focusedGroup.translateY
                    }, null))
                }
            }

            background.push(renderGroups);

            var stage = React.createElement('svg', {
                className: 'stage',
                key: 'stage',
                width: this.state.width,
                height: this.state.height,
                onContextMenu: function(e) {
                    e.nativeEvent.preventDefault();
                }
            }, background)

            var tools = React.createElement(app.ToolsComponent, {
                key: 'tool_list',
                onGroup: this.handleGroup,
                onUngroup: this.handleUngroup
            });

            var panelList = React.createElement(app.PanelListComponent, {
                nodes: this.state.selected,
                key: 'panel_list',
            });


            var children = [stage, groupList, panelList, tools];

            if (this.props.model.focusedGroup !== null) {
                var clipboard = React.createElement(app.ClipboardComponent, {
                    selected: this.props.model.recurseSelection(this.state.selected),
                    focus: this.state.controlKey,
                    key: 'clipboard',
                    group: this.props.model.focusedGroup.data.id,
                    setSelection: this.setSelection,
                });
                children.push(clipboard);
            }

            if (this.state.library.enabled === true) {
                children.push(React.createElement(app.AutoCompleteComponent, {
                    key: 'autocomplete',
                    x: this.state.library.x,
                    y: this.state.library.y,
                    options: this.props.model.library,
                    onChange: this.createNode,
                }, null));
            }


            var container = React.createElement('div', {
                className: 'app',
            }, children);

            return container
        }
    })
})();
