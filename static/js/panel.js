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
            })
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
            var children = [];
            var value = this.props.value.length === 0 ? '<empty>' : this.props.value;

            if (this.state.isEditing) {
                children.push(React.createElement('input', {
                    defaultValue: this.state.value,
                    onKeyUp: this.handleKeyUp,
                    onBlur: this.handleBlur,
                }, null));
            } else {
                children.push(React.createElement('div', {
                    onClick: this.handleClick
                }, value));
            }

            return React.createElement('div', {}, children);
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
                React.createElement(app.PanelEditableComponent, {
                    key: 'label',
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
                this.props.model.inputs.map(function(r, i) {
                    return React.createElement(app.PanelEditableComponent, {
                            value: JSON.stringify(r.data.value),
                            className: 'editable',
                            onChange: function(value) {
                                app.Utils.request(
                                    "PUT",
                                    this.props.model.instance() + "s/" + this.props.model.data.id + "/routes/" + i,
                                    JSON.parse(value),
                                    null
                                )
                            }.bind(this)
                        },
                        null);
                }.bind(this))
            ]);
        }
    })
})();
