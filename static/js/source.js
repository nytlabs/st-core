var app = app || {};

(function(){
    app.SourceComponent = React.createClass({
        displayName: "SourceComponent",
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