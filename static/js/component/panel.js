var app = app || {};

/* PanelComponent & PanelEditableComponent
 * Produces a list of fields that are the current representation of input
 * values for blocks/groups that are sent to the component.
 *
 */


/* TODO: all routes, source parameters, and labels *should* use the same 
 * component for input. This is ridiculous
 */

(function() {
    app.RoutePanelInput = React.createClass({
        getInitialState: function() {
            return {
                name: '',
                type: '',
                value: '',
            }
        },
        componentDidMount: function() {
            app.RouteStore.getRoute(this.props.id).addListener(this._update);
            this._update();
        },
        componentWillUnmount: function() {
            app.RouteStore.getRoute(this.props.id).removeListener(this._update);
        },
        _update: function() {
            var route = app.RouteStore.getRoute(this.props.id);
            var value = '';
            if (route.data.value !== null) {
                value = JSON.stringify(route.data.value.data);
            }

            this.setState({
                name: route.data.name,
                type: route.data.type,
                value: value,
            })
        },
        _handleChange: function(event) {
            this.setState({
                value: event.target.value
            });
        },
        _onKeyDown: function(event) {
            if (event.keyCode !== 13) return;

            var value = null;
            if (this.state.value !== null) {
                try {
                    value = {
                        data: JSON.parse(this.state.value)
                    }
                } catch (e) {
                    console.log(e);
                }
            }

            this.refs.value.getDOMNode().blur();

            app.Dispatcher.dispatch({
                action: app.Actions.APP_REQUEST_ROUTE_UPDATE,
                id: this.props.id,
                value: value
            })
        },
        render: function() {
            return React.createElement('div', {}, [
                React.createElement('div', {
                    className: 'label',
                    key: 'label',
                }, this.state.name),
                React.createElement('input', {
                    type: 'text',
                    ref: 'value',
                    key: 'value',
                    value: this.state.value,
                    onChange: this._handleChange,
                    onKeyDown: this._onKeyDown,
                }, null)
            ]);
        }
    });
})();

(function() {
    app.GroupNameComponent = React.createClass({
        displayName: 'GroupNameComponent',
        getInitialState: function() {
            return {
                value: JSON.stringify(app.NodeStore.getNode(this.props.id).data.label)
            }
        },
        componentDidMount: function() {
            app.NodeStore.getNode(this.props.id).addListener(this._update);
        },
        componentWillUnmount: function() {
            app.NodeStore.getNode(this.props.id).removeListener(this._update);
        },
        _update: function() {
            this.setState({
                value: JSON.stringify(app.NodeStore.getNode(this.props.id).data.label)
            })
        },
        _handleChange: function(event) {
            this.setState({
                value: event.target.value
            })
        },
        _onKeyDown: function() {
            if (event.keyCode !== 13) return;

            var value = null;
            if (this.state.value !== null) {
                try {
                    value = JSON.parse(this.state.value)
                } catch (e) {
                    console.log(e);
                }
            }

            this.refs.value.getDOMNode().blur();

            app.Dispatcher.dispatch({
                action: app.Actions.APP_REQUEST_NODE_LABEL,
                id: this.props.id,
                label: value
            })
        },
        render: function() {
            return React.createElement('div', {}, [
                React.createElement('div', {
                    className: 'label',
                    key: 'label',
                }, 'label'),
                React.createElement('input', {
                    type: 'text',
                    ref: 'value',
                    key: 'value',
                    value: this.state.value,
                    onChange: this._handleChange,
                    onKeyDown: this._onKeyDown,
                }, null)
            ]);

        }
    })
})();

/* Input component for a source paramter 
 */
(function() {
    app.SourceParameterComponent = React.createClass({
        displayName: 'SourceParameterComponent',
        getInitialState: function() {
            // this is -- the worst.
            // needs to be sorted out on API level
            return {
                value: JSON.stringify(app.NodeStore.getNode(this.props.id).data.params.filter(function(param) {
                    return param.name === this.props.name
                }.bind(this))[0].value)
            }
        },
        componentDidMount: function() {
            app.NodeStore.getNode(this.props.id).addListener(this._update);
        },
        componentWillUnmount: function() {
            app.NodeStore.getNode(this.props.id).removeListener(this._update);
        },
        _update: function() {
            // the absolute worst.
            this.setState({
                value: JSON.stringify(app.NodeStore.getNode(this.props.id).data.params.filter(function(param) {
                    return param.name === this.props.name
                }.bind(this))[0].value)
            })
        },
        _handleChange: function(event) {
            this.setState({
                value: event.target.value
            })
        },
        _onKeyDown: function() {
            if (event.keyCode !== 13) return;

            var value = null;
            if (this.state.value !== null) {
                try {
                    value = JSON.parse(this.state.value)
                } catch (e) {
                    console.log(e);
                    return;
                }
            }

            this.refs.value.getDOMNode().blur();

            app.Dispatcher.dispatch({
                action: app.Actions.APP_REQUEST_SOURCE_PARAMS,
                id: this.props.id,
                name: this.props.name,
                value: value
            })
        },
        render: function() {
            return React.createElement('div', {}, [
                React.createElement('div', {
                    className: 'label',
                    key: 'label',
                }, this.props.name),
                React.createElement('input', {
                    type: 'text',
                    ref: 'value',
                    key: 'value',
                    value: this.state.value,
                    onChange: this._handleChange,
                    onKeyDown: this._onKeyDown,
                }, null)
            ]);

        }
    })
})();

/* Panel Component for an individual node */
(function() {
    app.RoutesPanelComponent = React.createClass({
        displayName: 'PanelComponent',
        componentDidMount: function() {
            app.NodeStore.getNode(this.props.id).addListener(this._update);
            this._update();
        },
        componentWillUnmount: function() {
            app.NodeStore.getNode(this.props.id).removeListener(this._update);
        },
        _update: function() {
            this.render();
            /*var route = app.RouteStore.getRoute(this.props.id);
            var value = '';
            if (route.data.value !== null) {
                value = JSON.stringify(route.data.value.data);
            }

            this.setState({
                name: route.data.name,
                type: route.data.type,
                value: value,
            })*/
        },
        render: function() {
            var block = app.NodeStore.getNode(this.props.id);

            var children = [
                React.createElement('div', {
                    key: 'block_header',
                    className: 'block_header',
                }, block.data.type),
                React.createElement(app.GroupNameComponent, {
                    id: this.props.id
                }, null)
            ];


            // TODO: optimize this!
            // this retrieves _all_ routes for a block, seems unnecesary
            children = children.concat(block.routes.filter(function(id) {
                return app.RouteStore.getRoute(id).direction === 'input';
            }).map(function(id) {
                return React.createElement(app.RoutePanelInput, {
                    blockId: this.props.id,
                    id: id,
                    key: id,
                }, null)
            }.bind(this)));

            if (block instanceof app.Source) {
                children = children.concat(block.data.params.map(function(param) {
                    return React.createElement(app.SourceParameterComponent, {
                        name: param.name,
                        value: param.value,
                        id: block.data.id
                    }, null)
                }))
            }

            return React.createElement('div', {
                className: 'panel'
            }, children);
        }
    })
})();
