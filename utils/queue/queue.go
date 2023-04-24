package queue

import "sync"

//a simple queue
type ConcurrentQueue struct {
	queue []interface{}
	mu    sync.Mutex
}

func (q *ConcurrentQueue) Enqueue(item interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue = append(q.queue, item)
}

func (q *ConcurrentQueue) Dequeue() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()
	temp := q.queue[0]
	q.queue = q.queue[1:]
	return temp
}

func (q *ConcurrentQueue) Front() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.queue[0]
}

func (q *ConcurrentQueue) Empty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.queue) == 0
}
