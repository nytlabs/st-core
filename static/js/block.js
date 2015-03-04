var app = app || {};

(function() {
    app.RouteComponent = React.createClass({
        displayName: "RouteComponent",
        render: function() {
            var style = {}
            return React.createElement('g', {
                transform: 'translate(' + this.props.x + ', ' + this.props.y + ')',
            }, [
                React.createElement('text', {
                    className: "route_label unselectable",
                    textAnchor: this.props.left ? "start" : "end",
                    style: style
                }, this.props.model.name),
            ])
        },
        componentDidMount: function() {
            // console.log(this.getDOMNode().getComputedTextLength());
        }
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
                height: this.props.model.height + 5,
                className: classes,
                key: 'bg'
            }, null));
            children.push(React.createElement('text', {
                x: 0,
                y: 0,
                className: 'label unselectable',
                key: 'label'
            }, this.props.model.data.type));

            children = children.concat(this.props.model.data.inputs.map(function(e, i) {
                return React.createElement(app.RouteComponent, {
                    model: e,
                    geometry: this.props.model.inputs[i],
                    x: 0,
                    y: this.props.model.inputs[i].routeY,
                    left: true
                }, null)
            }.bind(this)));

            children = children.concat(this.props.model.data.outputs.map(function(e, i) {
                return React.createElement(app.RouteComponent, {
                    model: e,
                    geometry: this.props.model.outputs[i],
                    x: this.props.model.width,
                    y: this.props.model.outputs[i].routeY,
                    left: false,
                }, null)
            }.bind(this)));

            return React.createElement('g', {}, children);
        }
    })
})();
