var app = app || {};

(function() {
    app.PanelEditableComponent = React.createClass({
        displayName: "PanelEditableComponent",
        getInitialState: function() {
            return {
                isEditing: false,
                // yes yes anti-pattern but yet another circumstance where it
                // totally fine to have 2 separate models representing the 
                // "same" state.
                //
                // One state is meant to be user-editable and is sent as a req
                // to the server, the other state is always displaying the 
                // current state synced w/ server. we populate default for 
                // state that is meant to make the req with what is set by
                // the server.
                value: this.props.value
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
            var value = this.props.value.length === 0 ? '<empty>' : this.props.value;
            var inputStyle = {
                display: this.state.isEditing ? 'block' : 'none'
            }
            var style = {
                display: this.state.isEditing ? 'none' : 'block'
            }

            return React.createElement('div', {
                className: 'editable'
            }, [
                React.createElement('input', {
                    defaultValue: this.state.value,
                    onKeyUp: this.handleKeyUp,
                    onBlur: this.handleBlur,
                    style: inputStyle,
                    ref: "editableInput",
                    key: "editableInput"
                }, null),
                React.createElement('div', {
                    onClick: this.handleClick,
                    style: style,
                    key: 'editableDisplay'
                }, value)
            ]);
        }
    })
})();


(function() {
    app.PanelComponent = React.createClass({
        displayName: "PanelComponent",
        render: function() {
            return React.createElement('div', {
                className: 'panel'
            }, [
                React.createElement('div', {
                    key: 'block_header',
                    className: 'block_header',
                }, this.props.model.data.type),
                React.createElement('div', {
                    key: 'block_Label',
                    className: "label",
                }, "label"),
                React.createElement(app.PanelEditableComponent, {
                    key: 'route_label',
                    className: 'editable',
                    value: this.props.model.data.label,
                    onChange: function(value) {
                        app.Utils.request(
                            "PUT",
                            this.props.model.instance() + "s/" + this.props.model.data.id + "/label",
                            value,
                            null
                        )
                    }.bind(this)
                }, null),
                this.props.model.routes.filter(function(r) {
                    return r.direction === 'input'
                }).map(function(r, i) {
                    return [
                        React.createElement('div', {
                            className: 'label',
                        }, r.data.name),
                        React.createElement(app.PanelEditableComponent, {
                                key: r.id + r.data.name + r.index,
                                value: JSON.stringify(r.data.value),
                                onChange: function(value) {
                                    app.Utils.request(
                                        "PUT",
                                        "blocks/" + r.id + "/routes/" + r.index,
                                        JSON.parse(value),
                                        null
                                    )
                                }.bind(this)
                            },
                            null)
                    ]
                }.bind(this))
            ]);
        }
    })
})();
