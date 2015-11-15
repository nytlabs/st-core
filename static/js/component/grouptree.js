var app = app || {};

/* TreeComponent
 * Recursive component for displaying group hierarchy.
 */
(function() {
    app.TreeComponent = React.createClass({
        displayName: 'TreeComponent',
        _onMouseDown: function(e) {
            app.NodeStore.setRoot(this.props.tree.id);
        },
        render: function() {
            var style = '';
            if (app.NodeStore.getRoot() === this.props.tree.id) {
                style += 'current-group'
            }

            var label = this.props.tree.id;
            var node = app.NodeStore.getNode(this.props.tree.id);
            if (node.data.label.length > 0) {
                label = node.data.label;
            }

            var children = [
                React.createElement('span', {
                    onMouseDown: this._onMouseDown,
                    className: style,
                }, label)
            ]
            if (this.props.tree.children.length !== 0) {
                var list = this.props.tree.children.map(function(child) {
                    return React.createElement('li', {},
                        React.createElement(app.TreeComponent, {
                            key: child.id,
                            tree: child
                        }, null));
                })
                children.push(React.createElement('ul', {}, list));
            }
            return React.createElement('div', {
                className: 'group_tree'
            }, children);
        }
    })
})();

/* GroupTreeComponent
 * Sidebar widget for displaying group hierarchy, selection/moving between
 * groups.
 */
(function() {
    app.GroupTreeComponent = React.createClass({
        displayName: 'GroupTreeComponent',
        getInitialState: function() {
            return {
                tree: null
            }
        },
        componentDidMount: function() {
            app.NodeStore.addListener(this._update);
            this._update();
        },
        componentWillUnmount: function() {
            app.NodeStore.removeListener(this._update);
        },
        _update: function() {
            this.setState({
                tree: app.NodeStore.getTree()
            })
        },
        render: function() {
            var children = [
                React.createElement('div', {
                    key: 'block_header',
                    className: 'block_header',
                }, "groups"),
            ];

            if (this.state.tree !== null) {
                children.push(
                    React.createElement(app.TreeComponent, {
                        key: 'tree',
                        tree: this.state.tree,
                    }, null)
                )
            }

            return React.createElement('div', {
                className: 'panel unselectable'
            }, children);
        }
    })
})();
