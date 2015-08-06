var app = app || {};

(function() {
    var connections = {};

    function Connection(data) {
        this.data = data;

        app.BlockStore.getBlock(this.data.from.id).addListener(function() {
            this.render();
        }.bind(this));
        app.BlockStore.getBlock(this.data.to.id).addListener(function() {
            this.render()
        }.bind(this));

        this.render();
    }

    Connection.prototype = Object.create(app.Emitter.prototype);
    Connection.constructor = Connection;

    Connection.prototype.render = function() {
        // TODO: instead of blocks, this should somehow find the top-most visible geometry that
        // the route is apart of (for groups);
        var from = app.BlockStore.getBlock(this.data.from.id);
        var to = app.BlockStore.getBlock(this.data.to.id);

        var routeIdFrom = this.data.from.id + '_' + this.data.from.route + '_output';
        var routeIdTo = this.data.to.id + '_' + this.data.to.route + '_input';

        var routeIndexFrom = from.outputs.indexOf(routeIdFrom);
        var routeIndexTo = to.inputs.indexOf(routeIdTo);

        var yFrom = from.geometry.routeHeight * routeIndexFrom + (from.geometry.routeRadius * .5);
        var xFrom = from.geometry.routeRadius * .5 + from.geometry.width;

        var yTo = to.geometry.routeHeight * routeIndexTo - (to.geometry.routeRadius * .5);
        var xTo = to.geometry.routeRadius * -.5 + 0;

        xFrom += from.position.x;
        yFrom += from.position.y;
        xTo += to.position.x;
        yTo += to.position.y;

        var c = [xFrom, yFrom, xFrom + 50, yFrom, xTo - 50, yTo, xTo, yTo];
        this.curve = [
            'M',
            c[0], ' ',
            c[1], ' C ',
            c[2], ' ',
            c[3], ' ',
            c[4], ' ',
            c[5], ' ',
            c[6], ' ',
            c[7]
        ].join('');

        this.emit();
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
