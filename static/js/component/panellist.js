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
                ids: app.BlockStore.getSelected(),
            }
        },
        componentDidMount: function() {
            app.BlockSelection.addListener(this._onUpdate);
        },
        componentWillUnmount: function() {
            app.BlockSelection.addListener(this._onUpdate);
        },
        _onUpdate: function() {
            console.log("hello friends")
            this.setState({
                ids: app.BlockStore.getSelected(),
            })
        },
        render: function() {
            return React.createElement('div', {
                className: 'panel_list'
            }, this.state.ids.map(function(id) {
                console.log(id);
                return React.createElement(app.RoutesPanelComponent, {
                        key: id,
                        id: id
                    })
                    /*if (n instanceof app.Source) {
                        return React.createElement(app.ParametersPanelComponent, {
                            model: n,
                            key: id
                        }, null)
                    } else {
                        return React.createElement(app.RoutesPanelComponent, {
                            model: n,
                            key: n.data.id
                        }, null)
                    }*/
            }))
        },
    })
})();
