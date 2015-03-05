var app = app || {};

(function() {
    app.PanelComponent = React.createClass({
        displayName: "PanelComponent",
        render: function() {
            return React.createElement('div', {
                className: 'panel'
            }, this.props.model.inputs.map(function(r, i) {
                var status;
                if (r.connections.length > 0) {
                    status = "connected to " + r.connections.map(function(c) {
                        return c.data.from.id
                    }).join(",");
                } else
                if (r.data.value != null) {
                    status = JSON.stringify(r.data.value);
                } else {
                    value = "not set";
                }

                return React.createElement('div', {
                    className: 'input',
                    key: i
                }, r.data.name + ' ' + status)

            }));
        }
    })
})();
