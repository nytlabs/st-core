var app = app || {};

(function() {
    'use strict';

    app.CoreModel = function() {
        this.entities = {};
        this.onChanges = [];

        var ws = new WebSocket("ws://localhost:7071/updates");

        ws.onmessage = function(m) {
            this.update(JSON.parse(m.data));
        }.bind(this)

        ws.onopen = function() {
            ws.send('list');
        }
    }

    app.CoreModel.prototype.subscribe = function(onChange) {
        this.onChanges.push(onChange);
    }

    app.CoreModel.prototype.inform = function() {
        console.log(this.entities);
        this.onChanges.forEach(function(cb) {
            cb();
        });
    }

    app.CoreModel.prototype.update = function(m) {
        switch (m.action) {
            case 'update':
                break;
            case 'create':
                if( ['block','group','source'].indexOf(m.type) != -1 ){
                        this.entities[m.data[m.type].id] = m.data[m.type];
                }
                break;
            case 'delete':
                delete this.entities[m.data[m.type].id];
                break;
        }

        this.inform();
    }
})();

var m = new app.CoreModel();

var Entity = React.createClass({
        render: function(){
                var entity = this.props.model;
                if(entity.hasOwnProperty('inputs')){
                return(
                        <div>
                                {entity.id}
                                {entity.label}
                                {entity.type}
                                <ol>
                                        {entity.inputs.map(function(name,i){
                                                return <li key={i}>{name}</li>
                                        })}

                                </ol>

                        
                        </div>
                )
                } else {
                return <div>LOL WHO CARES</div>
                }
        }
})

var CoreApp = React.createClass({
    getInitialState: function() {
        return {
            group: 0,
        }
    },
    render: function() {
            var entities = this.props.model.entities; 
            return (
                    <div>
                            {Object.keys(entities).map(function(id){
                                    return <Entity key={id} model={entities[id]}/>
                             })}
                    </div>
            )
    }
})

function render() {
        React.render(<CoreApp model={m}/> , document.getElementById('example'));
}

m.subscribe(render);
render();
