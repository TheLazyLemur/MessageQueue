package queue

import "testing"

func TestQueue(t *testing.T) {
	q := NewQueue()

	q.Enqueue("1")
	q.Enqueue("2")
	q.Enqueue("3")
	q.Enqueue("4")

	if q.Dequeue() != "1" {
		t.Errorf("Expected 1, got %v", q.Dequeue())
	}

	if q.Dequeue() != "2" {
		t.Errorf("Expected 2, got %v", q.Dequeue())
	}

	if q.Dequeue() != "3" {
		t.Errorf("Expected 3, got %v", q.Dequeue())
	}

	if q.Dequeue() != "4" {
		t.Errorf("Expected 4, got %v", q.Dequeue())
	}

	if q.Dequeue() != "" {
		t.Errorf("Expected empty, got %v", q.Dequeue())
	}
}

func TestQueueIsEmpty(t *testing.T) {
	q := NewQueue()

	if !q.IsEmpty() {
		t.Errorf("Expected true, got %v", q.IsEmpty())
	}

}

func TestQueueSize(t *testing.T) {
	q := NewQueue()

	q.Enqueue("1")

	if q.Size() != 1 {
		t.Errorf("Expected 1, got %v", q.Size())
	}
}

func TestQueuePeek(t *testing.T) {
	q := NewQueue()

	q.Enqueue("1")
	q.Enqueue("2")

	if q.Peek() != "1" {
		t.Errorf("Expected 1, got %v", q.Peek())
	}
}
