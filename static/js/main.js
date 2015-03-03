(function(){
    var m = new app.CoreModel();

    function render() {
        React.render(React.createElement(app.CoreApp, {
            model: m
        }), document.getElementById('example'));
    }

    m.subscribe(render);
    render();
})();
