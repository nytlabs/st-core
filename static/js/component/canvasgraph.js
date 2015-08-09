var app = app || {};

(function() {
    'use strict';

    app.CanvasGraphComponent = React.createClass({
        displayName: "canvas",
        getInitialState: function() {
            return {
                shift: false,
                controlKey: false,
                button: null,
                bufferNodes: document.createElement('canvas'),
                bufferSelection: document.createElement('canvas'),
                mouseDownId: null,
                mouseDownX: null,
                mouseDownY: null,
                mouseLastX: null,
                mouseLastY: null,
                selecting: false,
                selection: []
            }
        },
        shouldComponentUpdate: function() {
            return false;
        },
        componentDidMount: function() {
            // this is a cardinal sin, but how are you supposed to handle
            // canvas-as-state?
            this.state.bufferNodes.width = this.props.width;
            this.state.bufferNodes.height = this.props.height;
            this.state.bufferSelection.width = this.props.width;
            this.state.bufferSelection.height = this.props.height;

            app.BlockStore.addListener(this._onNodesUpdate);
            document.addEventListener('keydown', this._onKeyDown);
            document.addEventListener('keyup', this._onKeyUp);
        },
        componentWillUnmount: function() {
            app.BlockStore.removeListener(this._onNodesUpdate);
            document.removeEventListener('keydown', this._onKeyDown);
            document.removeEventListener('keyup', this._onKeyUp);
        },
        _onKeyDown: function(e) {
            // only fire delete if we have the stage in focus
            if (e.keyCode === 8 && e.target === document.body) {
                e.preventDefault();
                e.stopPropagation();
                this.deleteSelection();
            }

            // only fire ctrl key state if we don't have anything in focus
            if (document.activeElement === document.body &&
                (e.keyCode === 91 || e.keyCode === 17)) {
                this.setState({
                    controlKey: true
                })
            }

            if (e.shiftKey === true) {
                this.setState({
                    shift: true
                })
            }
        },
        _onKeyUp: function(e) {
            if (e.keyCode === 91 || e.keyCode === 17) {
                this.setState({
                    controlKey: false
                })

            }
            if (e.shiftKey === false) this.setState({
                shift: false
            })
        },
        _onMouseDown: function(e) {
            this.setState({
                button: e.button,
                mouseDownX: e.pageX,
                mouseDownY: e.pageY
            })

            var ids = app.BlockStore.pickBlock(e.pageX, e.pageY);

            if (ids.length === 0) {
                if (this.state.shift === false) {
                    app.Dispatcher.dispatch({
                        action: app.Actions.APP_DESELECT_ALL,
                    });
                }
                this.setState({
                    mouseDownId: null
                })
                return
            }


            // pick the first ID
            var id = ids[0];
            if (this.state.shift === true) {
                app.Dispatcher.dispatch({
                    action: app.Actions.APP_SELECT_TOGGLE,
                    ids: [id]
                })
            } else if (app.BlockStore.getSelected().indexOf(id) === -1) {
                app.Dispatcher.dispatch({
                    action: app.Actions.APP_SELECT,
                    id: id
                })
            }
            this.setState({
                mouseDownId: id
            })
        },
        _onMouseUp: function(e) {
            this.setState({
                button: null
            });

            if (this.state.selecting === true) {
                this.setState({
                    selecting: false,
                    selection: []
                });
                this._selectionRectClear();
            }
        },
        _onClick: function(e) {},
        _onMouseMove: function(e) {
            this.setState({
                mouseLastX: e.pageX,
                mouseLastY: e.pageY
            });

            if (this.state.button === 0 && this.state.mouseDownId !== null &&
                this.state.shift === false) {
                app.Dispatcher.dispatch({
                    action: app.Actions.APP_SELECT_MOVE,
                    dx: e.pageX - this.state.mouseLastX,
                    dy: e.pageY - this.state.mouseLastY
                })
            }

            if (this.state.button === 0 && this.state.mouseDownId === null) {
                if (this.state.selected !== true) {
                    this.setState({
                        selecting: true
                    })
                }
                this._selectionRectUpdate(e.pageX, e.pageY);
            }
        },
        _selectionRectClear: function() {
            var ctx = this.state.bufferSelection.getContext('2d');
            ctx.clearRect(0, 0, this.props.width, this.props.height);

            this._renderBuffers();
        },
        _selectionRectUpdate: function(x, y) {
            var width = Math.abs(x - this.state.mouseDownX);
            var height = Math.abs(y - this.state.mouseDownY);
            var originX = Math.min(x, this.state.mouseDownX);
            var originY = Math.min(y, this.state.mouseDownY);
            var selectRect = app.BlockStore.pickArea(originX, originY, width, height);

            // get all nodes new to the selection rect
            var toggles = selectRect.filter(function(id) {
                return this.state.selection.indexOf(id) === -1
            }.bind(this))

            // get all nodes that have left the selection rect
            toggles = toggles.concat(this.state.selection.filter(function(id) {
                return selectRect.indexOf(id) === -1
            }));

            // toggle all new nodes, all nodes that have left the rect
            app.Dispatcher.dispatch({
                action: app.Actions.APP_SELECT_TOGGLE,
                ids: toggles
            })

            this.setState({
                selection: selectRect
            })

            var ctx = this.state.bufferSelection.getContext('2d');
            ctx.clearRect(0, 0, this.props.width, this.props.height);
            ctx.fillStyle = 'rgba(200,200,200,1)';
            ctx.fillRect(originX, originY, width, height);

            this._renderBuffers();
        },
        _onNodesUpdate: function() {
            var nodesCtx = this.state.bufferNodes.getContext('2d');
            nodesCtx.clearRect(0, 0, this.props.width, this.props.height);
            app.BlockStore.getBlocks().forEach(function(id, i) {
                var block = app.BlockStore.getBlock(id);
                nodesCtx.drawImage(block.canvas, block.position.x, block.position.y);
            })

            this._renderBuffers();
        },
        _renderBuffers: function() {
            var ctx = React.findDOMNode(this.refs.test).getContext('2d');
            ctx.clearRect(0, 0, this.props.width, this.props.height);
            ctx.drawImage(this.state.bufferSelection, 0, 0);
            ctx.drawImage(this.state.bufferNodes, 0, 0);
        },
        render: function() {
            return React.createElement('canvas', {
                ref: 'test',
                width: this.props.width,
                height: this.props.height,
                onMouseDown: this._onMouseDown,
                onMouseUp: this._onMouseUp,
                onDoubleClick: this.props.doubleClick,
                onClick: this._onClick,
                onMouseMove: this._onMouseMove,
                onDrag: this._onDrag
            }, null);
        }
    });
})();
