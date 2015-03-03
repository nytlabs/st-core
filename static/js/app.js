var app = app || {};

(function() {
    app.CoreApp = React.createClass({
        displayName: "CoreApp",
        getInitialState: function() {
            return {
                dragging: false,
                x: 0,
                y: 0,
                offX: null,
                offY: null,
                width: null,
                height: null,
                keys: {
                    shift: false,
                },
                selected: [],
                group: 0,
                selectionRect: {
                    x1: null,
                    y1: null,
                    x: null,
                    y: null,
                    width: 0,
                    height: 0,
                    enabled: false
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

                var selected = [];

                selected = this.props.model.entities[this.props.model.focusedGroup].data.children.filter(function(id) {
                    var node = this.props.model.entities[id].data;
                    return node.hasOwnProperty('position') &&
                        node.position.x + this.state.x >= rectX &&
                        node.position.x + this.state.x < rectX + width &&
                        node.position.y + this.state.y >= rectY &&
                        node.position.y + this.state.y < rectY + height
                }.bind(this));

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
                this.setState({
                    x: e.pageX - this.state.offX,
                    y: e.pageY - this.state.offY
                })
            }

            this.setState({
                width: document.body.clientWidth,
                height: document.body.clientHeight
            })
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
                offX: e.pageX - this.state.x,
                offY: e.pageY - this.state.y,
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
                this.setState({
                    x: e.pageX - this.state.offX,
                    y: e.pageY - this.state.offY,
                    dragging: false
                })
            }
        },
        nodeSelect: function(id) {
            if (this.state.keys.shift === true) {
                if (this.state.selected.indexOf(id) === -1) {
                    this.setState({
                        selected: this.state.selected.concat([id])
                    })
                } else {
                    this.setState({
                        selected: this.state.selected.slice().filter(function(i) {
                            return i != id;
                        })
                    })
                }
            } else {
                this.setState({
                    selected: [id],
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
                    selected: this.state.selected.indexOf(c.data.id) !== -1 ? true : false,
                }, null))
            }.bind(this));

            var edgeElements = this.props.model.focusedEdges.map(function(c) {
                return React.createElement(edges[c.instance()], {
                    key: c.data.id,
                    model: c,
                    graph: this.props.model
                }, null)
            }.bind(this));

            var renderGroups = React.createElement('g', {
                transform: 'translate(' + this.state.x + ', ' + this.state.y + ')',
                key: 'renderGroups'
            }, edgeElements.concat(nodeElements));

            var background = [];
            background.push(React.createElement("rect", {
                className: "background",
                x: "0",
                y: "0",
                width: this.state.width,
                height: this.state.height,
                onMouseDown: this.onMouseDown,
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
            }, null)

            var container = React.createElement("div", {
                className: "app",
            }, [stage, groupList]);

            return container
        }
    })
})();
