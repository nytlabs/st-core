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

    var ops = {
        create: function(obj, type) {
            entities[obj.id] = {
                data: obj,
                type: type
            };
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
                        entities[updated.id][key] = updated[key]
                    }
                }
            }
        }
    }

    Object.observe(entities, function(e) {
        switch (e[0].type) {
            case 'add':
                nodeEnter(e[0].object);
                break;
            case 'delete':
                // ??? who knows
        }
    });

    function nodeEnter()


})();
