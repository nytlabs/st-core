(function() {
    var types = [];

    function get(nodeType) {
        app.Utils.request(
            'GET',
            nodeType + '/library', {},
            function(e) {
                // TODO: fix the stupid schema in the server so we don't have to do this.
                types = types.concat(JSON.parse(e.response).map(function(t) {
                    return {
                        name: t.name,
                        type: nodeType,
                        source: t.source,
                    }
                }));
            })
    }

    function Library() {
        get('blocks');
        get('sources');
    }

    Library.prototype = Object.create(app.Emitter.prototype);
    Library.constructor = Library;

    Library.prototype.getLibrary = function() {
        return types;
    }

    var library = new Library();

    app.Dispatcher.register(function(event) {
        // TODO: in the case of live updates to a library that need to be
        // propagated to views, this is where they should live.
        switch (event.action) {}
    });

    app.LibraryStore = library;
})();
