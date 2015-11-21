var app = app || {};

/* AutoCompleteComponenet
 * Given a list of items and a position by parent.
 * Fires onChange event to parent node when item is selected.
 */

// TODO: somehow get current group's translation and add that to the X and Y 
// upon creation of a block!

(function() {
    'use strict';

    app.AutoCompleteComponent = React.createClass({
        displayName: 'AutoCompleteComponent',
        getInitialState: function() {
            var options = app.LibraryStore.getLibrary();
            return {
                text: '',
                options: options,
                filtered: [],
            }
        },
        componentWillMount: function() {
            this._autocomplete(this.state.text);
        },
        componentDidMount: function() {
            React.findDOMNode(this.refs.auto_input).focus()
        },
        _autocomplete: function(s) {
            var reBegin = new RegExp('^' + app.Utils.escape(s), 'i');
            var reIn = new RegExp(app.Utils.escape(s), 'i');

            var beginList = this.state.options.filter(function(o) {
                return o.name.match(reBegin)
            }).sort(function(a, b) {
                if (a.name > b.name) return 1;
                if (a.name < b.name) return -1;
                return 0
            });

            var inList = this.state.options.filter(function(o) {
                return o.name.match(reIn);
            }).sort(function(a, b) {
                if (a.name > b.name) return 1;
                if (a.name < b.name) return -1;
                return 0;
            })

            this.setState({
                filtered: beginList.concat(inList.filter(function(o) {
                    return beginList.indexOf(o) === -1;
                }))
            });
        },
        _handleChange: function(e) {
            this._autocomplete(e.target.value);
        },
        _handleKeyUp: function(e) {
            if (e.nativeEvent.keyCode === 13) {
                var match = this.state.options.filter(function(o) {
                    return o.name === e.target.value
                });
                if (match.length === 1) {
                    this._onEnter(match[0]);
                }
                this.props.onEnter();
            }
        },
        _onEnter: function(match) {
            // TODO: move this to a api actions receiver.
            app.Utils.request(
                'POST',
                match.type, {
                    'type': match.name,
                    'parent': app.NodeStore.getRoot(),
                    'position': {
                        'x': this.props.relativeX,
                        'y': this.props.relativeY
                    }
                },
                null)
        },
        render: function() {
            var input = React.createElement('input', {
                onChange: this._handleChange,
                onKeyUp: this._handleKeyUp,
                key: '_autocomplete_input',
                ref: 'auto_input'
            }, null);

            var options = this.state.filtered.map(function(o) {
                return React.createElement('li', {
                    key: o.name
                }, o.name);
            });

            var list = React.createElement('ul', {
                key: '_autocomplete_list',
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
