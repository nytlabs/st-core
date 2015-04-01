/* StageComponent
 * - renders the background and grid
 * - provides selection box and selection events
 * - group's translation state is updated in model
 *
 *           +-----+ group               +---------------+
 *           |     | width               |               |
 *           | app | height              | stage         |
 *           |     | +-----------------> | (no children) |
 *           |     |                     |               |
 *   +-----> |     | onSelectionChange() |               |
 *   |       |     | OnDoubleClick()     |               |
 *   |       |     | onMouseUp()         |               |
 *   |       |     | <-----------------+ |               |
 *   |       +-----+                     +---------------+
 *   |
 *   |       +-----+                     +
 *   |       |model| group.setTranslation|
 *   +-----+ |     | <-------------------+
 *           +-----+
 */

var app = app || {};

(function() {

    var GRID_PX = 200.0;

    app.StageComponent = React.createClass({
        displayName: 'StageComponent',
        getInitialState: function() {
            return {
                dragging: false,
                offX: null, // initial mouse position when dragging starts
                offY: null,
                selectionRect: {
                    x1: null, // used for calculation
                    y1: null,
                    x: null, // used for rendering
                    y: null,
                    width: null,
                    height: null,
                    enabled: false
                }
            }
        },
        componentWillMount: function() {
            window.addEventListener('mousemove', this.handleMouseMove);
            window.addEventListener('mouseup', this.handleMouseUp);
        },
        componentWillUnmount: function() {
            window.removeEventListener('mousemove', this.handleMouseMove);
            window.removeEventListener('mouseup', this.handleMouseUp);
        },
        handleMouseMove: function(e) {
            if (this.state.dragging) {
                this.props.group.setTranslation(e.pageX - this.state.offX, e.pageY - this.state.offY);
            }

            if (this.state.selectionRect.enabled === true) {
                // ensure that we don't have a negative width for the rect
                var x1 = this.state.selectionRect.x1;
                var y1 = this.state.selectionRect.y1;
                var x2 = e.pageX;
                var y2 = e.pageY;
                var rectX = x2 - x1 < 0 ? x2 : x1;
                var rectY = y2 - y1 < 0 ? y2 : y1;
                var width = Math.abs(x2 - x1);
                var height = Math.abs(y2 - y1);
                var translateX = this.props.group.translateX;
                var translateY = this.props.group.translateY;
                this.setState({
                    selectionRect: {
                        x1: x1,
                        y1: y1,
                        enabled: true,
                        x: rectX,
                        y: rectY,
                        width: width,
                        height: height,
                    }
                })

                // update parent components about what we're selecting
                this.props.onSelectionChange(rectX, rectY, width, height);
            }

        },
        handleMouseUp: function(e) {
            if (this.state.selectionRect.enabled === true) {
                this.setState({
                    selectionRect: {
                        enabled: false
                    }
                })
            }

            if (this.state.dragging === true) {
                this.handleMouseMove(e);
                this.setState({
                    dragging: false
                })
            }
        },
        handleMouseDown: function(e) {
            if (e.nativeEvent.button === 0) {
                this.setState({
                    selectionRect: {
                        x1: e.pageX,
                        y1: e.pageY,
                        enabled: true
                    },
                });
            }

            // TODO: make proper cancel selection event
            // like onSelectionCancel
            this.props.onSelectionChange(0, 0, 0, 0);

            if (e.nativeEvent.button === 2) {
                this.setState({
                    dragging: true,
                    offX: e.pageX - this.props.group.translateX,
                    offY: e.pageY - this.props.group.translateY,
                });
            }
        },
        render: function() {
            var nodes = [];

            // the background rect
            var background = React.createElement('rect', {
                className: 'background',
                x: '0',
                y: '0',
                width: this.props.width,
                height: this.props.height,
                onMouseDown: this.handleMouseDown,
                onDoubleClick: this.handleDoubleClick,
                key: 'background'
            }, null);

            // the grid
            var translateX = this.props.group.translateX;
            var translateY = this.props.group.translateY;
            var x = translateX % GRID_PX;
            var y = translateY % GRID_PX;
            var lines = [];
            var hMax = Math.floor(this.props.width / GRID_PX);
            var vMax = Math.floor(this.props.height / GRID_PX);
            for (var i = 0; i <= hMax; i++) {
                lines.push(React.createElement('line', {
                    key: 'h' + i,
                    x1: x + (i * GRID_PX),
                    y1: 0,
                    x2: x + (i * GRID_PX),
                    y2: this.props.height,
                    stroke: 'rgba(220,220,220,1)'
                }, null));
            }

            for (var i = 0; i <= vMax; i++) {
                lines.push(React.createElement('line', {
                    key: 'v' + i,
                    x1: 0,
                    y1: y + (i * GRID_PX),
                    x2: this.props.width,
                    y2: y + (i * GRID_PX),
                    stroke: 'rgba(220,220,220,1)'
                }, null));
            }

            var lineGroup = React.createElement('g', {
                key: 'line_group',
            }, lines)

            // the origin point
            var origin = React.createElement('circle', {
                cx: translateX,
                cy: translateY,
                r: 5,
                fill: 'rgba(220,220,220,1)',
                key: 'origin',
            }, null);

            var nodes = [background, lineGroup, origin];

            // selection rect
            if (this.state.selectionRect.enabled === true) {
                var selectionRect = React.createElement('rect', {
                    x: this.state.selectionRect.x,
                    y: this.state.selectionRect.y,
                    width: this.state.selectionRect.width,
                    height: this.state.selectionRect.height,
                    className: 'selection_rect',
                    key: 'selection_rect'
                }, null);

                nodes.push(selectionRect);
            }

            return React.createElement('g', {
                onDoubleClick: this.props.onDoubleClick,
                onMouseUp: this.props.onMouseUp
            }, nodes)
        }
    })
})();
