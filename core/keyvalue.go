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

// retrieves a value from the key value store
func kvGet() Spec {
	return Spec{
		Inputs: []Pin{
			Pin{"key"},
		},
		Outputs: []Pin{
			Pin{"value"},
		},
		Shared: KEY_VALUE,
		Kernel: func(in MessageMap, out MessageMap, s Store, i chan Interrupt) Interrupt {
			kv := s.(*KeyValue)
			key, ok := in[0].(string)
			if !ok {
				out[0] = NewError("Key is not type string")
				return nil
			}

			if value, ok := kv.kv[key]; !ok {
				out[0] = NewError("Key not found")
			} else {
				out[0] = value
			}
			return nil
		},
	}
}

// sets an entry in a key value store
// if the entry is new, emits true
func kvSet() Spec {
	return Spec{
		Inputs: []Pin{
			Pin{"key"},
			Pin{"value"},
		},
		Outputs: []Pin{
			Pin{"new"},
		},
		Shared: KEY_VALUE,
		Kernel: func(in MessageMap, out MessageMap, s Store, i chan Interrupt) Interrupt {
			kv := s.(*KeyValue)
			key, ok := in[0].(string)
			if !ok {
				out[0] = NewError("Key is not type string")
				return nil
			}

			if _, ok := kv.kv[key]; !ok {
				out[0] = true
			} else {
				out[0] = false
			}

			kv.kv[in[0].(string)] = in[1]
			return nil
		},
	}
}

// clears the entire map
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
		Shared: KEY_VALUE,
		Kernel: func(in MessageMap, out MessageMap, s Store, i chan Interrupt) Interrupt {
			kv := s.(*KeyValue)
			kv.kv = make(map[string]Message)
			out[0] = true
			return nil
		},
	}
}

// dumps the entire map into a message
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
		Shared: KEY_VALUE,
		Kernel: func(in MessageMap, out MessageMap, s Store, i chan Interrupt) Interrupt {
			kv := s.(*KeyValue)
			outMap := make(map[string]Message)
			for k, v := range kv.kv {
				outMap[k] = v
			}
			out[0] = outMap
			return nil
		},
	}
}

// deletes an entry in a key value store
func kvDelete() Spec {
	return Spec{
		Inputs: []Pin{
			Pin{"key"},
		},
		Outputs: []Pin{
			Pin{"deleted"},
		},
		Shared: KEY_VALUE,
		Kernel: func(in MessageMap, out MessageMap, s Store, i chan Interrupt) Interrupt {
			kv := s.(*KeyValue)
			key, ok := in[0].(string)
			if !ok {
				out[0] = NewError("Key is not type string")
				return nil
			}

			if _, ok := kv.kv[key]; !ok {
				out[0] = false
			} else {
				delete(kv.kv, key)
				out[0] = true
			}
			return nil
		},
	}
}
