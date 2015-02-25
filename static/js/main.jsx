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

    app.Group = function(data){
            for(var key in data){
                     this[key] = data[key]
            }
    }
   
    app.Block = function(data){
            for(var key in data){
                    this[key] = data[key]
            }
    }

    app.Source = function(data){
            for(var key in data){
                    this[key] = data[key];
            }
    }

    var nodes = {
        'block': app.Block,
        'source': app.Source,
        'group': app.Group
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
                if( nodes.hasOwnProperty(m.type) === true){
                        this.entities[m.data[m.type].id] = new nodes[m.type](m.data[m.type]);
            
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
        getInitialState: function(){
                return {
                        top: this.props.model.position.y,
                        left: this.props.model.position.x}
        },
        dragStart: function(e){
                e.dataTransfer.effectAllowed = "move";
                e.dataTransfer.setData("text/plain", JSON.stringify(this.props.model));
        },
        dragEnd: function(e){
                this.setState({left: e.pageX, top: e.pageY - e.nativeEvent.toElement.clientHeight})

                app.Utils.request(
                        "PUT", 
                        "blocks/" + this.props.model.id + "/position", 
                        {x: e.pageX, y: e.pageY - e.nativeEvent.toElement.clientHeight }, 
                        null
                );
        },
        render: function(){
                var entity = this.props.model;
                if(entity.hasOwnProperty('inputs')){
                return(
                        <div className="block" style={this.state} onDragStart={this.dragStart} draggable="true" onDragEnd={this.dragEnd}>
                                {entity.id}<br />
                                {entity.label}<br />
                                {entity.type}<br />
                                [{JSON.stringify(this.state)}]<br/>
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
    dragOver: function(e){
            e.preventDefault();
    },
    render: function() {
            var entities = this.props.model.entities; 
            return (
                    <div className="stage" onDragOver={this.dragOver}>
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
