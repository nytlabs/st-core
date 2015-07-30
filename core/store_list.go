package core

import (
	"errors"
	"sync"
)

func ListStore() SourceSpec {
	return SourceSpec{
		Name: "list",
		Type: LIST,
		New:  NewList,
	}
}

func NewList() Source {
	return &List{
		list: make([]interface{}, 0),
		quit: make(chan bool),
	}
}

type List struct {
	list []interface{}
	quit chan bool
	sync.Mutex
}

func (k List) GetType() SourceType {
	return LIST
}

func (l *List) Get() interface{} {
	return l.list
}

func (l *List) Set(v interface{}) error {
	list, ok := v.([]interface{})
	if !ok {
		return errors.New("not a slice")
	}
	l.list = list
	return nil
}

// retrieves an element from the list by index
func listGet() Spec {
	return Spec{
		Name: "listGet",
		Inputs: []Pin{
			Pin{"index", NUMBER},
		},
		Outputs: []Pin{
			Pin{"element", ANY},
		},
		Source: LIST,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			l := s.(*List)
			indexFloat, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("List index is not a Number ")
				return nil
			}
			index := int(indexFloat)
			/*
				// TODO should we be happy with int() just flooring the index?
				if math.Mod(a, math.Floor(a)) > 0 {
					out[0] = NewError("List index must be an integer")
					return nil
				}
			*/
			if index < 0 {
				out[0] = NewError("List index must be ≥ 0")
				return nil
			}
			if index > len(l.list) {
				out[0] = NewError("List index out of range")
				return nil
			}
			out[0] = l.list[index]
			return nil
		},
	}
}

// sets a value in the list by index
func listSet() Spec {
	return Spec{
		Name: "listSet",
		Inputs: []Pin{
			Pin{"index", NUMBER}, Pin{"element", ANY},
		},
		Outputs: []Pin{
			Pin{"out", ANY},
		},
		Source: LIST,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			l := s.(*List)
			indexFloat, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("List index is not a Number ")
				return nil
			}
			index := int(indexFloat)
			/*
				// TODO should we be happy with int() just flooring the index?
				if math.Mod(a, math.Floor(a)) > 0 {
					out[0] = NewError("List index must be an integer")
					return nil
				}
			*/
			if index < 0 {
				out[0] = NewError("List index must be ≥ 0")
				return nil
			}
			if index > len(l.list)-1 {
				out[0] = NewError("List index out of range")
				return nil
			}
			l.list[index] = in[1]
			out[0] = true
			return nil
		},
	}
}

// listAppend appends an element to the end of a list
func listAppend() Spec {
	return Spec{
		Name: "listAppend",
		Inputs: []Pin{
			Pin{"element", ANY},
		},
		Outputs: []Pin{
			Pin{"out", BOOLEAN},
		},
		Source: LIST,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			l := s.(*List)
			out[0] = true
			l.list = append(l.list, in[0])
			return nil
		},
	}
}

// listPop pops an element off the end of the list
func listPop() Spec {
	return Spec{
		Name: "listPop",
		Inputs: []Pin{
			Pin{"trigger", ANY},
		},
		Outputs: []Pin{
			Pin{"element", ANY},
		},
		Source: LIST,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			l := s.(*List)
			if len(l.list) == 0 {
				out[0] = NewError("empty list")
				return nil
			}
			out[0], l.list = l.list[len(l.list)-1], l.list[:len(l.list)-1]
			return nil
		},
	}
}

// listShift adds an element to the front of a list
func listShift() Spec {
	return Spec{
		Name: "listShift",
		Inputs: []Pin{
			Pin{"element", ANY},
		},
		Outputs: []Pin{
			Pin{"out", BOOLEAN},
		},
		Source: LIST,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			l := s.(*List)
			newList := make([]interface{}, len(l.list)+1)
			newList[0] = in[0]
			copy(newList[1:], l.list)
			l.list = newList
			out[0] = true
			return nil
		},
	}
}

// listDump returns the list
func listDump() Spec {
	return Spec{
		Name:    "listDump",
		Inputs:  []Pin{Pin{"trigger", ANY}},
		Outputs: []Pin{Pin{"list", ARRAY}},
		Source:  LIST,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			l := s.(*List)
			out[0] = l.list
			return nil
		},
	}
}
