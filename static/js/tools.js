var app = app || {};

(function() {
    app.ToolButton = React.createClass({
        displayName: 'ToolButton',
        render: function() {
            return React.createElement('div', {
                className: 'tool_button',
                onClick: this.props.onClick
            }, this.props.label)
        }
    })
})();


(function() {
    app.ToolsComponent = React.createClass({
        displayName: 'ToolsComponent',
        render: function() {
            var groupButton = React.createElement(app.ToolButton, {
                onClick: this.props.onGroup,
                key: 'group',
                label: 'group'
            });

            var ungroupButton = React.createElement(app.ToolButton, {
                onClick: this.props.onUngroup,
                key: 'ungroup',
                label: 'ungroup'
            }, null)

            var tools = [groupButton, ungroupButton]

            return React.createElement('div', {
                className: 'tools'
            }, tools);
        }
    });
})();
