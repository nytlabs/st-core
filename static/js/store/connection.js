var app = app || {};

(function() {
    var connections = {};

    function Connection(data) {
        this.data = data;
        console.log(app.BlockStore.getBlock(this.data.from.id));

        app.BlockStore.getBlock(this.data.from.id).addListener(function() {
            this.render();
        }.bind(this));
        app.BlockStore.getBlock(this.data.to.id).addListener(function() {
            this.render()
        }.bind(this));
    }

    Connection.prototype = Object.create(app.Emitter.prototype);
    Connection.constructor = Connection;

    Connection.prototype.render = function() {
        console.log("WTF!!!");
        console.log("I NEED TO BE RENDERED! - a connectin: ", this);
    }

    function ConnectionStore() {}
    ConnectionStore.prototype = Object.create(app.Emitter.prototype);
    ConnectionStore.constructor = ConnectionStore;

    ConnectionStore.prototype.getConnection = function(id) {
        return connections[id];
    }

    var rs = new ConnectionStore();

    function createConnection(connection) {
        if (connections.hasOwnProperty(connection.id) === true) {
            console.warn('could not create connection:', connection.id, ' already exists');
            return
        }
        connections[connection.id] = new Connection(connection);
    }

    function deleteConnection(id) {
        if (connections.hasOwnProperty(id) === false) {
            console.warn('could not delete connections: ', id, ' does not exist');
            return
        }
        delete connections[id]
    }

    app.Dispatcher.register(function(event) {
        switch (event.action) {
            case app.Actions.WS_CONNECTION_CREATE:
                console.log(event.action);
                createConnection(event.data);
                rs.emit();
                break;
            case app.Actions.WS_CONNECTION_DELETE:
                console.log(event.action);
                deleteConnection(event.id);
                rs.emit();
                break;
        }
    })

    app.ConnectionStore = rs;
}())
