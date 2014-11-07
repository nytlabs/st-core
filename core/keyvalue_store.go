package core

type KeyValue struct {
	*Store
	data    map[interface{}]interface{}
	setChan chan setRequest
	getChan chan getRequest
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
		setChan: make(chan setRequest),
		getChan: make(chan getRequest),
	}
}

func (s KeyValue) Serve() {
	for {
		select {
		case keyvalue := <-s.setChan:
			s.data[keyvalue.key] = keyvalue.value
		case req := <-s.getChan:
			req.respChan <- s.data[req.key]
		}
	}
}
