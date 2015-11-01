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
            case app.Actions.APP_SELECT_TOGGLE:
                selectToggle(event.ids);
                rs.emit();
                break;
            case app.Actions.APP_DESELECT_ALL:
                deselectAll();
                rs.emit();
                break;
        }
    })

    app.SelectionStore = rs;
}())
