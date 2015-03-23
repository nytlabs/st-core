var app = app || {};

(function() {
    app.RouteComponent = React.createClass({
        displayName: "RouteComponent",
        handleMouseUp: function() {
            this.props.onChange(this.props.model)
        },
        render: function() {
            var children = [];

            children.push(
                React.createElement('text', {
                    className: "route_label unselectable",
                    textAnchor: this.props.model.routeAlign,
                    key: "route_label",
                }, this.props.model.data.name)
            )

            var circleClass = "route_circle" + " " + this.props.model.data.type;
            children.push(
                React.createElement('circle', {
                    onMouseUp: this.handleMouseUp,
                    cx: this.props.model.routeCircleX,
                    cy: this.props.model.routeCircleY,
                    r: this.props.model.routeRadius,
                    className: circleClass,
                    key: "route_circle",
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
                        cx: this.props.model.routeCircleX,
                        cy: this.props.model.routeCircleY,
                        r: this.props.model.routeRadius - 2,
                        className: "route_circle_filled",
                        key: "route_circle_filled"
                    }, null)
                )
            }

            return React.createElement('g', {
                transform: 'translate(' +
                    this.props.model.routeX + ', ' +
                    this.props.model.routeY + ')',
            }, children)
        },
    })
})();

(function() {
    app.BlockComponent = React.createClass({
        displayName: "BlockComponent",
        onChange: function(r) {
            this.props.onRouteEvent({
                id: this.props.model.data.id,
                route: r.routeIndex,
                direction: r.routeDirection
            })
        },
        render: function() {
            var classes = "block";
            if (this.props.selected === true) classes += " selected";
            var children = [];
            children.push(React.createElement('rect', {
                x: 0,
                y: 0,
                width: this.props.model.width,
                height: this.props.model.height,
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

            var routes = this.props.model.inputs.concat(this.props.model.outputs);

            children = children.concat(routes.map(function(r, i) {
                return React.createElement(app.RouteComponent, {
                    onChange: this.onChange,
                    model: r,
                    key: r.routeAlign + i
                }, null)
            }.bind(this)));

            return React.createElement('g', {}, children);
        }
    })
})();
