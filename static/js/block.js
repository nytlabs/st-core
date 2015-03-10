var app = app || {};

(function() {
    app.RouteComponent = React.createClass({
        displayName: "RouteComponent",
        render: function() {
            var children = [];

            children.push(
                React.createElement('text', {
                    className: "route_label unselectable",
                    textAnchor: this.props.model.routeAlign,
                    key: "route_label",
                }, this.props.model.data.name)
            )

            children.push(
                React.createElement('circle', {
                    cx: this.props.model.routeCircleX,
                    cy: this.props.model.routeCircleY,
                    r: this.props.model.routeRadius,
                    className: "route_circle",
                    key: "route_circle"
                }, null)
            )

            if ((this.props.model.data.hasOwnProperty('value') && this.props.model.data.value !== null) || this.props.model.connections.length !== 0) {
                children.push(
                    React.createElement('circle', {
                        cx: this.props.model.routeCircleX,
                        cy: this.props.model.routeCircleY,
                        r: this.props.model.routeRadius - 2,
                        className: "route_circle_filled",
                        key: "route_circle_filled"
                    }, null)
                )
            }

            return React.createElement('g', {
                transform: 'translate(' + this.props.model.routeX + ', ' + this.props.model.routeY + ')',
            }, children)
        },
    })
})();

(function() {
    app.BlockComponent = React.createClass({
        displayName: "BlockComponent",
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
                    model: r,
                    key: r.routeAlign + i
                }, null)
            }.bind(this)));

            return React.createElement('g', {}, children);
        }
    })
})();
