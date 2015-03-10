var app = app || {};

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
        updateRoute: function(e, index, value) {
            app.Utils.request(
                "PUT",
                this.props.model.instance() + "s/" + this.props.model.data.id + "/routes/" + index,
                JSON.parse(value),
                null
            )
        },
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
                    return React.createElement(app.PanelInputComponent, {
                        updateRoute: this.updateRoute,
                        model: r,
                        index: i,
                    }, null);
                }.bind(this))
            ]);
        }
    })
})();
