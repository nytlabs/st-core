var app = app || {};

(function() {
    var selected = [];

    function SelectionStore() {}
    SelectionStore.prototype = Object.create(app.Emitter.prototype);
    SelectionStore.constructor = SelectionStore;

    SelectionStore.prototype.getSelected = function() {
        // return a copy, not the original
        return selected.slice().map(function(element) {
            return element.data.id
        });
    }

    SelectionStore.prototype.isSelected = function(id) {
        return selected.indexOf(id) !== -1 ? true : false;
    }

    SelectionStore.prototype.getIdsByKind = function(kind) {
        return selected.filter(function(item) {
            return item instanceof kind;
        }).map(function(element) {
            return element.data.id
        });
    }

    // returns a pattern from selected nodes to be imported/exported
    SelectionStore.prototype.getPattern = function() {
        //TODO: new pattern type should just consist of nodes and edges
        //this should be reflected in API as well. 

        var pattern = {
            blocks: [],
            sources: [],
            links: [],
            connections: [],
            groups: []
        };

        pattern.blocks = selected.filter(function(n) {
            return n instanceof app.Node && !(n instanceof app.Group) && !(n instanceof app.Source)
        }).map(function(o) {
            return o.data;
        });

        pattern.sources = selected.filter(function(s) {
            return s instanceof app.Source
        }).map(function(o) {
            return o.data;
        });

        pattern.links = selected.filter(function(l) {
            return l instanceof app.Link
        }).map(function(o) {
            return o.data;
        });

        pattern.connections = selected.filter(function(c) {
            return c instanceof app.Connection
        }).map(function(o) {
            return o.data;
        });

        function recurseData(id) {
            app.NodeStore.getNode(id).data.children.forEach(function(childId) {
                var node = app.NodeStore.getNode(childId);
                if (node instanceof app.Group) {
                    recurseData(childId);
                } else {
                    if (node instanceof app.Source) {
                        pattern.sources.push(node.data);
                    } else if (node instanceof app.Node) {
                        pattern.blocks.push(node.data);
                    }
                }
            });
        }

        pattern.groups = selected.filter(function(c) {
            return c instanceof app.Group
        }).map(function(o) {
            recurseData(o.data.id);
            return o.data;
        });

        return pattern;
    }

    var rs = new SelectionStore();

    function deselect(ids) {
        selected = selected.slice().filter(function(id) {
            return ids.indexOf(id) !== -1 ? false : true;
        });

        ids.forEach(function(id) {
            app.Dispatcher.dispatch({
                action: app.Actions.APP_RENDER,
                id: id
            })
        });
    }

    function selectToggle(ids) {
        ids.forEach(function(id) {
            if (selected.indexOf(id) === -1) {
                selected.push(id);
            } else {
                selected = selected.slice().filter(function(i) {
                    return i !== id;
                });
            }

            app.Dispatcher.dispatch({
                action: app.Actions.APP_RENDER,
                id: id.data.id
            });
        })
    }

    function deselectAll() {
        var toRender = selected.slice();
        selected = [];

        toRender.forEach(function(id) {
            app.Dispatcher.dispatch({
                action: app.Actions.APP_RENDER,
                id: id.data.id
            });
        });
    }

    function deleteSelected() {
        selected.forEach(function(s) {
            var type;
            if (s instanceof app.Node) type = 'blocks';
            if (s instanceof app.Group) type = 'groups';
            if (s instanceof app.Source) type = 'sources';
            if (s instanceof app.Connection) type = 'connections';
            if (s instanceof app.Link) type = 'links';

            app.Utils.request(
                'DELETE',
                type + '/' + s.data.id, {},
                null
            )
        });
    }

    app.Dispatcher.register(function(event) {
        switch (event.action) {
            case app.Actions.APP_DESELECT:
                deselect([event.id]);
                rs.emit();
                break;
            case app.Actions.APP_SELECT:
                deselectAll();
                selectToggle([event.id]);
                rs.emit();
                break;
            case app.Actions.APP_SELECT_ALL:
                deselectAll();
                selectToggle(event.ids);
                rs.emit();
                break;
            case app.Actions.APP_SELECT_TOGGLE:
                selectToggle(event.ids);
                rs.emit();
                break;
            case app.Actions.APP_DESELECT_ALL:
                deselectAll();
                rs.emit();
                break;
            case app.Actions.APP_DELETE_SELECTION:
                deleteSelected();
                rs.emit();
                break;
        }
    })

    app.SelectionStore = rs;
}())
