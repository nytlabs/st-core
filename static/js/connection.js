var app = app || {};

(function() {

    app.ConnectionToolComponent = React.createClass({
        displayName: 'ConnectionToolComponent',
        componentWillMount: function() {
            window.addEventListener('mousemove', this.handleMouseMove);
        },
        componentWillUnmount: function() {
            window.removeEventListener('mousemove', this.handleMouseMove);
        },
        getInitialState: function() {
            return {
                x: null,
                y: null
            }
        },
        handleMouseMove: function(e) {
            this.setState({
                x: e.pageX,
                y: e.pageY,
            })
        },
        render: function() {
            var lineStyle = {
                stroke: "black",
                strokeWidth: 2,
                fill: 'none'
            }

            var node = this.props.node;
            var routing = this.props.connecting;

            var routeX = node.data.position.x +
                node[routing.direction + 's'][routing.route].routeX +
                node[routing.direction + 's'][routing.route].routeCircleX +
                this.props.translateX;

            var routeY = node.data.position.y +
                node[routing.direction + 's'][routing.route].routeY +
                node[routing.direction + 's'][routing.route].routeCircleY +
                this.props.translateY;

            // if the tool is enabled but the mouse has not moved, set null
            // state as route position
            var target = {
                x: this.state.x === null ? routeX : this.state.x,
                y: this.state.y === null ? routeY : this.state.y,
            }

            var c = [routeX, routeY, routeX, routeY, target.x, target.y, target.x, target.y];

            if (routing.direction === 'output') {
                c[2] += 50.0;
                c[4] -= 50.0;
            } else {
                c[4] += 50.0;
                c[2] -= 50.0;
            }

            return React.createElement('path', {
                style: lineStyle,
                strokeDasharray: "5,5",
                d: ['M', c[0], ' ', c[1], ' C ', c[2], ' ', c[3], ' ', c[4], ' ', c[5], ' ', c[6], ' ', c[7]].join(''),
            }, null)
        }
    })

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
            return React.createElement("path", {
                style: lineStyle,
                d: this.props.model.path,
                onMouseUp: this.onMouseUp,
            }, null)
        }
    })
})();
