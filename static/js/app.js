var app = app || {};

// TODO 
// This file desperately needs to be refactored. The portion of CoreApp that 
// is related to the stage, the background lines, and the selection of nodes
// can be put into its own component. 

(function() {
    'use strict';

    app.CoreApp = React.createClass({
        displayName: "CoreApp",
        getInitialState: function() {
            return {
                dragging: false,
                offX: null,
                offY: null,
                width: null,
                height: null,
                keys: {
                    shift: false,
                },
                selected: [],
                selectionRect: {
                    x1: null,
                    y1: null,
                    x: null,
                    y: null,
                    width: 0,
                    height: 0,
                    enabled: false
                },
                connecting: null,
                connectionFrom: null,
                connectionTo: null,
                connectionRoute: null,
                library: {
                    x: null,
                    y: null,
                    enabled: false,
                }
            }
        },
        documentMouseMove: function(e) {
            // if we don't have a focus group we need to bail
            if (this.props.model.focusedGroup === null) return;
            if (this.state.selectionRect.enabled === true) {
                var x1 = this.state.selectionRect.x1;
                var y1 = this.state.selectionRect.y1;
                var x2 = e.pageX;
                var y2 = e.pageY;
                var rectX = x2 - x1 < 0 ? x2 : x1;
                var rectY = y2 - y1 < 0 ? y2 : y1;
                var width = Math.abs(x2 - x1);
                var height = Math.abs(y2 - y1);
                var translateX = this.props.model.focusedGroup.translateX;
                var translateY = this.props.model.focusedGroup.translateY;
                var selected = [];

                // check to see which nodes are currently in the selection box
                selected = this.props.model.focusedNodes.filter(function(node) {
                    if (!node.data.hasOwnProperty('position')) return false; // we may be able to get rid of this now.
                    var position = node.data.position;
                    return app.Utils.pointInRect(rectX, rectY, width, height, position.x + translateX, position.y + translateY);
                }.bind(this));

                // check to see which edges are in selection box
                selected = selected.concat(this.props.model.focusedEdges.filter(function(node) {
                    if (!node.hasOwnProperty('from')) return false; // we may be able to get rid of this now.
                    var p1 = node.from;
                    var p2 = node.to;
                    return (app.Utils.pointInRect(rectX, rectY, width, height, p1.x + translateX, p1.y + translateY) &&
                        app.Utils.pointInRect(rectX, rectY, width, height, p2.x + translateX, p2.y + translateY));
                }.bind(this)));

                // update the state of the selection box
                this.setState({
                    selected: selected,
                    selectionRect: {
                        x1: x1,
                        y1: y1,
                        enabled: true,
                        x: rectX,
                        y: rectY,
                        width: width,
                        height: height,
                    }
                })
            } else if (this.state.dragging === true) {
                this.props.model.focusedGroup.setTranslation(e.pageX - this.state.offX, e.pageY - this.state.offY);
            }

        },
        documentKeyDown: function(e) {
            if (e.shiftKey === true) this.setState({
                keys: {
                    shift: true
                }
            })
        },
        documentKeyUp: function(e) {
            if (e.shiftKey === false) this.setState({
                keys: {
                    shift: false
                }
            })
        },
        createBlock: function(b) {
            app.Utils.request(
                'POST',
                'blocks', {
                    'type': b,
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
        componentWillMount: function() {
            document.addEventListener('keydown', this.documentKeyDown);
            document.addEventListener('keyup', this.documentKeyUp);
            document.addEventListener('mousemove', this.documentMouseMove);
            this.setState({
                width: document.body.clientWidth,
                height: document.body.clientHeight
            })
        },
        onMouseDown: function(e) {
            e.nativeEvent.button === 0 ? this.setState({
                selectionRect: {
                    x1: e.pageX,
                    y1: e.pageY,
                    enabled: true
                },
                selected: []
            }) : this.setState({
                dragging: true,
                offX: e.pageX - this.props.model.focusedGroup.translateX,
                offY: e.pageY - this.props.model.focusedGroup.translateY,
                selected: [],
            })
        },
        onMouseUp: function(e) {
            if (this.state.selectionRect.enabled === true) {
                this.setState({
                    selectionRect: {
                        enabled: false
                    }
                })
            } else if (this.state.dragging === true) {
                this.props.model.focusedGroup.setTranslation(e.pageX - this.state.offX, e.pageY - this.state.offY);
                this.setState({
                    dragging: false
                })
            }
            this.setState({
                library: {
                    enabled: false
                }
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
            } else {
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
            if (r.direction == 'input') {
                this.setState({
                    connectionTo: this.props.model.entities[r.id],
                    connectionRoute: r.route,
                })
            }

            if (r.direction == 'output') {
                this.setState({
                    connectionFrom: this.props.model.entities[r.id],
                    connectionRoute: r.route,
                })
            }

        },
        render: function() {
            var nodes = {
                'source': app.SourceComponent,
                'group': app.GroupComponent,
                'block': app.BlockComponent
            }

            var edges = {
                'link': app.LinkComponent,
                'connection': app.ConnectionComponent
            }

            var nodeElements = this.props.model.focusedNodes.map(function(c) {
                return React.createElement(app.DragContainer, {
                    key: c.data.id,
                    model: c,
                    x: c.data.position.x,
                    y: c.data.position.y,
                    nodeSelect: this.nodeSelect
                }, React.createElement(nodes[c.instance()], {
                    key: c.data.id,
                    model: c,
                    onRouteEvent: this.handleRouteEvent,
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
                    transform: 'translate(' + this.props.model.focusedGroup.translateX + ', ' + this.props.model.focusedGroup.translateY + ')',
                    key: 'renderGroups'
                }, edgeElements.concat(nodeElements));
            }

            var background = [];
            background.push(React.createElement("rect", {
                className: "background",
                x: "0",
                y: "0",
                width: this.state.width,
                height: this.state.height,
                onMouseDown: this.onMouseDown,
                onDoubleClick: this.handleDoubleClick,
                key: 'background'
            }))


            // draws panning background grid
            if (this.props.model.focusedGroup !== null) {
                var translateX = this.props.model.focusedGroup.translateX;
                var translateY = this.props.model.focusedGroup.translateY;
                var x = translateX % 200.0;
                var y = translateY % 200.0;
                var lines = [];
                var hMax = Math.floor(this.state.width / 200.0);
                var vMax = Math.floor(this.state.height / 200.0);
                for (var i = 0; i <= hMax; i++) {
                    lines.push(React.createElement('line', {
                        key: 'h' + i,
                        x1: x + (i * 200),
                        y1: 0,
                        x2: x + (i * 200),
                        y2: this.state.height,
                        stroke: 'rgba(220,220,220,1)'
                    }, null));
                }

                for (var i = 0; i <= vMax; i++) {
                    lines.push(React.createElement('line', {
                        key: 'v' + i,
                        x1: 0,
                        y1: y + (i * 200),
                        x2: this.state.width,
                        y2: y + (i * 200),
                        stroke: 'rgba(220,220,220,1)'
                    }, null));
                }

                var lineGroup = React.createElement('g', {
                    key: 'line_group',
                }, lines)

                background.push(React.createElement('circle', {
                    cx: translateX,
                    cy: translateY,
                    r: 5,
                    fill: 'rgba(220,220,220,1)',
                    key: 'origin',
                }, null))

                background.push(lines);
            }

            if (this.state.selectionRect.enabled === true) {
                background.push(React.createElement("rect", {
                    x: this.state.selectionRect.x,
                    y: this.state.selectionRect.y,
                    width: this.state.selectionRect.width,
                    height: this.state.selectionRect.height,
                    className: 'selection_rect',
                    key: 'selection_rect'
                }, null))
            }

            if (this.state.connectionFrom != null ||
                this.state.connectionTo != null) {
                background.push(React.createElement(app.ConnectionToolComponent, {
                    key: 'tool',
                    from: this.state.connectionFrom,
                    to: this.state.connectionTo,
                    route: this.state.connectionRoute
                }, null))
            }

            background.push(renderGroups);

            var stage = React.createElement("svg", {
                className: "stage",
                key: "stage",
                onMouseUp: this.onMouseUp,
                onContextMenu: function(e) {
                    e.nativeEvent.preventDefault();
                }
            }, background)

            if (this.props.model.focusedGroup !== null) {
                var groupList = React.createElement(app.GroupSelectorComponent, {
                    focusedGroup: this.props.model.focusedGroup.data.id,
                    groups: this.props.model.groups,
                    key: "group_list",
                }, null);
            }

            var panelList = React.createElement(app.PanelListComponent, {
                nodes: this.state.selected,
                key: 'panel_list',
            });

            var children = [stage, groupList, panelList];

            if (this.state.library.enabled === true) {
                children.push(React.createElement(app.AutoCompleteComponent, {
                    key: 'autocomplete',
                    x: this.state.library.x,
                    y: this.state.library.y,
                    options: this.props.model.blockLibrary,
                    onChange: this.createBlock,
                }, null));
            }


            var container = React.createElement("div", {
                className: "app",
            }, children);

            return container
        }
    })
})();
