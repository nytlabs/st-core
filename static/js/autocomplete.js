var app = app || {};

(function() {
    app.AutoCompleteComponent = React.createClass({
        displayName: "AutoCompleteComponent",
        getInitialState: function() {
            return {
                text: "",
                options: [],
            }
        },
        componentWillMount: function() {
            this.autocomplete(this.state.options);
        },
        autocomplete: function(s) {
            var reBegin = new RegExp('^' + s, 'i');
            var reIn = new RegExp(s, 'i');

            var beginList = this.props.options.filter(function(o) {
                return o.type.match(reBegin)
            }).sort(function(a, b) {
                if (a.type > b.type) return 1;
                if (a.type < b.type) return -1;
                return 0
            });

            var inList = this.props.options.filter(function(o) {
                return o.type.match(reIn);
            }).sort(function(a, b) {
                if (a.type > b.type) return 1;
                if (a.type < b.type) return -1;
                return 0;
            })

            this.setState({
                options: beginList.concat(inList.filter(function(o) {
                    return beginList.indexOf(o) === -1;
                }))
            });
        },
        handleChange: function(e) {
            this.autocomplete(e.target.value);
        },
        render: function() {
            var input = React.createElement('input', {
                onChange: this.handleChange,
                key: 'autocomplete_input'
            }, null);

            var options = this.state.options.map(function(o) {
                return React.createElement('li', {
                    key: o.type
                }, o.type);
            });

            var list = React.createElement('ul', {
                key: 'autocomplete_list',
            }, options);

            return React.createElement('div', {
                style: {
                    width: '100',
                    top: this.props.y,
                    left: this.props.x,
                    position: 'absolute',
                    zIndex: 100,
                },
            }, [input, list]);
        }
    });
})();
