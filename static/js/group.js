var app = app || {};

(function(){
    app.GroupComponent = React.createClass({
        displayName: "GroupComponent",
        render: function() {
            return (
                React.createElement("rect", {
                    className: "block",
                    x: "0",
                    y: "0",
                    width: "100",
                    height: "10"
                })
            )
        }
    })
})();