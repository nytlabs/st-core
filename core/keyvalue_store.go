package core

type KeyValue struct {
	*Store
	data    map[interface{}]interface{}
	setChan chan setRequest
	getChan chan Connection
}

type setRequest struct {
	key   interface{}
	value interface{}
}

type getRequest struct {
	key      interface{}
	respChan chan interface{}
}

func NewKeyValue(name string) KeyValue {
	s := NewStore(name)
	return KeyValue{
		Store:   s,
		data:    make(map[interface{}]interface{}),
		setChan: make(Connection),
		getChan: make(chan Connection),
	}
}

func (s KeyValue) Serve() {
	for {
		select {
		case keyvalue := <-s.setChan:
			kv, _ := keyvalue.([]interface{})
			s.data[kv[0]] = kv[1]
		case key := <-s.getChan:
		v:
		}
	}
}
