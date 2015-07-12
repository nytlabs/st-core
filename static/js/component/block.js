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
            var m = this.props.model;
            return {
                blocked: app.RouteStore.getRoute(m.id + '_' + m.index + '_' + m.direction).data.blocked
            }
        },
        componentDidMount: function() {
            var m = this.props.model;
            var id = m.id + '_' + m.index + '_' + m.direction;
            app.RouteStore.getRoute(id).addListener(this._onChange);
        },
        componentWillUnmount: function() {
            var m = this.props.model;
            var id = m.id + '_' + m.index + '_' + m.direction;
            app.RouteStore.getRoute(id).removeListener(this._onChange);
        },
        shouldComponentUpdate: function(props, state) {
            return state.blocked != this.state.blocked
        },
        _onChange: function() {
            var m = this.props.model;
            var id = m.id + '_' + m.index + '_' + m.direction;
            this.setState({
                blocked: app.RouteStore.getRoute(id).blocked
            })
        },
        handleMouseUp: function() {
            this.props.onChange(this.props.model)
        },
        render: function() {
            var children = [];
            var direction = this.props.model.direction;

            children.push(
                React.createElement('text', {
                    className: 'route_label unselectable',
                    textAnchor: direction === 'input' ? 'start' : 'end',
                    key: 'route_label',
                }, this.props.model.data.name)
            )

            var waiting = this.state.blocked ? ' waiting' : '';
            var circleClass = 'route_circle' + ' ' + this.props.model.data.type + waiting;
            var cx = this.props.geometry.routeRadius * (direction === 'input' ? -.5 : .5);
            var cy = this.props.geometry.routeRadius * -.5;

            children.push(
                React.createElement('circle', {
                    onMouseUp: this.handleMouseUp,
                    cx: cx,
                    cy: cy,
                    r: this.props.geometry.routeRadius,
                    className: circleClass,
                    key: 'route_circle',
                }, null)
            )

            // if this route has a value set for it OR we are connected to
            // something, then fill the route.
            if ((this.props.model.data.hasOwnProperty('value') &&
                    this.props.model.data.value !== null) ||
                this.props.model.connections.length !== 0) {
                children.push(
                    React.createElement('circle', {
                        onMouseUp: this.handleMouseUp,
                        cx: cx,
                        cy: cy,
                        r: this.props.geometry.routeRadius - 2,
                        className: this.props.model.connections.length !== 0 ? 'route_circle_filled' : 'route_circle_white',
                        key: 'route_circle_filled'
                    }, null)
                )
            }

            return React.createElement('g', {
                transform: 'translate(' +
                    (direction === 'input' ? 0 : this.props.geometry.width) + ', ' +
                    ((1 + this.props.model.displayIndex) * this.props.geometry.routeHeight) + ')',
            }, children)
        },
    })
})();

(function() {
    app.CrankComponent = React.createClass({
        getInitialState: function() {
            return {
                lastCrank: null,
                tick: 0,
            }
        },
        componentWillReceiveProps: function(props) {
            if (props.lastCrank !== this.state.lastCrank) {
                var nextTick = this.state.tick < 7 ? this.state.tick + .5 : 0;
                this.setState({
                    tick: nextTick,
                    lastCrank: props.lastCrank,
                })
            }
        },
        render: function() {
            var state = '';
            if (this.state.lastCrank !== null && this.state.lastCrank !== undefined) {
                state = this.state.lastCrank.type;
            }
            var children = [
                React.createElement('circle', {
                    cx: 0,
                    cy: 0,
                    r: 5.0,
                    className: 'tick_circle' + (state === 'kernel' ? ' ' + state : ''),
                    key: 'tick_bg'
                }, null),
                React.createElement('circle', {
                    cx: 5,
                    cy: 0,
                    r: 3.0,
                    key: 'tick',
                    fill: 'red',
                    className: (state === 'running' ? state : 'ticker_' + Math.floor(this.state.tick)),
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
        render: function() {
            var classes = 'block';
            if (this.props.selected === true) classes += ' selected';
            var children = [];
            children.push(React.createElement('rect', {
                x: 0,
                y: 0,
                width: this.props.model.geometry.width,
                height: this.props.model.geometry.height,
                className: classes,
                key: 'bg'
            }, null));

            // TODO: style this better
            var title = this.props.model.data.type;
            title += ' ' + this.props.model.data.label;

            children.push(React.createElement('text', {
                x: 0,
                y: 0,
                className: 'type unselectable',
                key: 'type'
            }, title));

            children.push(React.createElement(app.CrankComponent, {
                x: this.props.model.geometry.width * .5,
                y: this.props.model.geometry.height + 10,
                lastCrank: this.props.model.lastCrank,
                key: 'crank'
            }));

            children = children.concat(this.props.model.routes.map(function(r, i) {
                return React.createElement(app.RouteComponent, {
                    onChange: this.onChange,
                    model: r,
                    geometry: r.parentNode.geometry, //this.props.model.geometry,
                    key: i
                })
            }.bind(this)));

            return React.createElement('g', {}, children);
        }
    })
})();
