var app = app || {};

/* TODO: SOURCES */

(function() {
    app.LinkComponent = React.createClass({
        displayName: "LinkComponent",
        render: function() {
            return (
                React.createElement("rect", {
                    className: "block",
                    x: "0",
                    y: "0",
                    width: "10",
                    height: "10"
                })
            )
        }
    })
})();
