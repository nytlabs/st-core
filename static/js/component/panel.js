var app = app || {};

/* PanelComponent & PanelEditableComponent
 * Produces a list of fields that are the current representation of input
 * values for blocks/groups that are sent to the component.
 *
 * TODO: fix the {'data': ...} nonsense
 */

(function() {
    'use strict';

    app.PanelEditableComponent = React.createClass({
        displayName: 'PanelEditableComponent',
        getInitialState: function() {
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
            var value = this.state.value.length === 0 ? '<empty>' : this.state.value;

            var inputStyle = {
                display: this.state.isEditing ? 'block' : 'none'
            }
            var style = {
                display: this.state.isEditing ? 'none' : 'block'
            }

            return React.createElement('div', {}, [
                React.createElement('div', {
                    className: 'label',
                }, this.props.name),
                React.createElement('div', {
                    className: 'editable'
                }, [
                    React.createElement('input', {
                        defaultValue: this.state.value,
                        onKeyUp: this.handleKeyUp,
                        onBlur: this.handleBlur,
                        style: inputStyle,
                        ref: 'editableInput',
                        key: 'editableInput'
                    }, null),
                    React.createElement('div', {
                        onClick: this.handleClick,
                        style: style,
                        key: 'editableDisplay'
                    }, value)
                ])
            ])
        }
    })
})();

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

(function() {
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
})();

(function() {
    app.RoutesPanelComponent = React.createClass({
        displayName: 'PanelComponent',
        getInitialState: function() {
            var block = app.BlockStore.getBlock(this.props.id);
            return {
                label: block.data.label
            }
        },
        componentDidMount: function() {
            app.BlockStore.addListener(this._onUpdate);
        },
        componentWillUnmount: function() {
            app.BlockStore.removeListener(this._onUpdate);
        },
        shouldComponentUpdate: function(props, state) {
            return state.label != this.state.label
        },
        _onUpdate: function() {
            var block = app.BlockStore.getBlock(this.props.id);
            this.setState({
                label: block.data.label
            })
        },
        _requestChange: function(a) {
            console.log("NEW LABEL:", a);
        },
        render: function() {
            var block = app.BlockStore.getBlock(this.props.id);
            return React.createElement('div', {
                className: 'panel'
            }, [
                React.createElement('div', {
                    key: 'block_header',
                    className: 'block_header',
                }, block.data.type),
                React.createElement(app.PanelEditableComponent, {
                    key: 'label',
                    value: block.data.label,
                    name: 'label',
                    onChange: this._requestChange
                }, null),
                block.inputs.map(function(r, i) {
                    return React.createElement(app.RouteEditableComponent, {
                        key: r.id,
                        id: r.id,
                    }, null)
                }.bind(this))
            ]);
        }
    })
})();
