var app = app || {};

(function() {
    app.ClipboardComponent = React.createClass({
        componentDidMount: function() {
            this.focus();
        },
        componentDidUpdate: function(props) {
            this.focus();
        },
        focus: function() {
            if (this.props.focus === true) {
                React.findDOMNode(this.refs.textarea).focus();
                React.findDOMNode(this.refs.textarea).select();
            } else {
                React.findDOMNode(this.refs.textarea).blur();
            }
        },
        stringifyProps: function() {
            var pattern = {
                blocks: [],
                sources: [],
                links: [],
                connections: [],
                groups: []
            };

            console.log(this.props.selected);

            pattern.blocks = this.props.selected.filter(function(n) {
                return n instanceof app.Block
            }).map(function(o) {
                return o.data;
            });

            pattern.sources = this.props.selected.filter(function(s) {
                return s instanceof app.Source
            }).map(function(o) {
                return o.data;
            });

            pattern.links = this.props.selected.filter(function(l) {
                return l instanceof app.Link
            }).map(function(o) {
                return o.data;
            });

            pattern.connections = this.props.selected.filter(function(c) {
                return c instanceof app.Connection
            }).map(function(o) {
                return o.data;
            });

            pattern.groups = this.props.selected.filter(function(c) {
                return c instanceof app.Group
            }).map(function(o) {
                return o.data;
            });

            return pattern;
        },
        shouldComponentUpdate: function(props) {
            React.findDOMNode(this.refs.textarea).value = JSON.stringify(this.stringifyProps());
            return true;
        },
        onPaste: function(e) {
            //console.log(React.findDOMNode(this.refs.textarea).value);
            app.Utils.request(
                'post',
                'groups/' + this.props.group + '/import',
                JSON.parse(React.findDOMNode(this.refs.textarea).value),
                function(e) {
                    console.log(e);
                });
        },
        render: function() {
            var pattern = this.stringifyProps();
            return React.createElement('textarea', {
                ref: 'textarea',
                className: 'clipboard',
                onChange: this.onPaste,
                defaultValue: JSON.stringify(pattern)
            }, null);
        }
    });
})();
