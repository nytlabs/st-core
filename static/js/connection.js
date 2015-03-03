var app = app || {};

(function() {
    app.ConnectionComponent = React.createClass({
        displayName: "ConnectionComponent",
        render: function() {
            var from = this.props.graph.entities[this.props.model.data.from.id].data
            var to = this.props.graph.entities[this.props.model.data.to.id].data
            var lineStyle = {
                stroke: "black",
                strokeWidth: 2,
                fill: 'transparent'
            }
            var path = 'M' + (50 + from.position.x) + ' ' + from.position.y + ' C ';
            path += (from.position.x + 100) + ' ' + from.position.y + ', '
            path += (to.position.x - 50) + ' ' + to.position.y + ', '
            path += to.position.x + ' ' + to.position.y;

            return (
                React.createElement("path", {
                    style: lineStyle,
                    d: path
                })
            )
        }
    })
})();
