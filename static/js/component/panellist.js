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
              if (n instanceof app.Source){
                return React.createElement(app.ParametersPanelComponent, { model: n, key: n.data.id }, null)
              } else {
                return React.createElement(app.RoutesPanelComponent, { model: n, key: n.data.id }, null)
              }
            }))
        },
    })
})();
