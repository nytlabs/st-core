package core

import (
	"errors"
	"sync"
)

func KeyValueStore() SourceSpec {
	return SourceSpec{
		Name: "key_value",
		Type: KEY_VALUE,
		New:  NewKeyValue,
	}
}

func NewKeyValue() Source {
	return &KeyValue{
		kv:   make(map[string]interface{}),
		quit: make(chan bool),
	}
}

func (k KeyValue) GetType() SourceType {
	return KEY_VALUE
}

type KeyValue struct {
	kv   map[string]interface{}
	quit chan bool
	sync.Mutex
}

func (k *KeyValue) Get() interface{} {
	return k.kv
}

func (k *KeyValue) Set(v interface{}) error {
	kv, ok := v.(map[string]interface{})
	if !ok {
		return errors.New("not a map")
	}
	k.kv = kv
	return nil
}

// retrieves a value from the key value store
func kvGet() Spec {
	return Spec{
		Name: "kvGet",
		Inputs: []Pin{
			Pin{"key", STRING},
		},
		Outputs: []Pin{
			Pin{"value", ANY},
		},
		Source: KEY_VALUE,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
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
		Name: "kvSet",
		Inputs: []Pin{
			Pin{"key", STRING},
			Pin{"value", ANY},
		},
		Outputs: []Pin{
			Pin{"new", BOOLEAN},
		},
		Source: KEY_VALUE,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
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
		Name: "kvClear",
		Inputs: []Pin{
			Pin{"clear", ANY},
		},
		Outputs: []Pin{
			Pin{"cleared", BOOLEAN},
		},
		Source: KEY_VALUE,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			kv := s.(*KeyValue)
			kv.kv = make(map[string]interface{})
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
		Name: "kvDump",
		Inputs: []Pin{
			Pin{"dump", ANY},
		},
		Outputs: []Pin{
			Pin{"object", OBJECT},
		},
		Source: KEY_VALUE,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			kv := s.(*KeyValue)
			outMap := make(map[string]interface{})
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
		Name: "kvDelete",
		Inputs: []Pin{
			Pin{"key", STRING},
		},
		Outputs: []Pin{
			Pin{"deleted", BOOLEAN},
		},
		Source: KEY_VALUE,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
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
