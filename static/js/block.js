var app = app || {};

(function(){
    app.BlockComponent = React.createClass({
        displayName: "BlockComponent",
        render: function() {
            var classes = "block";
            if (this.props.selected === true) classes += " selected";

            var children = [];
            children.push(React.createElement('rect', {
                x: 0,
                y: 0,
                width: 50,
                height: 20,
                className: classes,
                key: 'bg'
            }, null));
            children.push(React.createElement('text', {
                x: 0,
                y: 10,
                className: 'label unselectable',
                key: 'label'
            }, this.props.model.type));
            return React.createElement('g', {}, children);
        }
    })
})();