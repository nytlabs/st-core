var app = app || {};

(function() {
    app.PanelEditableComponent = React.createClass({
        displayName: "PanelEditableComponent",
        getInitialState: function() {
            return {
                isEditing: false,
                value: this.props.value
            }
        },
        handleClick: function() {
            this.setState({
                isEditing: true
            })
        },
        handleKeyUp: function(e) {
            if (e.nativeEvent.keyCode === 13) {
                this.props.onChange(e.target.value);
                this.setState({
                    isEditing: false,
                })
            }
        },
        render: function() {
            var children = [];
            if (this.state.isEditing) {
                children.push(React.createElement('input', {
                    defaultValue: this.state.value,
                    onKeyUp: this.handleKeyUp,
                }, null));
            } else {
                children.push(React.createElement('div', {
                    onClick: this.handleClick
                }, this.state.value));
            }

            return React.createElement('div', {}, children);
        }
    })
})();

(function() {
    app.PanelInputComponent = React.createClass({
        displayName: "PanelInputCompnent",
        getInitialState: function() {
            return {
                isEditing: false
            }
        },
        onKeyPress: function(e) {
            if (e.nativeEvent.keyCode === 13) {
                this.props.updateRoute(e, this.props.index, e.target.value);
                this.setState({
                    isEditing: false
                })
            }
        },
        handleClick: function() {
            this.setState({
                isEditing: true
            })
        },
        render: function() {
            var child = [];

            child.push(React.createElement('div', {
                key: 'route_name',
                className: 'route_name'
            }, this.props.model.data.name));

            if (this.props.model.connections.length > 0) {
                child.push(React.createElement('div', {
                    key: 'connected',
                    className: 'connected'
                }, this.props.model.connections.map(function(c) {
                    return c.data.from.id
                }).join(", ")));
            } else {
                if (this.state.isEditing) {
                    child.push(React.createElement('input', {
                        key: 'route_value',
                        className: 'route_value',
                        defaultValue: JSON.stringify(this.props.model.data.value),
                        onKeyPress: this.onKeyPress,
                    }))
                } else {
                    child.push(React.createElement('div', {
                        key: 'route_value',
                        className: 'route_value',
                        onClick: this.handleClick,
                    }, JSON.stringify(this.props.model.data.value)))
                }
            }

            return React.createElement('div', {
                className: 'input',
                key: this.props.index,
            }, child);
        }
    })
})();


(function() {
    app.PanelComponent = React.createClass({
        displayName: "PanelComponent",
        updateLabel: function(e) {
            if (e.nativeEvent.keyCode === 13) {
                app.Utils.request(
                    "PUT",
                    this.props.model.instance() + "s/" + this.props.model.data.id + "/label",
                    e.target.value,
                    null
                )
            }
        },
        updateRoute: function(e, index, value) {},
        render: function() {
            return React.createElement('div', {
                className: 'panel'
            }, [
                React.createElement('div', {
                    key: 'block_header',
                    className: 'block_header',
                }, this.props.model.data.type),
                React.createElement('div', {
                    key: 'label',
                    className: 'route_name',
                }, "label"),
                React.createElement('input', {
                    key: 'label_input',
                    className: 'route_value',
                    defaultValue: this.props.model.data.label,
                    onKeyPress: this.updateLabel
                }, null),
                this.props.model.inputs.map(function(r, i) {
                    console.log(r.data.value);
                    return React.createElement(app.PanelEditableComponent, {
                            value: JSON.stringify(r.data.value),
                            onChange: function(value) {
                                console.log("NEW:", value)
                                app.Utils.request(
                                    "PUT",
                                    this.props.model.instance() + "s/" + this.props.model.data.id + "/routes/" + i,
                                    JSON.parse(value),
                                    null
                                )
                            }.bind(this)
                        },
                        null)

                    /*return React.createElement(app.PanelInputComponent, {
                        updateRoute: this.updateRoute,
                        model: r,
                        index: i,
                    }, null);*/
                }.bind(this))
            ]);
        }
    })
})();
