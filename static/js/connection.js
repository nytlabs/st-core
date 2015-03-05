var app = app || {};

(function() {
    app.ConnectionComponent = React.createClass({
        displayName: "ConnectionComponent",
        onMouseUp: function(e) {
            this.props.nodeSelect(this.props.model.data.id);
        },
        render: function() {
            var lineStyle = {
                stroke: "black",
                strokeWidth: 2,
                fill: 'none'
            }
            if (this.props.selected === true) lineStyle.stroke = "blue";
            return React.createElement('g', {},
                React.createElement("path", {
                    style: lineStyle,
                    d: this.props.model.path,
                    onMouseUp: this.onMouseUp,
                }, null),
                React.createElement('circle', {
                    cx: this.props.model.to.x,
                    cy: this.props.model.to.y,
                    r: this.props.model.routeRadius,
                    className: "route_circle_filled",
                    key: "route_circle_filled_to"
                }, null),
                React.createElement('circle', {
                    cx: this.props.model.from.x,
                    cy: this.props.model.from.y,
                    r: this.props.model.routeRadius,
                    className: "route_circle_filled",
                    key: "route_circle_filled_from"
                }, null)
            )
        }
    })
})();
