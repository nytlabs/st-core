(function() {
    var entities = {};

    var ws = new WebSocket("ws://localhost:7071/updates");
    ws.onopen = function(e) {
        ws.send("list");
    }
    ws.onmessage = function(e) {
        var update = JSON.parse(e.data);
        ops[update.action](update.data[update.type], update.type);
    }

    var update = {
        route: function() {},
        param: function() {},
        child: function() {},
    }

    Object.observe(entities, function(e) {
        e[0].object[e[0].name][e[0].type]();
    });

    window.Test = function() {
        entities[0].data.foo = +new Date();
        console.log(entities);
    }

    function Entity(obj, type) {
        var _this = this;
        this.data = obj;
        this.type = type;
        Object.observe(this.data, function(e) {
            _this.render();
        })
    }

    Entity.prototype.add = function() {
        console.log("Hello, noob", this)
    }

    Entity.prototype.update = function() {
        console.log("say what?");
    }

    Entity.prototype.delete = function() {
        console.log("lol kcya", this.type)
    }

    Entity.prototype.render = function() {

    }

    /*function Block(obj, type) {
        Entity.call(this, obj, type)
    }

    Block.prototype = Entity.prototype;

    function Connection(obj, type) {
        Entity.call(this, obj, type);
    }

    Connection.prototype = Entity.prototype;

    function Group(obj, type) {
        Entity.call(this, obj, type);

    }

    Group.prototype = Entity.prototype;

    function Source(obj, type) {
        Entity.call(this, obj, type);

    }

    Source.prototype = Entity.prototype;

    function Link(obj, type) {
        Entity.call(this, obj, type);

    }

    Link.prototype = Entity.prototype;

    var types = {
        'block': Block,
        'group': Group,
        'connection': Connection,
        'source': Source,
        'link': Link,
    };*/

    var ops = {
        create: function(obj, type) {
            //entities[obj.id] = new types[type](obj, type);
            entities[obj.id] = new Entity(obj, type);
        },
        delete: function(obj) {
            delete entities[obj.id];
        },
        update: function(obj) {
            if (update.hasOwnProperty(obj.type)) {
                update[obj.type](obj)
            } else {
                updated = obj;
                for (var key in updated) {
                    if (key != "id") {
                        entities[updated.id].data[key] = updated[key]
                    }
                }
            }
        }
    }

})();
