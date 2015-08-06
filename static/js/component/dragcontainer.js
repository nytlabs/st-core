var app = app || {};

/* DragContainerComponent
 * This component wraps around an element of model Entity and provides
 * mouse-drag functionality.
 */

(function() {
    'use strict';

    app.DragContainer = React.createClass({
        displayName: 'DragContainer',
        getInitialState: function() {
            var block = app.BlockStore.getBlock(this.props.id);
            return {
                x: block.position.x,
                y: block.position.y,
                dragging: false
            }
        },
        componentDidMount: function() {
            app.BlockStore.getBlock(this.props.id).addListener(this._onChange);
        },
        componentWillUnmount: function() {
            app.BlockStore.getBlock(this.props.id).removeListener(this._onChange);
        },
        shouldComponentUpdate: function(props, state) {
            return false;
        },
        _onChange: function() {
            var position = app.BlockStore.getBlock(this.props.id).position;
            this.setState({
                x: position.x,
                y: position.y,
            }, function() {
                React.findDOMNode(this.refs.container).setAttribute('transform', 'translate(' + this.state.x + ', ' + this.state.y + ')');
            })
        },
        onMouseDown: function(e) {
            this.props.nodeSelect(this.props.model.data.id);
            this.setState({
                dragging: true,
            });
            document.addEventListener('mousemove', this.onMouseMove)
            document.addEventListener('mouseup', this.onMouseUp)
        },
        onMouseUp: function(e) {
            document.removeEventListener('mousemove', this.onMouseMove)
            document.removeEventListener('mouseup', this.onMouseUp)
            this.setState({
                dragging: false,
            });
            this.props.onDragStop();
        },
        onMouseMove: function(e) {
            if (this.state.dragging) {
                this.props.onDrag(e.movementX, e.movementY);
            }
        },
        render: function() {
            return (
                React.createElement('g', {
                        ref: 'container',
                        transform: 'translate(' + this.state.x + ', ' + this.state.y + ')',
                        onMouseDown: this.onMouseDown,
                        onMouseUp: this.onMouseUp,
                    },
                    this.props.children
                )
            )

        }
    })
})();
