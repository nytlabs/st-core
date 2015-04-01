var app = app || {};

(function() {
    'use strict';

    var dm = new app.Utils.DebounceManager();

    app.Entity = function() {
        this.isDragging = false;
        this.parentNode = null;
    }

    app.Entity.prototype.setPosition = function(p) {
        this.data.position.x = p.x;
        this.data.position.y = p.y;

        // this function refreshes all connection geometry in view
        // it may be better to have a specific call for just connections that
        // are touching this particular entity.
        //this.model.refreshFocusedEdgeGeometry();
        this.model.inform()
        dm.push(this.id, function() {
            app.Utils.request(
                'PUT',
                this.instance() + 's/' + this.data.id + '/position', // would be nice to change API to not have the 'S' in it!
                p,
                null
            );
        }.bind(this), 50)
    }

    app.Entity.prototype.setDragging = function(e) {
        this.isDragging = e;
    }

})();
