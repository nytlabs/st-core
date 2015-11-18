package core

import (
	"sync"

	"github.com/oleiade/lane"
)

type queue []*PQMessage

type PriorityQueue struct {
	queue *lane.PQueue
	sync.Mutex
}

type PQMessage struct {
	val   interface{}
	t     int
	index int
}

func PriorityQueueStore() SourceSpec {
	return SourceSpec{
		Name: "priority-queue",
		Type: PRIORITY,
		New:  NewPriorityQueue,
	}
}

func NewPriorityQueue() Source {
	var pq *lane.PQueue = lane.NewPQueue(lane.MINPQ)
	return &PriorityQueue{
		queue: pq,
	}
}

func (pq PriorityQueue) GetType() SourceType {
	return PRIORITY
}

func pqPush() Spec {
	return Spec{
		Name: "pqPush",
		Inputs: []Pin{
			Pin{"in", ANY},
			Pin{"priority", NUMBER},
		},
		Outputs: []Pin{
			Pin{"out", BOOLEAN},
		},
		Source: PRIORITY,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			pq := s.(*PriorityQueue)
			priority, ok := in[1].(float64)
			if !ok {
				out[0] = NewError("pqPush needs a Number for a priority")
				return nil
			}
			pq.queue.Push(in[0], int(priority))
			out[0] = true
			return nil
		},
	}
}

func pqPop() Spec {
	return Spec{
		Name: "pqPop",
		Inputs: []Pin{
			Pin{"trigger", ANY},
		},
		Outputs: []Pin{
			Pin{"out", ANY},
			Pin{"priority", NUMBER},
		},
		Source: PRIORITY,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			pq := s.(*PriorityQueue)
			if pq.queue.Size() == 0 {
				out[0] = NewError("empty PriorityQueue")
				return nil
			}
			msg, priority := pq.queue.Pop()
			out[0] = msg
			out[1] = float64(priority)
			return nil
		},
	}
}

func pqPeek() Spec {
	return Spec{
		Name: "pqPeek",
		Inputs: []Pin{
			Pin{"trigger", ANY},
		},
		Outputs: []Pin{
			Pin{"out", ANY},
			Pin{"priority", NUMBER},
		},
		Source: PRIORITY,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			pq := s.(*PriorityQueue)
			if pq.queue.Size() == 0 {
				out[0] = NewError("empty PriorityQueue")
				return nil
			}
			msg, priority := pq.queue.Head()
			out[0] = msg
			out[1] = float64(priority)
			return nil
		},
	}
}

func pqLen() Spec {
	return Spec{
		Name: "pqLen",
		Inputs: []Pin{
			Pin{"trigger", ANY},
		},
		Outputs: []Pin{
			Pin{"length", NUMBER},
		},
		Source: PRIORITY,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			pq := s.(*PriorityQueue)
			out[0] = float64(pq.queue.Size())
			return nil
		},
	}
}

func pqClear() Spec {
	return Spec{
		Name: "pqClear",
		Inputs: []Pin{
			Pin{"clear", ANY},
		},
		Outputs: []Pin{
			Pin{"cleared", BOOLEAN},
		},
		Source: PRIORITY,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			pq := s.(*PriorityQueue)
			for pq.queue.Size() != 0 {
				_, _ = pq.queue.Pop()
			}
			out[0] = true
			return nil
		},
	}
}
