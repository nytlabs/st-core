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
                //value: this.props.value
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
            /*var value = this.props.value.length === 0 ? '<empty>' : this.props.value;
            var inputStyle = {
                display: this.state.isEditing ? 'block' : 'none'
            }
            var style = {
                display: this.state.isEditing ? 'none' : 'block'
            }

            var children = [React.createElement('div', {
                className: 'label',
            }, r.data.name)]

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
            ]);*/
            return React.createElement('div', {}, null);
        }
    })
})();

(function() {

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

}());

(function() {
    app.RoutesPanelComponent = React.createClass({
        displayName: 'PanelComponent',
        render: function() {
            var block = app.BlockStore.getBlock(this.props.id);
            console.log(this.props.key);
            console.log(block);
            return React.createElement('div', {
                className: 'panel'
            }, [
                React.createElement('div', {
                    key: 'block_header',
                    className: 'block_header',
                }, block.data.type),
                React.createElement('div', {
                    key: 'block_Label',
                    className: 'label',
                }, 'label'),
                React.createElement(app.PanelEditableComponent, {
                    key: 'route_label',
                    className: 'editable',
                    value: block.data.label,
                    onChange: function() {},
                }, null),
                block.inputs.map(function(r, i) {
                    console.log(r);
                    return [
                        React.createElement(app.PanelEditableComponent, {
                            key: r.id,
                            onChange: function() {},
                        }, null)
                    ]
                }.bind(this))
            ]);
        }
    })
})();
