package core

import "sync"

func NewKeyValue() Store {
	return &KeyValue{
		kv: make(map[string]Message),
	}
}

type KeyValue struct {
	kv map[string]Message
	sync.Mutex
}

// TODO: proper type checking on input key
func kvGet() Spec {
	return Spec{
		Inputs: []Pin{
			Pin{"key"},
		},
		Outputs: []Pin{
			Pin{"value"},
		},
		Kernel: func(in MessageMap, out MessageMap, s Store, i chan Interrupt) Interrupt {
			s.Lock()
			kv := s.(*KeyValue)
			if value, ok := kv.kv[in[0].(string)]; !ok {
				out[0] = "error"
			} else {
				out[0] = value
			}
			s.Unlock()
			return nil
		},
	}
}

// TODO: proper type checking on input key
func kvSet() Spec {
	return Spec{
		Inputs: []Pin{
			Pin{"key"},
			Pin{"value"},
		},
		Outputs: []Pin{
			Pin{"inserted"},
		},
		Kernel: func(in MessageMap, out MessageMap, s Store, i chan Interrupt) Interrupt {
			s.Lock()
			kv := s.(*KeyValue)
			if _, ok := kv.kv[in[0].(string)]; !ok {
				out[0] = true
			} else {
				out[0] = false
			}

			kv.kv[in[0].(string)] = in[1]
			s.Unlock()
			return nil
		},
	}
}

// TODO: prefer "empty"
// change interface{} to message
func kvClear() Spec {
	return Spec{
		Inputs: []Pin{
			Pin{"clear"},
		},
		Outputs: []Pin{
			Pin{"cleared"},
		},
		Kernel: func(in MessageMap, out MessageMap, s Store, i chan Interrupt) Interrupt {
			s.Lock()
			kv := s.(*KeyValue)
			kv.kv = make(map[string]Message)
			out[0] = true
			s.Unlock()
			return nil
		},
	}
}

// should output be named "object" ?
// TODO: convert interface{} to message
// !! should probably double check this to ensure that we don't need a deep copy
func kvDump() Spec {
	return Spec{
		Inputs: []Pin{
			Pin{"dump"},
		},
		Outputs: []Pin{
			Pin{"object"},
		},
		Kernel: func(in MessageMap, out MessageMap, s Store, i chan Interrupt) Interrupt {
			s.Lock()
			kv := s.(*KeyValue)
			outMap := make(map[string]Message)
			for k, v := range kv.kv {
				outMap[k] = v
			}
			out[0] = outMap
			s.Unlock()
			return nil
		},
	}
}

// TODO: proper type checking on input key
func kvDelete() Spec {
	return Spec{
		Inputs: []Pin{
			Pin{"key"},
		},
		Outputs: []Pin{
			Pin{"deleted"},
		},
		Kernel: func(in MessageMap, out MessageMap, s Store, i chan Interrupt) Interrupt {
			s.Lock()
			kv := s.(*KeyValue)
			if _, ok := kv.kv[in[0].(string)]; !ok {
				out[0] = false
			} else {
				delete(kv.kv, in[0].(string))
				out[0] = true
			}
			s.Unlock()
			return nil
		},
	}
}
