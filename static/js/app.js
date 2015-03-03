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
                var width = e.pageX - this.state.selectionRect.x;
                var height = e.pageY - this.state.selectionRect.y;
                var selected = [];

                selected = this.props.model.entities[this.props.model.focusedGroup].children.filter(function(id) {
                    var node = this.props.model.entities[id];
                    return node.hasOwnProperty('position') &&
                        node.position.x + this.state.x >= this.state.selectionRect.x &&
                        node.position.x + this.state.x < this.state.selectionRect.x + width &&
                        node.position.y + this.state.y >= this.state.selectionRect.y &&
                        node.position.y + this.state.y < this.state.selectionRect.y + height
                }.bind(this));

                this.setState({
                    selected: selected,
                    selectionRect: {
                        enabled: true,
                        x: this.state.selectionRect.x,
                        y: this.state.selectionRect.y,
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
                    x: e.pageX,
                    y: e.pageY,
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
        selectGroup: function(e) {
            this.setState({
                group: e.id
            })
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
                    key: c.id,
                    model: c,
                    x: c.position.x,
                    y: c.position.y,
                    nodeSelect: this.nodeSelect
                }, React.createElement(nodes[c.instance()], {
                    key: c.id,
                    model: c,
                    selected: this.state.selected.indexOf(c.id) !== -1 ? true : false,
                }, null))
            }.bind(this));

            var edgeElements = this.props.model.focusedEdges.map(function(c) {
                return React.createElement(edges[c.instance()], {
                    key: c.id,
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
