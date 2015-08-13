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
            return {}
        },
        render: function() {
            var canvasGraph = React.createElement(app.CanvasGraphComponent, {
                key: 'canvas_graph',
                width: function() {},
                height: function() {},
                doubleClick: function() {}
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

            var children = [canvasGraph, /*groupList,*/ panelList, tools];

            var container = React.createElement('div', {
                className: 'app',
            }, children);

            return container
        }
    })
})();
