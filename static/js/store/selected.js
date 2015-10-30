var app = app || {};

(function() {
    var selected = [];

    function SelectionStore() {}
    SelectionStore.prototype = Object.create(app.Emitter.prototype);
    SelectionStore.constructor = SelectionStore;

    SelectionStore.prototype.getSelected = function() {
        // return a copy, not the original
        return selected.slice();
    }

    SelectionStore.prototype.isSelected = function(id) {
        return selected.indexOf(id) !== -1 ? true : false;
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
                    return i != id;
                });
            }

            app.Dispatcher.dispatch({
                action: app.Actions.APP_RENDER,
                id: id
            });
        })
    }

    function deselectAll() {
        var toRender = selected.slice();
        selected = [];

        toRender.forEach(function(id) {
            app.Dispatcher.dispatch({
                action: app.Actions.APP_RENDER,
                id: id
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
