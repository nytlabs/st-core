var app = app || {};

/* BlockComponent and RouteComponent
 * BlockComponent renders both Block and Group models. (Group is just an
 * extended version of Block). BlockComponent creates RouteComponents to
 * render each individual route. Routes can almost work entirely
 * independantly of Blocks/Groups, as they are self-contained in state
 * and the nodes that they reference. The only bit of data that they depend
 * on from their parent node (block or group) is the geometry information
 * that is required for laying out in that particular blockor group (width,
 * height # of total routes, etc).
 */

(function() {
    'use strict';

    app.RouteComponent = React.createClass({
        displayName: 'RouteComponent',
        getInitialState: function() {
            var route = app.RouteStore.getRoute(this.props.id);
            return {
                blocked: route.blocked,
                active: route.active,
            }
        },
        componentDidMount: function() {
            app.RouteStore.getRoute(this.props.id).addListener(this._onChange);
        },
        componentWillUnmount: function() {
            app.RouteStore.getRoute(this.props.id).removeListener(this._onChange);
        },
        shouldComponentUpdate: function(props, state) {
            var u = (state.blocked != this.state.blocked || state.active != this.state.active)
            return u
        },
        _onChange: function() {
            var route = app.RouteStore.getRoute(this.props.id);
            this.setState({
                blocked: route.blocked,
                active: route.active,
            })
        },
        /*handleMouseUp: function() {
            this.props.onChange(this.props.model)
        },*/
        render: function() {
            var children = [];

            children.push(React.createElement('text', {
                className: 'route_label unselectable',
                textAnchor: this.props.direction === 'input' ? 'start' : 'end',
                key: 'route_label',
            }, this.props.displayName))

            var waiting = this.state.blocked ? ' waiting' : ' open';
            var circleClass = 'route_circle' + ' ' + this.props.displayType + waiting;
            var cx = this.props.geometry.routeRadius * (this.props.direction === 'input' ? -.5 : .5);
            var cy = this.props.geometry.routeRadius * -.5;

            children.push(React.createElement('circle', {
                onMouseUp: this.handleMouseUp,
                cx: cx,
                cy: cy,
                r: this.props.geometry.routeRadius,
                className: circleClass,
                key: 'route_circle',
            }, null))

            // if this route has a value set for it OR we are connected to
            // something, then fill the route.
            if (this.state.active || this.props.connections.length !== 0) {
                children.push(React.createElement('circle', {
                    onMouseUp: this.handleMouseUp,
                    cx: cx,
                    cy: cy,
                    r: this.props.geometry.routeRadius - 2,
                    className: this.props.connections.length !== 0 ? 'route_circle_filled' : 'route_circle_white',
                    key: 'route_circle_filled'
                }, null))
            }

            return React.createElement('g', {
                transform: 'translate(' +
                    (this.props.direction === 'input' ? 0 : this.props.geometry.width) + ', ' +
                    ((1 + this.props.displayIndex) * this.props.geometry.routeHeight) + ')',
            }, children)
        },
    })
})();

(function() {
    app.CrankComponent = React.createClass({
        getInitialState: function() {
            return {
                status: null,
            }
        },
        componentDidMount: function() {
            app.BlockStore.getBlock(this.props.id).crank.addListener(this._onChange);
        },
        componentWillUnmount: function() {
            app.BlockStore.getBlock(this.props.id).crank.removeListener(this._onChange);
        },
        shouldComponentUpdate: function(props, state) {
            return this.status != state.status
        },
        _onChange: function() {
            this.setState({
                status: app.BlockStore.getBlock(this.props.id).crank.status,
            })
        },
        render: function() {
            var children = [
                React.createElement('circle', {
                    cx: 0,
                    cy: 0,
                    r: 5.0,
                    className: 'tick_circle ' + (this.state.status === 'kernel' ? 'kernel' : ''),
                    key: 'tick_bg'
                }, null),
                React.createElement('circle', {
                    cx: 5,
                    cy: 0,
                    r: 3.0,
                    key: 'tick',
                    fill: 'red',
                    className: 'ticker_0 ' + (this.state.status === 'running' ? 'running' : ''),
                }, null),
            ]
            return React.createElement('g', {
                transform: 'translate(' + this.props.x + ', ' + this.props.y + ')',
            }, children)
        }
    })
})();

(function() {
    'use strict';

    app.BlockComponent = React.createClass({
        displayName: 'BlockComponent',
        onChange: function(r) {
            this.props.onRouteEvent(r)
        },
        shouldComponentUpdate: function(props, state) {
            return props.selected != this.props.selected
        },
        render: function() {
            var block = app.BlockStore.getBlock(this.props.id);

            var classes = 'block';
            if (this.props.selected === true) classes += ' selected';
            var children = [];
            children.push(React.createElement('rect', {
                x: 0,
                y: 0,
                width: block.geometry.width,
                height: block.geometry.height,
                className: classes,
                key: 'bg'
            }, null));

            // TODO: style this better
            var title = block.data.type;
            title += ' ' + block.data.label;

            children.push(React.createElement('text', {
                x: 0,
                y: 0,
                className: 'type unselectable',
                key: 'type'
            }, title));

            children.push(React.createElement(app.CrankComponent, {
                x: block.geometry.width * .5,
                y: block.geometry.height + 10,
                id: block.data.id,
                key: 'crank'
            }));

            // add the input routes to the block
            children = children.concat(block.inputs.map(function(routeName, i) {
                var route = app.RouteStore.getRoute(routeName);
                return React.createElement(app.RouteComponent, {
                    onChange: this.onChange,
                    id: routeName,
                    geometry: block.geometry,
                    direction: 'input',
                    displayIndex: i,
                    displayName: route.data.name,
                    displayType: route.data.type,
                    connections: route.connections,
                    key: routeName
                })
            }.bind(this)));

            // add the output routes to the block
            children = children.concat(block.outputs.map(function(routeName, i) {
                var route = app.RouteStore.getRoute(routeName);
                return React.createElement(app.RouteComponent, {
                    onChange: this.onChange,
                    id: routeName,
                    direction: 'output',
                    geometry: block.geometry,
                    displayIndex: i,
                    displayName: route.data.name,
                    displayType: route.data.type,
                    connections: route.connections,
                    key: routeName
                })
            }.bind(this)));

            return React.createElement('g', {}, children);
        }
    })
})();
