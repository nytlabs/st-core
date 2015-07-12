var app = app || {};

(function() {
    var blocks = {};

    function BlockCollection() {}
    BlockCollection.prototype = Object.create(app.Emitter.prototype);
    BlockCollection.constructor = BlockCollection;

    BlockCollection.prototype.getBlock = function(id) {
        return blocks[id];
    }

    var rs = new BlockCollection();

    function createBlock(block) {
        if (blocks.hasOwnProperty(block.id) === true) {
            console.warn('could not create block:', block.id, ' already exists');
            return
        }
        blocks[block.id] = block;
    }

    function deleteBlock(id) {
        if (blocks.hasOwnProperty(id) === false) {
            console.warn('could not delete block: ', id, ' does not exist');
            return
        }
        delete blocks[id]
    }

    function updateBlock(block) {
        if (blocks.hasOwnProperty(block.id) === false) {
            console.warn('could not update block: ', block.id, ' does not exist');
            return
        }
        block[block.id] = block;
    }

    app.Dispatcher.register(function(event) {
        console.log(event);
        switch (event.action) {
            case app.Actions.WS_BLOCK_CREATE:
                console.log(event.action);
                createBlock(event.data);
                rs.emit();
                break;
            case app.Actions.WS_BLOCK_DELETE:
                console.log(event.action);
                deleteBlock(action.id);
                rs.emit();
                break;
            case app.Actions.WS_BLOCK_UPDATE:
                console.log(event.action);
                rs.emit();
                break;
            case app.Actions.WS_BLOCK_UPDATE_STATUS:
                console.log(event.action);
                rs.emit();
                break;
        }
    })

    app.BlockCollection = rs;
}())
