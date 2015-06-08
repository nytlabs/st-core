package core

import (
	"container/heap"
	"sync"
)

type queue []*PQMessage

type PriorityQueue struct {
	queue *queue
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
	pq := &queue{}
	heap.Init(pq)
	return &PriorityQueue{
		queue: pq,
	}
}

func (pq PriorityQueue) GetType() SourceType {
	return PRIORITY
}

func (pq queue) Len() int {
	return len(pq)
}

func (pq queue) Less(i, j int) bool {
	return pq[i].t > pq[j].t
}

func (pq queue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *queue) Push(x interface{}) {
	n := len(*pq)
	msg := x.(*PQMessage)
	msg.index = n
	*pq = append(*pq, msg)
}

func (pq *queue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *queue) Peek() interface{} {
	if pq.Len() == 0 {
		return nil
	}
	return (*pq)[0]
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
			msg := &PQMessage{
				val: in[0],
				t:   int(priority),
			}
			heap.Push(pq.queue, msg)
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
			msgI := pq.queue.Pop()
			msg, ok := msgI.(*PQMessage)
			if !ok {
				out[0] = NewError("pulled something weird off the PriorityQueue")
				return nil
			}
			out[0] = msg.val
			out[1] = float64(msg.t)
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
			msgI := pq.queue.Peek()
			msg, ok := msgI.(*PQMessage)
			if !ok {
				out[0] = NewError("pulled something weird off the PriorityQueue")
				return nil
			}
			out[0] = msg.val
			out[1] = float64(msg.t)
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
			out[0] = len(*pq.queue)
			return nil
		},
	}
}
