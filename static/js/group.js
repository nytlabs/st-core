var app = app || {};

(function() {
    app.GroupComponent = React.createClass({
        displayName: "GroupComponent",
        render: function() {
            var classes = "block"
            if (this.props.selected === true) classes += " selected";

            var children = [];
            children.push(React.createElement('rect', {
                x: 0,
                y: 0,
                width: 50,
                height: 50,
                className: classes,
                key: 'bg'
            }, null));
            children.push(React.createElement('text', {
                x: 0,
                y: 10,
                className: 'label unselectable',
                key: 'label'
            }, 'group ' + this.props.model.id));
            return React.createElement('g', {}, children);
        }
    })
})();
