var app = app || {};

// TODO 
// This file desperately needs to be refactored. The portion of CoreApp that 
// is related to the stage, the background lines, and the selection of nodes
// can be put into its own component. 

(function() {
    'use strict';

    app.CoreApp = React.createClass({
        displayName: 'CoreApp',
        getInitialState: function() {
            return {
                autoCompleteVisible: false,
                autoCompleteX: null,
                autoCompleteY: null,
            }
        },
        _showAutoComplete: function(x, y, relativeX, relativeY) {
            this.setState({
                autoCompleteVisible: true,
                autoCompleteX: x,
                autoCompleteY: y,
                relativeX: relativeX,
                relativeY: relativeY,
            });
        },
        _hideAutoComplete: function() {
            this.setState({
                autoCompleteVisible: false,
            });
        },
        render: function() {
            var canvasGraph = React.createElement(app.CanvasGraphComponent, {
                key: 'canvas_graph',
                width: function() {},
                height: function() {},
                showAutoComplete: this._showAutoComplete,
                onClick: this._hideAutoComplete,
            }, null);

            var tools = React.createElement(app.ToolsComponent, {
                key: 'tool_list',
                onGroup: function() {},
                onUngroup: function() {},
            });

            var panelList = React.createElement(app.PanelListComponent, {
                nodes: [],
                key: 'panel_list',
            });

            var children = [canvasGraph, panelList, tools];

            if (this.state.autoCompleteVisible === true) {
                children.push(React.createElement(app.AutoCompleteComponent, {
                    key: 'autocomplete',
                    x: this.state.autoCompleteX,
                    y: this.state.autoCompleteY,
                    relativeX: this.state.relativeX,
                    relativeY: this.state.relativeY,
                    onEnter: this._hideAutoComplete,
                }));
            }

            var container = React.createElement('div', {
                className: 'app',
            }, children);

            return container
        }
    })
})();
