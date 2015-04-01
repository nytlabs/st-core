var app = app || {};

/* GroupSelectorComponent
 * Displays menu that lists all groups.
 * Clicking a group updates CoreModel with the currently in-view group.
 * Changes to CoreModel update which group is selected, and which nods are
 * currently in view.
 */

(function() {
    'use strict';

    app.GroupSelectorComponent = React.createClass({
        displayName: 'GroupSelectorComponent',
        onClick: function(group) {
            group.setFocusedGroup();
        },
        render: function() {
            return React.createElement('div', {
                className: 'group_list',
            }, this.props.groups.map(function(g) {
                var classes = 'group';
                if (this.props.focusedGroup === g.data.id) classes += ' focused';
                return React.createElement('div', {
                    className: classes,
                    key: g.data.id,
                    onClick: this.onClick.bind(null, g),
                }, g.data.label + ' ' + g.data.id);
            }.bind(this)))
        }
    })
})();
