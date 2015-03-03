var app = app || {};

(function(){
    app.DragContainer = React.createClass({
        displayName: "DragContainer",
        getInitialState: function() {
            return {
                dragging: false,
                offX: null,
                offY: null,
                debounce: 0,
            }
        },
        onMouseDown: function(e) {
            //m.select(this.props.model.id);
            this.props.nodeSelect(this.props.model.id);

            this.setState({
                dragging: true,
                offX: e.pageX - this.props.x,
                offY: e.pageY - this.props.y
            })
        },
        componentDidUpdate: function(props, state) {
            if (this.state.dragging && !state.dragging) {
                document.addEventListener('mousemove', this.onMouseMove)
                document.addEventListener('mouseup', this.onMouseUp)
            } else if (!this.state.dragging && state.dragging) {
                document.removeEventListener('mousemove', this.onMouseMove)
                document.removeEventListener('mouseup', this.onMouseUp)
            }
        },
        onMouseUp: function(e) {
            this.props.model.setPosition({
                x: e.pageX - this.state.offX,
                y: e.pageY - this.state.offY
            })

            this.setState({
                dragging: false,
            })
        },
        onMouseMove: function(e) {
            if (this.state.dragging) {
                this.props.model.setPosition({
                    x: e.pageX - this.state.offX,
                    y: e.pageY - this.state.offY
                })
            }
        },
        render: function() {
            return (
                React.createElement("g", {
                        transform: 'translate(' + this.props.x + ', ' + this.props.y + ')',
                        onMouseMove: this.onMouseMove,
                        onMouseDown: this.onMouseDown,
                        onMouseUp: this.onMouseUp,
                    },
                    this.props.children
                )
            )

        }
    })
})();