var app = app || {};

/* ToolComponent & ToolButtonComponent
 * Simple toolbar (lower left)  that may go away at some point
 */

(function() {
    'use strict';

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
    'use strict';

    app.ToolsComponent = React.createClass({
        displayName: 'ToolsComponent',
        shouldComponentUpdate: function(props, state) {
            return false;
        },
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
