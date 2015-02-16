package core

import (
	"container/heap"
	"sync"
	"time"
)

type queue []*PQMessage

type PriorityQueue struct {
	queue *queue
	quit  chan bool
	in    chan interface{}
	sync.Mutex
}

type PQMessage struct {
	val   interface{}
	t     time.Time
	index int
}

func PriorityQueueSource() SourceSpec {
	return SourceSpec{
		Name: "priority-queue",
		Type: PRIORITY,
		New:  NewPriorityQueue,
	}
}

func NewPriorityQueue() Source {
	pq := &queue{}
	return &PriorityQueue{
		queue: pq,
		quit:  make(chan bool),
		in:    make(chan interface{}),
	}
}

func (pq PriorityQueue) GetType() SourceType {
	return PRIORITY
}

func (pq PriorityQueue) Describe() map[string]string {
	return map[string]string{}
}

func (pq PriorityQueue) SetSourceParameter(key, value string) {
}

func (pq PriorityQueue) Serve() {
	waitTimer := time.NewTimer(100 * time.Millisecond)
	window := time.Duration(0)
	heap.Init(pq.queue)
	for {
		select {
		case <-waitTimer.C:
		case val := <-pq.in:
			queueMessage := &PQMessage{
				val: val,
				t:   time.Now(),
			}
			heap.Push(pq.queue, queueMessage)
		case <-pq.quit:
			return
		}
		for {
			pqMsg, diff := pq.queue.PeekAndShift(time.Now(), window)
			if pqMsg == nil {
				// either the queue is empty, or it"s not time to emit
				if diff == 0 {
					// then the queue is empty. Pause for 5 seconds before checking again
					diff = time.Duration(500) * time.Millisecond
				}
				waitTimer.Reset(diff)
				break
			}
		}
	}
}

func (pq PriorityQueue) Stop() {
	pq.quit <- true
}

func (pq queue) Len() int {
	return len(pq)
}

func (pq queue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].t.Before(pq[j].t)
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

func (pq *queue) PeekAndShift(max time.Time, lag time.Duration) (interface{}, time.Duration) {
	if pq.Len() == 0 {
		return nil, 0
	}

	item := (*pq)[0]

	if item.t.Add(lag).Before(max) {
		heap.Remove(pq, 0)
		return item, 0
	}

	return nil, lag - max.Sub(item.t)
}

func pqPush() Spec {
	return Spec{
		Name: "pqPush",
		Inputs: []Pin{
			Pin{"in"},
			Pin{"timestamp"},
		},
		Outputs: []Pin{
			Pin{"out"},
		},
		Source: PRIORITY,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			pq := s.(*PriorityQueue)
			timestamp, ok := in[1].(float64)
			if !ok {
				out[0] = NewError("pqPush needs a Number for a timestamp")
				return nil
			}
			nsec := timestamp * 1000000
			t := time.Unix(0, int64(nsec))
			msg := PQMessage{
				val: in[0],
				t:   t,
			}
			pq.queue.Push(msg)
			out[0] = true
			return nil
		},
	}
}

func pqPop() Spec {
	return Spec{
		Name: "pqPop",
		Inputs: []Pin{
			Pin{"trigger"},
		},
		Outputs: []Pin{
			Pin{"out"},
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
			return nil
		},
	}
}

func pqPeek() Spec {
	return Spec{
		Name: "pqPeek",
		Inputs: []Pin{
			Pin{"trigger"},
		},
		Outputs: []Pin{
			Pin{"out"},
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
			return nil
		},
	}
}

func pqLen() Spec {
	return Spec{
		Name: "pqLen",
		Inputs: []Pin{
			Pin{"trigger"},
		},
		Outputs: []Pin{
			Pin{"length"},
		},
		Source: PRIORITY,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			pq := s.(*PriorityQueue)
			out[0] = len(pq.queue)
			return nil
		},
	}
}

// pqPeekAndShift blocks until a message is ready on the priority queue
func pqPeekAndShift() Spec {
	return Spec{
		Name: "pqPeek",
		Inputs: []Pin{
			Pin{"window"}, // this is the time window of the priority queue
		},
		Outputs: []Pin{
			Pin{"message"}, // this is the next message that's ready
		},
		Source: PRIORITY,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			pq := s.(*PriorityQueue)
			waitTimer := time.NewTimer(100 * time.Millisecond)
			window := in[0].(float64)
			for {
				select {
				case <-waitTimer.C:
				case interrupt := <-i:
					return interrupt
				}

				pqMsg, diff := pq.queue.PeekAndShift(time.Now(), window)
				if pqMsg == nil {
					// either the queue is empty, or it"s not time to emit
					if diff == 0 {
						// then the queue is empty. Pause for 5 seconds before checking again
						diff = time.Duration(500) * time.Millisecond
					}
					waitTimer.Reset(diff)
					continue
				}
			}
			out[0] = pqMsg
		},
	}
}
