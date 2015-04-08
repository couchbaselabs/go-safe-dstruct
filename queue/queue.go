package queue

import (
	"container/list"
	"sync"
)

// Thread-safe producer/consumer queue.  Copied from
// github.com/couchbaselabs/walrus
type Queue struct {
	list *list.List
	cond *sync.Cond
}

func NewQueue() *Queue {
	return &Queue{
		list: list.New(),
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

// Pushes a value into the queue. (Never blocks: the queue has no size limit.)
func (q *Queue) Push(value interface{}) {
	q.cond.L.Lock()
	q.list.PushFront(value)
	if q.list.Len() == 1 {
		q.cond.Signal()
	}
	q.cond.L.Unlock()
}

// Removes the last/oldest value from the queue; if the queue is empty, blocks.
func (q *Queue) Pull() interface{} {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	for q.list != nil && q.list.Len() == 0 {
		q.cond.Wait()
	}
	if q.list == nil {
		return nil // queue is closed
	}
	last := q.list.Back()
	q.list.Remove(last)
	return last.Value
}

func (q *Queue) Close() {
	q.cond.L.Lock()
	if q.list != nil {
		q.list = nil
		q.cond.Broadcast()
	}
	q.cond.L.Unlock()
}
