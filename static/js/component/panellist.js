var app = app || {};

/* PanelListComponent
 * The right sidebar that is contains a series of panels.
 */

(function() {
    'use strict';

    app.PanelListComponent = React.createClass({
        displayName: 'PanelListComponent',
        render: function() {
            return React.createElement('div', {
                className: 'panel_list'
            }, this.props.nodes.filter(function(r) {
                return r instanceof app.Entity
            }).map(function(n) {
                return React.createElement(app.PanelComponent, {
                    model: n,
                    key: n.data.id
                }, null)
            }))
        },
    })
})();
