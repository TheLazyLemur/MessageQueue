package server

type Queue struct {
	items []string
	front int
	rear  int
}

func NewQueue() *Queue {
	return &Queue{make([]string, 0), 0, -1}
}

func (q *Queue) Enqueue(item string) {
	q.rear++
	if q.rear == len(q.items) {
		q.items = append(q.items, item)
	} else {
		q.items[q.rear] = item
	}
}

func (q *Queue) Dequeue() string {
	if q.IsEmpty() {
		return ""
	}
	item := q.items[q.front]
	q.front++
	return item
}

func (q *Queue) IsEmpty() bool {
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
