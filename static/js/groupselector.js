var app = app || {};

/* 
 * GroupSelectorComponent
 * Displays menu that lists all groups.
 * Clicking a group updates CoreModel with the currently in-view group.
 * Changes to CoreModel update which group is selected, and which nods are
 * currently in view.
 */

(function() {
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
                if (this.props.focusedGroup === g.id) classes += ' focused';
                return React.createElement('div', {
                    className: classes,
                    key: g.id,
                    onClick: this.onClick.bind(null, g),
                }, g.label + ' ' + g.id);
            }.bind(this)))
        }
    })
})();
