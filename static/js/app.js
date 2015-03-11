var app = app || {};

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
                library: {
                    x: null,
                    y: null,
                    enabled: false,
                }
            }
        },
        documentMouseMove: function(e) {
            if (this.state.selectionRect.enabled === true) {
                var x1 = this.state.selectionRect.x1;
                var y1 = this.state.selectionRect.y1;
                var x2 = e.pageX;
                var y2 = e.pageY;
                var rectX = x2 - x1 < 0 ? x2 : x1;
                var rectY = y2 - y1 < 0 ? y2 : y1;
                var width = Math.abs(x2 - x1);
                var height = Math.abs(y2 - y1);
                var translateX = this.props.model.entities[this.props.model.focusedGroup].translateX;
                var translateY = this.props.model.entities[this.props.model.focusedGroup].translateY;
                var selected = [];

                selected = this.props.model.focusedNodes.filter(function(node) {
                    if (!node.data.hasOwnProperty('position')) return false; // we may be able to get rid of this now.
                    var position = node.data.position;
                    return app.Utils.pointInRect(rectX, rectY, width, height, position.x + translateX, position.y + translateY);
                }.bind(this));

                selected = selected.concat(this.props.model.focusedEdges.filter(function(node) {
                    if (!node.hasOwnProperty('from')) return false; // we may be able to get rid of this now.
                    var p1 = node.from;
                    var p2 = node.to;
                    return (app.Utils.pointInRect(rectX, rectY, width, height, p1.x + translateX, p1.y + translateY) &&
                        app.Utils.pointInRect(rectX, rectY, width, height, p2.x + translateX, p2.y + translateY));
                }.bind(this)));

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
                this.props.model.entities[this.props.model.focusedGroup].setTranslation(e.pageX - this.state.offX, e.pageY - this.state.offY);
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
                offX: e.pageX - this.props.model.entities[this.props.model.focusedGroup].translateX,
                offY: e.pageY - this.props.model.entities[this.props.model.focusedGroup].translateY,
                selected: [],
            })
            this.setState({
                library: {
                    enabled: false
                }
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
                this.props.model.entities[this.props.model.focusedGroup].setTranslation(e.pageX - this.state.offX, e.pageY - this.state.offY);
                this.setState({
                    dragging: false
                })
            }
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
            if (this.props.model.entities.hasOwnProperty(this.props.model.focusedGroup) === true) {
                renderGroups = React.createElement('g', {
                    transform: 'translate(' + this.props.model.entities[this.props.model.focusedGroup].translateX + ', ' + this.props.model.entities[this.props.model.focusedGroup].translateY + ')',
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

            background.push(renderGroups);

            var stage = React.createElement("svg", {
                className: "stage",
                key: "stage",
                onMouseUp: this.onMouseUp,
                onContextMenu: function(e) {
                    e.nativeEvent.preventDefault();
                }
            }, background)

            var groupList = React.createElement(app.GroupSelectorComponent, {
                focusedGroup: this.props.model.focusedGroup,
                groups: this.props.model.groups,
                key: "group_list",
            }, null)

            var panelList = React.createElement(app.PanelListComponent, {
                nodes: this.state.selected,
                key: 'panel_list',
            });

            var children = [stage, groupList, panelList];



            var container = React.createElement("div", {
                className: "app",
            }, children);

            return container
        }
    })
})();
