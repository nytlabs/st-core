var app = app || {};

(function() {
    var Emitter = function() {}

    Emitter.prototype.addListener = function(listener) {
        if (!this._events) this._events = [];
        this._events.push(listener);
        return this;
    }

    Emitter.prototype.removeListener = function(listener) {
        if (!this._events) return this;
        var index = this._events.indexOf(listener);
        if (index < 0) return this;
        this._events.splice(index, 1);
        return this;
    }

    Emitter.prototype.emit = function() {
        if (!this._events) return this;
        this._events.forEach(function(listener) {
            listener();
        })
        return this;
    }

    app.Emitter = Emitter;
})();
