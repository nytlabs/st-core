var app = app || {};

(function() {
    var groups = {};

    function GroupStore() {}
    GroupStore.prototype = Object.create(app.Emitter.prototype);
    GroupStore.constructor = GroupStore;

    GroupStore.prototype.getGroup = function(id) {
        return groups[id];
    }

    var rs = new GroupStore();

    function createGroup(group) {
        if (groups.hasOwnProperty(group.id) === true) {
            console.warn('could not create group:', group.id, ' already exists');
            return
        }
        groups[group.id] = group;
    }

    function deleteGroup(id) {
        if (groups.hasOwnProperty(id) === false) {
            console.warn('could not delete group: ', id, ' does not exist');
            return
        }
        delete groups[id]
    }

    function updateGroup(group) {
        if (groups.hasOwnProperty(group.id) === false) {
            console.warn('could not update group: ', group.id, ' does not exist');
            return
        }
        group[group.id] = group;
    }

    app.Dispatcher.register(function(event) {
        switch (event.action) {
            case app.Actions.WS_GROUP_CREATE:
                console.log(event.action);
                createGroup(event.data);
                rs.emit();
                break;
            case app.Actions.WS_GROUP_UPDATE:
                console.log(event.action);
                deleteGroup(action.id);
                rs.emit();
                break;
            case app.Actions.WS_GROUP_DELETE:
                console.log(event.action);
                rs.emit();
                break;
            case app.Actions.WS_GROUP_ADD_CHILD:
                console.log(event.action);
                rs.emit();
                break;
            case app.Actions.WS_GROUP_REMOVE_CHILD:
                console.log(event.action);
                rs.emit();
                break;
        }
    })

    app.GroupStore = rs;
}())
