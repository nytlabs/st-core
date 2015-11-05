var app = app || {};

/* PanelListComponent
 * The right sidebar that is contains a series of panels.
 */

(function() {
    'use strict';

    app.PanelListComponent = React.createClass({
        displayName: 'PanelListComponent',
        getInitialState: function() {
            return {
                ids: app.SelectionStore.getIdsByKind(app.Node),
            }
        },
        componentDidMount: function() {
            app.SelectionStore.addListener(this._onUpdate);
        },
        componentWillUnmount: function() {
            app.SelectionStore.removeListener(this._onUpdate);
        },
        _onUpdate: function() {
            this.setState({
                ids: app.SelectionStore.getIdsByKind(app.Node),
            });
        },
        render: function() {
            var children = [
                React.createElement(app.GroupTreeComponent, {})
            ];

            children = children.concat(this.state.ids.map(function(id) {
                return React.createElement(app.RoutesPanelComponent, {
                    key: id,
                    id: id
                })
            }));

            return React.createElement('div', {
                className: 'panel_list',
                children
            })
        },
    })
})();
