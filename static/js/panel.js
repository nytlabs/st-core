var app = app || {};

(function() {
    app.PanelComponent = React.createClass({
        displayName: "PanelComponent",
        render: function() {
            return React.createElement('div', {
                className: 'panel'
            }, this.props.model.data.inputs.map(function(r, i) {
                return React.createElement('div', {
                    className: 'input',
                    key: i
                }, r.name + ' ' + JSON.stringify(r.value))

            }));
        }
    })
})();
