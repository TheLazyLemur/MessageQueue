package queue

import "sync"

type Queue struct {
	items []string
	front int
	rear  int
	lock  sync.RWMutex
}

func NewQueue() *Queue {
	return &Queue{
		items: make([]string, 0),
		front: 0,
		rear:  -1,
		lock:  sync.RWMutex{},
	}
}

func (q *Queue) Enqueue(item string) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.rear++
	if q.rear == len(q.items) {
		q.items = append(q.items, item)
	} else {
		q.items[q.rear] = item
	}
}

func (q *Queue) Dequeue() string {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.isEmpty() {
		return ""
	}

	item := q.items[q.front]
	q.front++
	return item
}

func (q *Queue) IsEmpty() bool {
	q.lock.RLock()
	defer q.lock.RUnlock()

	return q.isEmpty()
}

func (q *Queue) isEmpty() bool {
	return q.front > q.rear
}

func (q *Queue) Size() int {
	return q.rear - q.front + 1
}

func (q *Queue) Peek() string {
	if q.IsEmpty() {
		panic("queue is empty")
	}
	return q.items[q.front]
}
