var app = app || {};

/* TODO: SOURCES */

(function() {
    'use strict';

    app.SourceComponent = React.createClass({
        displayName: 'SourceComponent',
        render: function() {
              console.log("hello!")
            return (
                React.createElement('circle', {
                    className: 'source',
                    x: '0',
                    y: '0',
                    width: '10',
                    height: '10'
                })
            )
        },
    })
})();
