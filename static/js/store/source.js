var app = app || {};

(function() {
    var sources = {};

    function SourceStore() {}
    SourceStore.prototype = Object.create(app.Emitter.prototype);
    SourceStore.constructor = SourceStore;

    SourceStore.prototype.getSource = function(id) {
        return sources[id];
    }

    var rs = new SourceStore();

    function createSource(source) {
        if (sources.hasOwnProperty(source.id) === true) {
            console.warn('could not create source:', source.id, ' already exists');
            return
        }
        sources[source.id] = source;
    }

    function deleteSource(id) {
        if (sources.hasOwnProperty(id) === false) {
            console.warn('could not delete source: ', id, ' does not exist');
            return
        }
        delete sources[id]
    }

    function updateSource(source) {
        if (sources.hasOwnProperty(source.id) === false) {
            console.warn('could not update source: ', source.id, ' does not exist');
            return
        }
        source[source.id] = source;
    }

    app.Dispatcher.register(function(event) {
        switch (event.action) {
            case app.Actions.WS_SOURCE_CREATE:
                console.log(event.action);
                createSource(event.data);
                rs.emit();
                break;
            case app.Actions.WS_SOURCE_DELETE:
                console.log(event.action);
                deleteSource(action.id);
                rs.emit();
                break;
            case app.Actions.WS_SOURCE_UPDATE:
                console.log(event.action);
                rs.emit();
                break;
            case app.Actions.WS_SOURCE_UPDATE_PARAMS:
                console.log(event.action);
                rs.emit();
                break;
        }
    })

    app.SourceStore = rs;
}())
