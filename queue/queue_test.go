package queue

import (
	"testing"

	"github.com/ameliaikeda/crawler/node"
)

func TestQueueAddingNodes(t *testing.T) {
	t.Parallel()

	q := &queue{c: make(chan *node.Node, 400)}
	n := &node.Node{}

	q.Add(n)

	if len(q.c) != 1 {
		t.Error("Queue's channel length should be 1")
	}

	if q.remaining != 1 {
		t.Error("Queue's outstanding value should be 1")
	}

	r, ok := <-q.Get()

	if !ok {
		t.Error("queuing work returned a closed channel")
	}

	if r != n {
		t.Error("returned node from Queue.C isn't the same as the one queued")
	}
}

func TestQueueDeletingNodes(t *testing.T) {
	t.Parallel()

	q := &queue{c: make(chan *node.Node, 400)}
	n := &node.Node{}

	q.Add(n)

}

func TestDeletingWithoutNodes(t *testing.T) {
	t.Parallel()

	q := &queue{c: make(chan *node.Node, 400)}
	err := q.Del()

	if err != ErrQueueDrained {
		t.Error("deleting from a drained queue should return ErrQueueDrained")
	}
}

func TestNewQueue(t *testing.T) {
	t.Parallel()

	New(func() {})
}
