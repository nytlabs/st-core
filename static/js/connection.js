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
                fill: 'transparent'
            }
            if (this.props.selected === true) lineStyle.stroke = "blue";
            return (
                React.createElement("path", {
                    style: lineStyle,
                    d: this.props.model.path,
                    onMouseUp: this.onMouseUp,
                })
            )
        }
    })
})();
