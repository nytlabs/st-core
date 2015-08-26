var app = app || {};

/* PanelComponent & PanelEditableComponent
 * Produces a list of fields that are the current representation of input
 * values for blocks/groups that are sent to the component.
 *
 * TODO: fix the {'data': ...} nonsense
 */

(function() {
    app.RoutePanelInput = React.createClass({
        getInitialState: function() {
            return {
                name: '',
                type: '',
                value: '',
            }
        },
        componentDidMount: function() {
            app.RouteStore.getRoute(this.props.id).addListener(this._update);
            this._update();
        },
        componentWillUnmount: function() {
            app.RouteStore.getRoute(this.props.id).removeListener(this._update);
        },
        _update: function() {
            var route = app.RouteStore.getRoute(this.props.id);
            var value = '';
            if (route.data.value !== null) {
                value = JSON.stringify(route.data.value.data);
            }

            this.setState({
                name: route.data.name,
                type: route.data.type,
                value: value,
            })
        },
        _handleChange: function(event) {
            this.setState({
                value: event.target.value
            });
        },
        _onKeyDown: function(event) {
            if (event.keyCode !== 13) return;

            var value = null;
            if (this.state.value !== null) {
                try {
                    value = {
                        data: JSON.parse(this.state.value)
                    }
                } catch (e) {
                    return
                }
            }

            app.Dispatcher.dispatch({
                action: app.Actions.APP_REQUEST_ROUTE_UPDATE,
                id: this.props.id,
                value: value
            })
        },
        render: function() {
            return React.createElement('div', {}, [
                React.createElement('div', {
                    className: 'label',
                    key: 'label',
                }, this.state.name),
                React.createElement('input', {
                    type: 'text',
                    ref: 'value',
                    key: 'value',
                    value: this.state.value,
                    onChange: this._handleChange,
                    onKeyDown: this._onKeyDown,
                }, null)
            ]);
        }
    });

})();

(function() {
    app.RoutesPanelComponent = React.createClass({
        displayName: 'PanelComponent',
        render: function() {
            var block = app.NodeStore.getNode(this.props.id);

            var children = [
                React.createElement('div', {
                    key: 'block_header',
                    className: 'block_header',
                }, block.data.type),
            ];

            children = children.concat(block.inputs.map(function(r) {
                return React.createElement(app.RoutePanelInput, {
                    id: r.id,
                    key: r.id,
                }, null)
            }));

            return React.createElement('div', {
                className: 'panel'
            }, children);
        }
    })
})();
