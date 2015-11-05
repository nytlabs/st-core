var app = app || {};


(function() {
    app.TreeComponent = React.createClass({
        displayName: 'TreeComponent',
        _onMouseDown: function(e) {
            app.NodeStore.setRoot(this.props.tree.id);
        },
        render: function() {
            var children = [
                React.createElement('span', {
                    onMouseDown: this._onMouseDown,
                }, this.props.tree.id),
            ]
            if (this.props.tree.children.length !== 0) {
                var list = this.props.tree.children.map(function(child) {
                    return React.createElement('li', {},
                        React.createElement(app.TreeComponent, {
                            tree: child
                        }, null));
                })
                children.push(React.createElement('ul', {}, list));
            }
            return React.createElement('div', {}, children);
        }
    })
})();

(function() {
    app.GroupTreeComponent = React.createClass({
        displayName: 'GroupTreeComponent',
        componentDidMount: function() {
            app.NodeStore.addListener(this._update);
        },
        componentWillUnmount: function() {
            app.NodeStore.removeListener(this._update);
        },
        _update: function() {
            this.render();
        },
        render: function() {
            var children = [
                React.createElement('div', {
                    key: 'block_header',
                    className: 'block_header',
                }, "groups"),
                React.createElement(app.TreeComponent, {
                    tree: app.NodeStore.getTree(),
                }, null)
            ];

            return React.createElement('div', {
                className: 'panel'
            }, children);
        }
    })
})();
