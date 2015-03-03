var app = app || {};

// TODO:
// create a standard model API that the rest of the components can use
// this standard API should use WS to communicate back to server

(function() {
    'use strict';

    app.CoreModel = function() {
        this.entities = {};
        this.list = [];
        this.groups = [];
        this.edges = [];
        this.onChanges = [];

        var ws = new WebSocket("ws://localhost:7071/updates");

        ws.onmessage = function(m) {
            this.update(JSON.parse(m.data));
        }.bind(this)

        ws.onopen = function() {
            ws.send('list');
        }
    }

    app.CoreModel.prototype.subscribe = function(onChange) {
        this.onChanges.push(onChange);
    }

    app.CoreModel.prototype.inform = function() {
        //console.log("updating model");
        this.onChanges.forEach(function(cb) {
            cb();
        });
    }

    app.Entity = function() {}

    function Debounce() {
        this.func = null;
        this.fire = null;
        this.last = null;
    }

    Debounce.prototype.push = function(e, duration) {
        if (this.last === null || this.last + duration < +new Date()) {
            this.last = +new Date();
            e();
            return;
        }
        this.func = e;
        if (this.fire != null) clearInterval(this.fire);
        this.fire = setTimeout(function() {
            this.func();
            this.last = +new Date()
        }.bind(this), duration);
    }

    function DebounceManager() {
        this.entities = {};
    }

    DebounceManager.prototype.push = function(id, f, duration) {
        if (!this.entities.hasOwnProperty(id)) {
            this.entities[id] = new Debounce();
        }
        this.entities[id].push(f, duration)
    }

    var dm = new DebounceManager();

    // TODO: put API methods on CoreModel
    app.Entity.prototype.setPosition = function(p) {
        this.position.x = p.x;
        this.position.y = p.y;
        this.__model.inform();
        dm.push(this.id, function() {
            app.Utils.request(
                "PUT",
                this.instance() + "s/" + this.id + "/position", // would be nice to change API to not have the "S" in it!
                p,
                null
            );
        }.bind(this), 50)
    }

    app.Group = function(data) {
        for (var key in data) {
            this[key] = data[key]
        }
    }

    app.Group.prototype = new app.Entity();

    app.Group.prototype.instance = function() {
        return "group";
    }

    app.Block = function(data) {
        for (var key in data) {
            this[key] = data[key]
        }
    }

    app.Block.prototype = new app.Entity();

    app.Block.prototype.instance = function() {
        return "block";
    }

    app.Source = function(data) {
        for (var key in data) {
            this[key] = data[key];
        }
    }

    app.Source.prototype = new app.Entity();

    app.Source.prototype.instance = function() {
        return "source";
    }

    app.Connection = function(data) {
        for (var key in data) {
            this[key] = data[key];
        }
    }

    app.Connection.prototype = new app.Entity();

    app.Connection.prototype.instance = function() {
        return "connection";
    }

    app.Link = function(data) {
        for (var key in data) {
            this[key] = data[key];
        }
    }

    app.Link.prototype = new app.Entity();

    app.Link.prototype.instance = function() {
        return "link";
    }

    var nodes = {
        'block': app.Block,
        'source': app.Source,
        'group': app.Group,
        'connection': app.Connection,
        'link': app.Link
    }

    // this takes an id and puts it at the very top of the list
    app.CoreModel.prototype.select = function(id) {
        this.list.push(this.list.splice(this.list.indexOf(this.entities[id]), 1)[0]);
        this.inform();
    }

    app.CoreModel.prototype.addChild = function(group, id) {
        this.entities[group].children.push(id);
        this.inform();
    }

    app.CoreModel.prototype.removeChild = function(group, id) {
        console.log(group, id, this.entities[group]);
        this.entities[group].children.splice(this.entities[group].children.indexOf(id), 1);
        this.inform();
    }

    app.CoreModel.prototype.update = function(m) {
        switch (m.action) {
            case 'update':
                for (var key in m.data[m.type]) {
                    if (key !== 'id') {
                        this.entities[m.data[m.type].id][key] = m.data[m.type][key]
                    }
                }
                break;
            case 'create':
                // create seperate action for child.
                if (m.type === "child") {
                    this.addChild(m.data.group.id, m.data.child.id);
                    return;
                }

                var n = new nodes[m.type](m.data[m.type]);
                // this reference allows all entities to inform() the model
                n.__model = this;
                this.entities[m.data[m.type].id] = n;
                this.list.push(this.entities[m.data[m.type].id]);

                if (m.type === "group") {
                    this.groups.push(n);
                }

                if (m.type === "connection" || m.type === "link") {
                    this.edges.push(n);
                }

                break;
            case 'delete':
                if (m.type === "child") {
                    this.removeChild(m.data.group.id, m.data.child.id); // this child nonsense is a mess
                    return
                }

                var i = this.list.indexOf(this.entities[m.data[m.type].id]);
                this.list.splice(i, 1);

                if (m.type === "group") {
                    var i = this.groups.indexOf(this.entities[m.data[m.type].id]);
                    this.groups.splice(i, 1);
                }

                if (m.type === "connection" || m.type == "link") {
                    var i = this.edges.indexOf(this.entities[m.data[m.type].id]);
                    this.edges.splice(i, 1);
                }

                delete this.entities[m.data[m.type].id];
                break;
        }

        this.inform();
    }
})();