var app = app || {};

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
        console.log("updating...");
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

    app.CoreModel.prototype.select = function(id){
        console.log(id);
        console.log(this.entities[id]);
        this.list.push(this.list.splice(this.list.indexOf(this.entities[id]), 1)[0]);
        console.log(id);
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
                        offY: null
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
                app.Utils.request(
                        "PUT", 
                        "blocks/" + this.props.model.id + "/position", 
                        {x: this.state.x, y: this.state.y }, 
                        null
                );
              
                this.setState({
                        dragging: false,
                })
        },
        onMouseMove: function(e){
                if(this.state.dragging){
                        this.setState({
                                x: e.pageX - this.state.offX,
                                y: e.pageY - this.state.offY,
                        })
                }
        },
        componentWillReceiveProps: function(props){
                this.setState({
                      x: props.x,
                      y: props.y       
                })
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



var Entity = React.createClass({
        render: function(){
                var entity = this.props.model;
                if(entity.hasOwnProperty('inputs')){
                return(
                        <DragContainer {...this.props} x={this.props.model.position.x} y={this.props.model.position.y}>
                        <rect className='block' x='0' y='0' width='100' height='100' />
                        </DragContainer>
                )
                } else {
                return <div>LOL WHO CARES</div>
                }
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
