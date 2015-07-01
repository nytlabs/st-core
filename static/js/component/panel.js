var app = app || {};

/* PanelComponent & PanelEditableComponent
 * Produces a list of fields that are the current representation of input
 * values for blocks/groups that are sent to the component.
 *
 * TODO: fix the {'data': ...} nonsense
 */

(function() {
    'use strict';

    app.PanelEditableComponent = React.createClass({
        displayName: 'PanelEditableComponent',
        getInitialState: function() {
            return {
                isEditing: false,
                value: this.props.value
            }
        },
        handleClick: function() {
            this.setState({
                isEditing: true,
                value: this.props.value
            }, function() {
                this.refs.editableInput.getDOMNode().focus();
                this.refs.editableInput.getDOMNode().select();
            });
        },
        handleKeyUp: function(e) {
            if (e.nativeEvent.keyCode === 13) {
                this.props.onChange(e.target.value);
                this.setState({
                    isEditing: false,
                });
            }
        },
        handleBlur: function() {
            this.setState({
                isEditing: false,
            });
        },
        render: function() {
            var value = this.props.value.length === 0 ? '<empty>' : this.props.value;
            var inputStyle = {
                display: this.state.isEditing ? 'block' : 'none'
            }
            var style = {
                display: this.state.isEditing ? 'none' : 'block'
            }

            return React.createElement('div', {
                className: 'editable'
            }, [
                React.createElement('input', {
                    defaultValue: this.state.value,
                    onKeyUp: this.handleKeyUp,
                    onBlur: this.handleBlur,
                    style: inputStyle,
                    ref: 'editableInput',
                    key: 'editableInput'
                }, null),
                React.createElement('div', {
                    onClick: this.handleClick,
                    style: style,
                    key: 'editableDisplay'
                }, value)
            ]);
        }
    })
})();

(function() {

app.ParametersPanelComponent = React.createClass({
  displayName: 'ParametersPanelComponent',
  render: function() {
    console.log(this.props.model.data.params)
    var id = this.props.model.data.id
    return React.createElement('div', {className:'panel'},[
      React.createElement('div', {key:'block_header', className:'block_header'}, this.props.model.data.type),
      React.createElement('div', {key:'block_label', className:'label'}, 'label'),
        React.createElement(app.PanelEditableComponent, {
          key: 'route_label',
          className: 'editable',
          value: this.props.model.data.label,
          onChange: function(value) {
              app.Utils.request(
                  'PUT',
                  this.props.model.instance() + 's/' + this.props.model.data.id + '/label',
                  value,
                  null
              )
          }.bind(this)
      }, null),
      this.props.model.data.params.map(function(p,i){
        console.log("boo!", p)
        return [
          React.createElement('div', { className: 'label', }, p.name),
          React.createElement(app.PanelEditableComponent, {
            key: id + p.name ,
            value: p.value,
            onChange: function(value) {
              console.log("hi",[{name:p.name, value:value}])
              app.Utils.request( 'PUT', 'sources/' + id + '/params', [{name:p.name, value:value}], null)
            }.bind(this)
          }, null)
        ]
        
      })
    ])
  }
})

}());

(function() {
    app.RoutesPanelComponent = React.createClass({
        displayName: 'PanelComponent',
        render: function() {
            return React.createElement('div', {
                className: 'panel'
            }, [
                React.createElement('div', {
                    key: 'block_header',
                    className: 'block_header',
                }, this.props.model.data.type),
                React.createElement('div', {
                    key: 'block_Label',
                    className: 'label',
                }, 'label'),
                React.createElement(app.PanelEditableComponent, {
                    key: 'route_label',
                    className: 'editable',
                    value: this.props.model.data.label,
                    onChange: function(value) {
                        app.Utils.request(
                            'PUT',
                            this.props.model.instance() + 's/' + this.props.model.data.id + '/label',
                            value,
                            null
                        )
                    }.bind(this)
                }, null),
                this.props.model.routes.filter(function(r) {
                    return r.direction === 'input'
                }).map(function(r, i) {
                    return [
                        React.createElement('div', {
                            className: 'label',
                        }, r.data.name),
                        React.createElement(app.PanelEditableComponent, {
                                key: r.id + r.data.name + r.index,
                                value: JSON.stringify(r.data.value),
                                onChange: function(value) {
                                    app.Utils.request(
                                        'PUT',
                                        'blocks/' + r.id + '/routes/' + r.index,
                                        JSON.parse(value),
                                        null
                                    )
                                }.bind(this)
                            },
                            null)
                    ]
                }.bind(this))
            ]);
        }
    })
})();
