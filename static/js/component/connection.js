var app = app || {};

/* ConnectionComponent
 *
 * TODO:
 * props.node no longer needed for this component, replaced with parentNode
 */

(function() {
    'use strict';

    function getCoords(node, route) {
        var cx = node.geometry.routeRadius * (route.direction === 'input' ? -.5 : .5);
        var cy = node.geometry.routeRadius * -.5;
        var routeX = (route.direction === 'input' ? 0 : node.geometry.width) +
            cx + node.data.position.x;
        var routeY = (1 + route.displayIndex) * node.geometry.routeHeight +
            cy + node.data.position.y;
        return {
            x: routeX,
            y: routeY
        }
    }

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
                stroke: 'black',
                strokeWidth: 2,
                fill: 'none'
            }

            var coord = getCoords(this.props.connecting.parentNode, this.props.connecting);
            coord.x += this.props.translateX;
            coord.y += this.props.translateY;

            // if the tool is enabled but the mouse has not moved, set null
            // state as route position
            var target = {
                x: this.state.x === null ? coord.x : this.state.x,
                y: this.state.y === null ? coord.y : this.state.y,
            }

            var c = [
                coord.x, coord.y, coord.x, coord.y,
                target.x, target.y, target.x, target.y
            ];

            if (this.props.connecting.direction === 'output') {
                c[2] += 50.0;
                c[4] -= 50.0;
            } else {
                c[4] += 50.0;
                c[2] -= 50.0;
            }

            return React.createElement('path', {
                style: lineStyle,
                strokeDasharray: '5,5',
                d: [
                    'M',
                    c[0], ' ',
                    c[1], ' C ',
                    c[2], ' ',
                    c[3], ' ',
                    c[4], ' ',
                    c[5], ' ',
                    c[6], ' ',
                    c[7]
                ].join(''),
            }, null)
        }
    })

    app.ConnectionComponent = React.createClass({
        displayName: 'ConnectionComponent',
        onMouseUp: function(e) {
            this.props.nodeSelect(this.props.model.data.id);
        },
        render: function() {
            var lineStyle = {
                stroke: 'black',
                strokeWidth: 2,
                fill: 'none'
            }

            var from = getCoords(this.props.model.from.node, this.props.model.from.route)
            var to = getCoords(this.props.model.to.node, this.props.model.to.route)

            var c = [from.x, from.y, from.x, from.y, to.x, to.y, to.x, to.y];
            c[2] += 50.0;
            c[4] -= 50.0;

            if (this.props.selected === true) lineStyle.stroke = 'blue';
            return React.createElement('path', {
                style: lineStyle,
                d: [
                    'M',
                    c[0], ' ',
                    c[1], ' C ',
                    c[2], ' ',
                    c[3], ' ',
                    c[4], ' ',
                    c[5], ' ',
                    c[6], ' ',
                    c[7]
                ].join(''),
                onMouseUp: this.onMouseUp,
            }, null)
        }
    })
})();
