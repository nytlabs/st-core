var app = app || {};

/* PanelComponent & PanelEditableComponent
 * Produces a list of fields that are the current representation of input
 * values for blocks/groups that are sent to the component.
 *
 * TODO: fix the {'data': ...} nonsense
 */

/*(function() {
    'use strict';

    app.PanelEditableComponent = React.createClass({
        displayName: 'PanelEditableComponent',
        getInitialState: function() {
            //console.log(app.RouteStore.getRoute(this.props.id));
            return {
                isEditing: false,
                value: this.props.value,
            }
        },
        handleClick: function() {
            this.setState({
                isEditing: true,
                value: this.props.value
            }, function() {
                this.refs.editableInput.getDOMNode().focus();
                this.refs.editableInput.getDOMNode().select();
            });
        },
        handleKeyUp: function(e) {
            if (e.nativeEvent.keyCode === 13) {
                this.props.onChange(e.target.value);
                this.setState({
                    isEditing: false,
                });
            }
        },
        handleBlur: function() {
            this.setState({
                isEditing: false,
            });
        },
        render: function() {
            return React.createElement('div', {}, [
                React.createElement('div', {
                    className: 'label',
                    key: 'label',
                }, null),
                React.createElement('input', {
                    ref: 'input',
                    key: 'input',
                }, null)
            ])
        }
    })
})();*/

/*(function() {

    app.ParametersPanelComponent = React.createClass({
        displayName: 'ParametersPanelComponent',
        render: function() {
            var id = this.props.model.data.id
            return React.createElement('div', {
                className: 'panel'
            }, [
                React.createElement('div', {
                    key: 'block_header',
                    className: 'block_header'
                }, this.props.model.data.type),
                React.createElement('div', {
                    key: 'block_label',
                    className: 'label'
                }, 'label'),
                React.createElement(app.PanelEditableComponent, {
                    key: 'route_label',
                    className: 'editable',
                    value: this.props.model.data.label,
                    onChange: function(value) {
                        app.Utils.request(
                            'PUT',
                            this.props.model.instance() + 's/' + this.props.model.data.id + '/label',
                            value,
                            null
                        )
                    }.bind(this)
                }, null),
                this.props.model.data.params.map(function(p, i) {
                    return [
                        React.createElement('div', {
                            className: 'label',
                        }, p.name),
                        React.createElement(app.PanelEditableComponent, {
                            key: id + p.name,
                            value: p.value,
                            onChange: function(value) {
                                app.Utils.request('PUT', 'sources/' + id + '/params', [{
                                    name: p.name,
                                    value: value
                                }], null)
                            }.bind(this)
                        }, null)
                    ]

                })
            ])
        }
    })

}());*/

/*(function() {
    'use strict';

    app.RouteEditableComponent = React.createClass({
        getInitialState: function() {
            var route = app.RouteStore.getRoute(this.props.id);
            return {
                name: route.data.name,
                type: route.data.type,
                value: route.data.value
            }
        },
        componentDidMount: function() {
            app.RouteStore.addListener(this._onUpdate);
        },
        componentWillUnmount: function() {
            app.RouteStore.addListener(this._onUpdate);
        },
        shouldComponentUpdate: function(props, state) {
            // TODO: double check objects passed through here.
            return this.state.value !== state.value
        },
        _onUpdate: function() {
            var route = this.RouteStore.getRoute(this.props.id);
            this.setState({
                name: route.data.name,
                type: route.data.type,
                value: route.data.value
            })
        },
        _requestChange: function(a) {
            console.log("I'm Changing:", a);
        },
        render: function() {
            return React.createElement(app.PanelEditableComponent, {
                name: this.state.name,
                value: JSON.stringify(this.state.value),
                onChange: this._requestChange,
            }, null);
        }
    });
})();*/

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
                    return
                }
            }

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
    app.RoutesPanelComponent = React.createClass({
        displayName: 'PanelComponent',
        render: function() {
            var block = app.BlockStore.getBlock(this.props.id);

            var children = [
                React.createElement('div', {
                    key: 'block_header',
                    className: 'block_header',
                }, block.data.type),
            ];

            children = children.concat(block.inputs.map(function(r) {
                return React.createElement(app.RoutePanelInput, {
                    id: r.id,
                    key: r.id,
                }, null)
            }));

            return React.createElement('div', {
                className: 'panel'
            }, children);
        }
    })
})();
