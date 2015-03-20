var app = app || {};

(function() {
    var lineStyle = {
        stroke: "black",
        strokeWidth: 2,
        fill: 'none'
    }

    app.ConnectionToolComponent = React.createClass({
        displayName: 'ConnectionToolComponent',
        componentWillMount: function() {
            document.addEventListener('mousemove', this.handleMouseMove);
        },
        getInitialState: function() {
            return {
                x: 0,
                y: 0
            }
        },
        handleMouseMove: function(e) {
            this.setState({
                x: e.pageX,
                y: e.pageY,
            })
        },
        render: function() {
            var x1, y1, cx1, cy1, cx2, cy2, x2, y2;
            if (this.props.from != null) {
                var from = this.props.from;
                x1 = from.width + from.data.position.x + from.outputs[this.props.route].routeCircleX;
                y1 = from.data.position.y + from.outputs[this.props.route].routeY + from.outputs[this.props.route].routeCircleY;
                cx1 = from.width + 50 + from.data.position.x;
                cy1 = from.data.position.y + from.outputs[this.props.route].routeY + from.outputs[this.props.route].routeCircleY;
                x2 = this.state.x;
                y2 = this.state.y;
                cx2 = -50 + this.state.x;
                cy2 = this.state.y;
            }

            if (this.props.to != null) {
                var to = this.props.to;
                x2 = to.data.position.x + to.inputs[this.props.route].routeX + to.inputs[this.props.route].routeCircleX;
                y2 = to.data.position.y + to.inputs[this.props.route].routeY + to.inputs[this.props.route].routeCircleY;
                cx2 = -50 + x2;
                cy2 = y2;
                x1 = this.state.x;
                y1 = this.state.y;
                cx1 = x1 + 50;
                cy1 = this.state.y;
            }

            return React.createElement('path', {
                style: lineStyle,
                d: ['M', x1, ' ', y1, ' C ', cx1, ' ', cy1, ' ', cx2, ' ', cy2, ' ', x2, ' ', y2].join(''),
            }, null)
        }
    })

    app.ConnectionComponent = React.createClass({
        displayName: "ConnectionComponent",
        onMouseUp: function(e) {
            this.props.nodeSelect(this.props.model.data.id);
        },
        render: function() {
            if (this.props.selected === true) lineStyle.stroke = "blue";
            return React.createElement("path", {
                style: lineStyle,
                d: this.props.model.path,
                onMouseUp: this.onMouseUp,
            }, null)
        }
    })
})();
