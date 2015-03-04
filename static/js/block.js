var app = app || {};

(function() {
    app.RouteComponent = React.createClass({
        displayName: "RouteComponent",
        render: function() {
            return React.createElement('g', {
                transform: 'translate(' + this.props.model.routeX + ', ' + this.props.model.routeY + ')',
            }, [
                React.createElement('text', {
                    className: "route_label unselectable",
                    textAnchor: this.props.model.routeAlign,
                    key: "route_label",
                }, this.props.model.data.name),
                React.createElement('circle', {
                    cx: this.props.model.routeCircleX,
                    cy: this.props.model.routeCircleY,
                    r: this.props.model.routeRadius,
                    className: "route_circle",
                    key: "route_circle"
                }, null)
            ])
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
            children.push(React.createElement('text', {
                x: 0,
                y: 0,
                className: 'label unselectable',
                key: 'label'
            }, this.props.model.data.type));

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
