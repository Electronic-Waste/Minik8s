package queue

import "sync"

type Item interface{}

type Queue struct {
	items []Item
	mu    sync.Mutex
}

func (q *Queue) New() *Queue {
	q.items = []Item{}
	return q
}

func (q *Queue) Enqueue(item Item) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, item)
}

func (q *Queue) Dequeue() Item {
	q.mu.Lock()
	defer q.mu.Unlock()
	temp := q.items[0]
	q.items = q.items[1:]
	return temp
}

func (q *Queue) Front() Item {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.items[0]
}

func (q *Queue) Empty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items) == 0
}
