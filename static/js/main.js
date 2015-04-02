(function() {
    'use strict';

    var m = new app.CoreModel();

    function render() {
        React.render(React.createElement(app.CoreApp, {
            model: m
        }), document.getElementById('st'));
    }

    m.subscribe(render);
    render();
})();
