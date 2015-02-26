var app = app || {};

// TODO:
// create a standard model API that the rest of the components can use
// this standard API should use WS to communicate back to server

(function() {
    'use strict';

    app.CoreModel = function() {
        this.entities = {};
        this.list = [];
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
        console.log("updating model");
        this.onChanges.forEach(function(cb) {
            cb();
        });
    }

    app.Entity = function(){
    }

    app.Entity.prototype.setPosition = function(p){
                    app.Utils.request(
                            "PUT", 
                            this.instance() + "s/" + this.id + "/position", 
                            p, 
                            null
                    );   
    }

    app.Group = function(data){
            for(var key in data){
                     this[key] = data[key]
            }
    }

    app.Group.prototype = new app.Entity();

    app.Group.prototype.instance = function(){
        return "group";
    }
   
    app.Block = function(data){
            for(var key in data){
                    this[key] = data[key]
            }
    }

    app.Block.prototype = new app.Entity();
    
    app.Block.prototype.instance = function(){
        return "block";
    }
   

    app.Source = function(data){
            for(var key in data){
                    this[key] = data[key];
            }
    }

    app.Source.prototype = new app.Entity();

    app.Source.prototype.instance = function(){
        return "source";
    }

    var nodes = {
        'block': app.Block,
        'source': app.Source,
        'group': app.Group
    }

    // this takes an id and puts it at the very top of the list
    app.CoreModel.prototype.select = function(id){
        this.list.push(this.list.splice(this.list.indexOf(this.entities[id]), 1)[0]);
        this.inform();
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
                        this.list.push(this.entities[m.data[m.type].id]) 
                }
                break;
            case 'delete':
                var i = this.list.indexOf(this.entities[m.data[m.type].id]);
                this.list.splice(i, 1);
                delete this.entities[m.data[m.type].id];
                break;
        }

        this.inform();
    }
})();

var m = new app.CoreModel();

var DragContainer = React.createClass({
        getInitialState: function(){
                return {
                        dragging: false,
                        x: this.props.x,
                        y: this.props.y,
                        offX: null,
                        offY: null,
                        debounce: 0,
                }
        },
        onMouseDown: function(e){
                m.select(this.props.model.id);
                
                this.setState({
                        dragging: true,
                        offX: e.pageX - this.state.x,
                        offY: e.pageY - this.state.y
                }) 
                console.log(this.state.offX, this.state.offY);
        },
        componentDidUpdate: function (props, state) {
                if (this.state.dragging && !state.dragging) {
                        document.addEventListener('mousemove', this.onMouseMove)
                        document.addEventListener('mouseup', this.onMouseUp)
                } else if (!this.state.dragging && state.dragging) {
                        document.removeEventListener('mousemove', this.onMouseMove)
                        document.removeEventListener('mouseup', this.onMouseUp)
                }
        },
        onMouseUp: function(e){
                this.props.model.setPosition({x: this.state.x, y: this.state.y})
                
                this.setState({
                        dragging: false,
                })
        },
        onMouseMove: function(e){
                this.setState({
                        debounce: this.state.debounce + 1,
                })

                if(this.state.debounce > 5){
                        this.setState({
                                debounce: 0,
                        })
                        
                        this.props.model.setPosition({x: this.state.x, y: this.state.y})
                }
                
                if(this.state.dragging){
                        this.setState({
                                x: e.pageX - this.state.offX,
                                y: e.pageY - this.state.offY,
                        })
                }
        },
        componentWillReceiveProps: function(props){
                if(!this.state.dragging){
                        this.setState({
                                x: props.x,
                                y: props.y       
                        })
                }
        },
        render: function(){
                return (
                        <g 
                        transform={'translate(' + this.state.x + ', ' + this.state.y + ')'} 
                        onMouseMove={this.onMouseMove}
                        onMouseDown={this.onMouseDown}
                        onMouseUp={this.onMouseUp}
                        >
                        {this.props.children}
                        </g>
                )

        }
})

var Block = React.createClass({
        render: function(){
                return (
                        <rect className='block' x='0' y='0' width='100' height='100' />
                )
        }
})

var Group = React.createClass({
        render: function(){
                return (
                        <rect className='block' x='0' y='0' width='100' height='10' />
                )
        }
})

var Source = React.createClass({
        render: function(){
                return (
                        <rect className='block' x='0' y='0' width='10' height='10' /> 
                )      
        }
})

var Entity = React.createClass({
        render: function(){
                var element;
                switch(this.props.model.instance()){
                        case 'block':
                        element = <Block {...this.props} />
                        break;
                        case 'group':
                        element = <Group {...this.props} />
                        break;
                        case 'source':
                        element = <Source {...this.props }/>
                        break;
                }

                return(
                        <DragContainer {...this.props} x={this.props.model.position.x} y={this.props.model.position.y}>
                                {element}
                        </DragContainer>
                )
        }
})

var CoreApp = React.createClass({
    render: function() {
            return (
                    <svg className="stage" onDragOver={this.dragOver}>
                    {this.props.model.list.map(function(e){
                        return <Entity key={e.id} model={e} />
                    })}
                    </svg>
            )
    }
})

function render() {
        React.render(<CoreApp model={m}/> , document.getElementById('example'));
}

m.subscribe(render);
render();
