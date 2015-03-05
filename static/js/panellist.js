var app = app || {};

(function() {
    app.PanelListComponent = React.createClass({
        displayName: "PanelListComponent",
        render: function() {
            return React.createElement('div', {
                className: "panel_list"
            }, this.props.nodes.map(function(n) {
                return React.createElement(app.PanelComponent, {
                    model: n,
                    key: n.data.id
                }, null)
            }))
        },
    })
})();
