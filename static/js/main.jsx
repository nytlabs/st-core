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
                for(var key in m.data[m.type]){
                    if(key !== 'id'){
                        this.entities[m.data[m.type].id][key] = m.data[m.type][key] 
                    }
                }
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
                var divStyle = {
                        top: entity.position.y,
                        left: entity.position.x,
                }
                if(entity.hasOwnProperty('inputs')){
                return(
                        <div className="block" style={divStyle}>
                                {entity.id}
                                {entity.label}
                                {entity.type}
                                [{entity.position.x},{entity.position.y}]
                                <ul>
                                        {entity.inputs.map(function(name,i){
                                                return <li key={i}>{name}</li>
                                        })}

                                </ul>
                                <ul>
                                        {entity.outputs.map(function(name,i){
                                                return <li key={i}>{name}</li>
                                        })}
                                </ul>

                        
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
