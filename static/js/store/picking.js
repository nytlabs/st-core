var app = app || {};

(function() {
    var colorToNode = {};
    var colorIdx = 1;
    var recycleStack = [];

    function PickingStore() {}

    PickingStore.prototype = Object.create(app.Emitter.prototype);
    PickingStore.constructor = PickingStore;

    PickingStore.prototype.getColor = function(node) {
        var c;
        if (recycleStack.length === 0) {
            var ret = [];
            // via http://stackoverflow.com/a/15804183
            if (colorIdx < 16777215) {
                ret.push(colorIdx & 0xff); // R
                ret.push((colorIdx & 0xff00) >> 8); // G 
                ret.push((colorIdx & 0xff0000) >> 16); // B
                colorIdx += 1;
            }
            c = "rgb(" + ret.join(',') + ")";
        } else {
            c = recycleStack.pop();
        }

        colorToNode[c] = node;
        return c;
    }

    PickingStore.prototype.removeColor = function(color) {
        delete colorToNode[color];
        recycleStack.push(color);
    }

    PickingStore.prototype.colorToNode = function(color) {
        if (!colorToNode.hasOwnProperty(color)) {
            return null
        }

        return colorToNode[color];
    }

    var rs = new PickingStore();

    app.Dispatcher.register(function(event) {
        switch (event.action) {
            //...
        }
    })

    app.PickingStore = rs;
}())
