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
	item := x.(*PQMessage)
	item.index = n
	*pq = append(*pq, item)
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
