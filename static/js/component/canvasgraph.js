var app = app || {};

(function() {




    app.CanvasGraphComponent = React.createClass({
        displayName: "canvas",
        getInitialState: function() {
            return {
                blocks: app.BlockStore.getBlocks()
            }
        },
        shouldComponentUpdate: function() {
            return false;
        },
        componentDidMount: function() {
            app.BlockStore.addListener(this._onChange);
        },
        componentWillUnmount: function() {
            app.BlockStore.removeListener(this._onChange);
        },
        _onChange: function() {
            var ctx = React.findDOMNode(this.refs.test).getContext('2d');
            ctx.clearRect(0, 0, this.props.width, this.props.height);
            app.BlockStore.getBlocks().forEach(function(id) {
                var block = app.BlockStore.getBlock(id);
                ctx.fillRect(block.position.x, block.position.y, 10, 10);
            })
        },
        render: function() {
            return React.createElement('canvas', {
                ref: 'test',
                width: this.props.width,
                height: this.props.height
            }, null);
        }
    });
})();
