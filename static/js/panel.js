var app = app || {};

(function() {
    app.PanelComponent = React.createClass({
        displayName: "PanelComponent",
        render: function() {
            return React.createElement('div', {
                className: 'panel'
            }, this.props.model.data.type);
        }
    })
})();
