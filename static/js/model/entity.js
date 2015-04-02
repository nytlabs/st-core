var app = app || {};

(function() {
    'use strict';

    app.Entity = function() {
        this.isDragging = false;
        this.parentNode = null;
    }

    app.Entity.prototype.setPosition = function(p) {
        this.data.position.x = p.x;
        this.data.position.y = p.y;

        this.model.inform()
    }

    app.Entity.prototype.postPosition = function() {
        app.Utils.request(
            'PUT',
            this.instance() + 's/' + this.data.id + '/position', // would be nice to change API to not have the 'S' in it!
            this.data.position,
            null
        );
    }

    app.Entity.prototype.setDragging = function(e) {
        this.isDragging = e;
    }

})();
