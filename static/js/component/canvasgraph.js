var app = app || {};

(function() {
    'use strict';

    app.CanvasGraphComponent = React.createClass({
        displayName: "canvas",
        getInitialState: function() {
            return {
                shift: false,
                control: false,
                button: null,
                bufferNodes: document.createElement('canvas'),
                bufferSelection: document.createElement('canvas'),
                bufferStage: document.createElement('canvas'),
                bufferEdgeTool: document.createElement('canvas'),
                bufferEdges: document.createElement('canvas'),
                bufferPicking: document.createElement('canvas'),
                mouseDownId: null,
                mouseDownX: null,
                mouseDownY: null,
                mouseLastX: null,
                mouseLastY: null,
                // the selection area/rect
                selecting: false,
                selection: [],
                // the panning offset
                translateX: 0, // TODO: do per-group translations, deprecate this
                translateY: 0,
                // the connection tool
                connectingBlock: null,
                connectingRoute: null,
                dirty: false,
                width: 0,
                height: 0,
                offX: 0,
                offY: 0
            }
        },
        shouldComponentUpdate: function() {
            return false;
        },
        componentDidMount: function() {
            app.NodeStore.addListener(this._onNodesUpdate);
            app.EdgeStore.addListener(this._onEdgesUpdate);

            window.addEventListener('keydown', this._onKeyDown);
            window.addEventListener('keyup', this._onKeyUp);
            window.addEventListener('resize', this._onResize);
            window.addEventListener('copy', this._onCopy);
            window.addEventListener('paste', this._onPaste);

            this._onResize();
            this._renderBuffers();
        },
        componentWillUnmount: function() {
            app.NodeStore.removeListener(this._onNodesUpdate);
            app.EdgeStore.removeListener(this._onEdgesUpdate);

            window.removeEventListener('keydown', this._onKeyDown);
            window.removeEventListener('keyup', this._onKeyUp);
            window.removeEventListener('resize', this._onResize);
            window.removeEventListener('copy', this._onCopy);
            window.removeEventListener('paste', this._onPaste);
        },
        _onCopy: function(e) {
            if (e.target.nodeName !== 'input') {
                e.preventDefault();
                e.clipboardData.setData('text/plain', JSON.stringify(app.SelectionStore.getPattern()));
            }
        },
        _onPaste: function(e) {
            if (e.target.nodeName !== 'input') {
                app.Dispatcher.dispatch({
                    action: app.Actions.APP_REQUEST_GROUP_IMPORT,
                    pattern: e.clipboardData.getData('text'),
                });
            }
        },
        _onResize: function(e) {
            var width = document.body.clientWidth;
            var height = document.body.clientHeight;

            // resize all the buffers
            this.state.bufferNodes.width = width;
            this.state.bufferNodes.height = height;
            this.state.bufferSelection.width = width;
            this.state.bufferSelection.height = height;
            this.state.bufferStage.width = width;
            this.state.bufferStage.height = height;
            this.state.bufferEdgeTool.width = width;
            this.state.bufferEdgeTool.height = height;
            this.state.bufferEdges.width = width;
            this.state.bufferEdges.height = height;
            this.state.bufferPicking.width = width;
            this.state.bufferPicking.height = height;

            // resize the main canvas
            React.findDOMNode(this.refs.test).width = width;
            React.findDOMNode(this.refs.test).height = height;

            // render everything again
            this.setState({
                dirty: true,
                width: width,
                height: height,
                offX: Math.floor(width * .5),
                offY: Math.floor(height * .5)
            }, function() {
                this._renderStage();
                this._renderEdges();
                this._renderNodes();
            }.bind(this));
        },
        _onKeyDown: function(e) {
            // only fire delete if we have the stage in focus
            if (e.keyCode === 8 && e.target === document.body) {
                e.preventDefault();
                e.stopPropagation();
                app.Dispatcher.dispatch({
                    action: app.Actions.APP_DELETE_SELECTION,
                })
            }

            // only fire ctrl key state if we don't have anything in focus
            if (document.activeElement === document.body &&
                (e.keyCode === 91 || e.keyCode === 17)) {
                this.setState({
                    control: true
                })
            }

            if (e.shiftKey === true) {
                this.setState({
                    shift: true
                })
            }

            if (document.activeElement === document.body && e.keyCode === 71) {
                app.Dispatcher.dispatch({
                    action: app.Actions.APP_GROUP_SELECTION
                })
            }

            if (document.activeElement === document.body && e.keyCode === 85) {
                app.Dispatcher.dispatch({
                    action: app.Actions.APP_UNGROUP_SELECTION
                })
            }
        },
        _onKeyUp: function(e) {
            if (e.keyCode === 91 || e.keyCode === 17) {
                this.setState({
                    control: false
                })

            }
            if (e.shiftKey === false) {
                this.setState({
                    shift: false
                })
            }
        },
        /* given coordinates, return a node */
        _pickBuffer: function(x, y) {
            var ctx = this.state.bufferPicking.getContext('2d');
            var col = ctx.getImageData(x, y, 1, 1).data;
            var picked = null;
            // important! this throws away anti-aliased parts of a line!
            if (col[3] === 255) {
                var colString = "rgb(" + col[0] + "," + col[1] + "," + col[2] + ")";
                picked = app.PickingStore.colorToNode(colString);
            }
            return picked;
        },
        _onMouseDown: function(e) {
            this.setState({
                button: e.button,
                mouseDownX: e.pageX,
                mouseDownY: e.pageY
            })

            this._renderPickingBuffer();

            var picked = this._pickBuffer(e.pageX, e.pageY);
            var isElement = picked instanceof app.Node || picked instanceof app.Edge;
            var isRoute = picked instanceof app.Route;

            if (this.state.connectingBlock !== null) {
                this._connectingClear();
            }

            if (picked === null) {
                if (this.state.shift === false) {
                    app.Dispatcher.dispatch({
                        action: app.Actions.APP_DESELECT_ALL,
                    });
                }
                this.setState({
                    mouseDownId: null,
                    connectingBlock: null,
                    connectingRoute: null
                })
                return
            }

            if (isRoute && this.state.connectingBlock === null) {
                var block = app.NodeStore.getNode(picked.visibleParent);
                var route = block.routeGeometry[picked.id];
                this.setState({
                    connectingBlock: block,
                    connectingRoute: route,
                })
                return
            } else if (isRoute && this.state.connectingBlock !== null) {
                var route = app.RouteStore.getRoute(this.state.connectingRoute.id);
                app.Dispatcher.dispatch({
                    action: app.Actions.APP_REQUEST_CONNECTION,
                    routes: [picked, route],
                });
            } else if (isElement && this.state.shift === true) {
                app.Dispatcher.dispatch({
                    action: app.Actions.APP_SELECT_TOGGLE,
                    ids: [picked]
                })
            } else if (isElement && !app.SelectionStore.isSelected(picked)) {
                app.Dispatcher.dispatch({
                    action: app.Actions.APP_SELECT,
                    id: picked
                })
            }

            this.setState({
                mouseDownId: true, // fix this its stupid, it doesn't need to nbe an id apparently?
                connectingBlock: null,
                connectingRoute: null
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

            if (this.state.mouseDownId !== null && this.state.shift === false) {
                app.Dispatcher.dispatch({
                    action: app.Actions.APP_REQUEST_NODE_MOVE,
                });
            }
        },
        _onDoubleClick: function(e) {
            var p = this._pickBuffer(e.pageX, e.pageY);
            if (p === null) {
                this.props.showAutoComplete(
                    e.pageX, e.pageY,
                    e.pageX - this.state.offX + -1 * this.state.translateX,
                    e.pageY - this.state.offY + -1 * this.state.translateY
                );
            }
            if (p instanceof app.Group) {
                app.NodeStore.setRoot(p.data.id);
            }
        },
        _onContextMenu: function(e) {
            e.nativeEvent.preventDefault();
        },
        _onMouseMove: function(e) {
            this.setState({
                mouseLastX: e.pageX,
                mouseLastY: e.pageY
            });

            if (this.state.connectingBlock !== null) {
                this._connectingUpdate(e.pageX, e.pageY);
            } else if (this.state.button === 0 && this.state.mouseDownId !== null &&
                this.state.shift === false) {
                app.Dispatcher.dispatch({
                    action: app.Actions.APP_SELECT_MOVE,
                    dx: e.pageX - this.state.mouseLastX,
                    dy: e.pageY - this.state.mouseLastY
                })
            } else if (this.state.button === 0 && this.state.mouseDownId === null) {
                if (this.state.selected !== true) {
                    this.setState({
                        selecting: true
                    })
                }
                this._selectionRectUpdate(e.pageX, e.pageY);
            } else if (this.state.button === 2) {
                var dx = e.pageX - this.state.mouseLastX;
                var dy = e.pageY - this.state.mouseLastY;
                // TODO: get rid of translate in favor of per-group translations
                this.setState({
                    translateX: this.state.translateX + dx,
                    translateY: this.state.translateY + dy
                }, function() {
                    this._onStageUpdate()
                }.bind(this));
            }
        },
        _selectionRectClear: function() {
            var ctx = this.state.bufferSelection.getContext('2d');
            ctx.clearRect(0, 0, this.state.width, this.state.height);

            this.setState({
                dirty: true
            });
        },
        _connectingClear: function() {
            var ctx = this.state.bufferEdgeTool.getContext('2d');
            ctx.clearRect(0, 0, this.state.width, this.state.height);

            //this._renderBuffers();
            this.setState({
                dirty: true
            });
        },
        _connectingUpdate: function(mx, my) {
            var block = this.state.connectingBlock;
            var x = block.position.x + this.state.translateX + this.state.connectingRoute.x;
            var y = block.position.y + this.state.translateY + this.state.connectingRoute.y;
            var ctx = this.state.bufferEdgeTool.getContext('2d');
            var direction = this.state.connectingRoute.direction === 'input' ? -1 : 1;

            ctx.clearRect(0, 0, this.state.width, this.state.height);
            ctx.beginPath()
            ctx.moveTo(this.state.offX + x, this.state.offY + y);
            ctx.setLineDash([5, 5]);
            ctx.lineWidth = 2.0
            ctx.bezierCurveTo(this.state.offX + x + (50 * direction),
                this.state.offY + y,
                mx + (-50 * direction),
                my,
                mx,
                my);
            ctx.stroke();

            this.setState({
                dirty: true
            });
        },
        _selectionRectUpdate: function(x, y) {
            var width = Math.abs(x - this.state.mouseDownX);
            var height = Math.abs(y - this.state.mouseDownY);
            var originX = Math.min(x, this.state.mouseDownX);
            var originY = Math.min(y, this.state.mouseDownY);
            // TODO: get rid of translate in favor of per-group translations
            var selectRect = app.NodeStore.getNodes().filter(function(id) {
                var block = app.NodeStore.getNode(id);
                return app.Utils.pointInRect(
                    originX - this.state.translateX - this.state.offX,
                    originY - this.state.translateY - this.state.offY,
                    width, height,
                    block.position.x, block.position.y);
            }.bind(this)).map(function(id) {
                return app.NodeStore.getNode(id)
            })

            var selectConn = app.EdgeStore.getEdges().filter(function(id) {
                var connection = app.EdgeStore.getEdge(id);
                return app.Utils.pointInRect(
                    originX - this.state.translateX - this.state.offX,
                    originY - this.state.translateY - this.state.offY,
                    width, height,
                    connection.position.x, connection.position.y);
            }.bind(this)).map(function(id) {
                return app.EdgeStore.getEdge(id)
            })

            var selectRect = selectRect.concat(selectConn);

            // get all nodes new to the selection rect
            var toggles = selectRect.filter(function(id) {
                return this.state.selection.indexOf(id) === -1
            }.bind(this))

            // get all nodes that have left the selection rect
            toggles = toggles.concat(this.state.selection.filter(function(id) {
                return selectRect.indexOf(id) === -1
            }));

            if (toggles.length > 0) {
                // toggle all new nodes, all nodes that have left the rect
                app.Dispatcher.dispatch({
                    action: app.Actions.APP_SELECT_TOGGLE,
                    ids: toggles
                })
            }

            this.setState({
                selection: selectRect
            })

            var ctx = this.state.bufferSelection.getContext('2d');
            ctx.clearRect(0, 0, this.state.width, this.state.height);
            ctx.fillStyle = 'rgba(200,200,200,.5)';
            ctx.fillRect(originX, originY, width, height);

            this.setState({
                dirty: true
            });
        },
        _onStageUpdate: function() {
            this._renderStage();
            this._renderEdges();
            this._renderNodes();
            this.setState({
                dirty: true
            });
        },
        _onNodesUpdate: function() {
            this._renderNodes();
            this.setState({
                dirty: true
            });
        },
        _onEdgesUpdate: function() {
            this._renderEdges();
            this.setState({
                dirty: true
            });
        },
        _renderStage: function() {
            var ctx = this.state.bufferStage.getContext('2d');
            var width = this.state.bufferStage.width;
            var height = this.state.bufferStage.height;
            var GRID_PX = 50.0;
            //TODO: get rid in favor of per group translations
            var translateX = this.state.translateX;
            var translateY = this.state.translateY;
            var x = (this.state.offX + translateX) % GRID_PX;
            var y = (this.state.offY + translateY) % GRID_PX;
            var lines = [];
            var hMax = Math.floor(width / GRID_PX);
            var vMax = Math.floor(height / GRID_PX);

            ctx.clearRect(0, 0, width, height);
            ctx.strokeStyle = 'rgb(220,220,220)';

            var grid = new Path2D();
            for (var i = 0; i <= hMax; i++) {
                grid.moveTo(x + (i * GRID_PX), 0);
                grid.lineTo(x + (i * GRID_PX), height);
            }
            for (var i = 0; i <= vMax; i++) {
                grid.moveTo(0, y + (i * GRID_PX));
                grid.lineTo(width, y + (i * GRID_PX));
            }
            ctx.stroke(grid);

            ctx.beginPath();
            ctx.fillStyle = 'rgba(255,0,0,1)';
            ctx.arc(this.state.offX + translateX,
                this.state.offY + translateY,
                3,
                0,
                2 * Math.PI
            );
            ctx.fill();
        },
        _renderNodes: function() {
            var nodesCtx = this.state.bufferNodes.getContext('2d');
            nodesCtx.clearRect(0, 0, this.state.width, this.state.height);
            app.NodeStore.getNodes().forEach(function(id, i) {
                var block = app.NodeStore.getNode(id);
                var x = block.position.x + this.state.translateX + this.state.offX;
                var y = block.position.y + this.state.translateY + this.state.offY;
                nodesCtx.drawImage(block.canvas, x, y);
            }.bind(this))
        },
        _renderEdges: function() {
            var ctx = this.state.bufferEdges.getContext('2d');
            ctx.clearRect(0, 0, this.state.width, this.state.height);
            app.EdgeStore.getEdges().forEach(function(id, i) {
                var connection = app.EdgeStore.getEdge(id);
                var x = connection.position.x + this.state.translateX + this.state.offX;
                var y = connection.position.y + this.state.translateY + this.state.offY;
                ctx.drawImage(connection.canvas, x, y);
            }.bind(this));
        },
        _renderPickingBuffer: function() {
            // renders all the picking-elements to a single canvas in the same
            // order as to how the non-picking elements are rendered.
            // TODO: d-r-y 
            var pickCtx = this.state.bufferPicking.getContext('2d');
            pickCtx.clearRect(0, 0, this.state.width, this.state.height);
            app.NodeStore.getNodes().forEach(function(id, i) {
                var block = app.NodeStore.getNode(id);
                var x = block.position.x + this.state.translateX + this.state.offX;
                var y = block.position.y + this.state.translateY + this.state.offY;
                pickCtx.drawImage(block.pickCanvas, x, y);
            }.bind(this));
            app.EdgeStore.getEdges().forEach(function(id, i) {
                var connection = app.EdgeStore.getEdge(id);
                // we need to update the picking image for each connection 
                // that has been moved.
                if (connection.dirtyPicking) {
                    app.Dispatcher.dispatch({
                        action: app.Actions.APP_RENDER_CONNECTION_PICKING,
                        id: id
                    })
                }
                var x = connection.position.x + this.state.translateX + this.state.offX;
                var y = connection.position.y + this.state.translateY + this.state.offY;
                pickCtx.drawImage(connection.pickCanvas, x, y);
            }.bind(this));
        },
        _renderBuffers: function() {
            // this is getting into the weeds of optimization, but the
            // state.dirty flag batches renders into 16.667ms frames so that we
            // don't encounter a situation where our model fires renders faster
            // than 60fps.
            //
            // for example: it's possible for model updates to occur much
            // faster than 60fps. in the case of a large import, or a single
            // event causing multiple renders, our CanvasGraphComponent may
            // receive events that are less than 16.667 apart. this can lock
            // the interface and cause lag.
            //
            // the cost is that state.dirty/requestAnimationFrame spin while
            // waiting for an update. on my machine, this amounts to ~6% CPU
            // while doing nothing. (as opposed to ~1% of when not calling
            // requestAnimationFrame)
            //
            // TODO: a potential optimization, to maybe save some battery and
            // ultimately the planet earth, would be to make it so that a
            // model update triggers a span of time, like 60 frames, until
            // requestAnimationFrame clears itself. each render would tick
            // down some of the "dirty time" until we reach  0, when we clear
            // requestAnimationFrame, effectively debouncing render events
            // and ceasing the spinning when not doing anything.
            window.requestAnimationFrame(this._renderBuffers);
            if (this.state.dirty) {
                this.setState({
                    dirty: false
                });
                var ctx = React.findDOMNode(this.refs.test).getContext('2d');
                ctx.clearRect(0, 0, this.state.width, this.state.height);
                ctx.drawImage(this.state.bufferStage, 0, 0);
                ctx.drawImage(this.state.bufferSelection, 0, 0);
                ctx.drawImage(this.state.bufferNodes, 0, 0);
                ctx.drawImage(this.state.bufferEdges, 0, 0);
                ctx.drawImage(this.state.bufferEdgeTool, 0, 0);
            }
        },
        render: function() {
            return React.createElement('canvas', {
                ref: 'test',
                width: this.state.width,
                height: this.state.height,
                onMouseDown: this._onMouseDown,
                onMouseUp: this._onMouseUp,
                onClick: this.props.onClick,
                onDoubleClick: this._onDoubleClick,
                onMouseMove: this._onMouseMove,
                onContextMenu: this._onContextMenu,
            }, null);
        }
    });
})();
