var app = app || {};

(function() {
    var _callbacks = [];
    var Dispatcher = function() {}
    Dispatcher.prototype.register = function(callback) {
        _callbacks.push(callback);
    }
    Dispatcher.prototype.dispatch = function(payload) {
        console.log(payload);
        _callbacks.forEach(function(callback) {
            callback(payload);
        })
    }
    app.Dispatcher = new Dispatcher();
})();
