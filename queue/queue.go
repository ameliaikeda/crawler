package queue

import (
	"errors"
	"sync"

	"github.com/ameliaikeda/crawler/node"
)

var ErrQueueDrained = errors.New("crawler: cannot delete from the queue when there are no entries left")

// Queue is a simple type that returns a channel for reading nodes.
type Queue interface {
	Get() <-chan *node.Node
	Del() error
	Add(*node.Node)
}

// queue is a channel struct with a counter attached to it for outstanding HTTP requests
type queue struct {
	sync.Mutex

	c         chan *node.Node
	remaining int
	cancel    func()
}

// Get returns a read channel for pulling additional nodes out of.
func (q *queue) Get() <-chan *node.Node {
	return q.c
}

// Add will add a node to the queue and mark a request as outstanding.
// A caller should defer a call to Del() to ensure it removes requests when done.
// Callers MUST queue additional work they encounter before exiting, or risk a race condition.
func (q *queue) Add(n *node.Node) {
	q.Lock()
	defer q.Unlock()

	q.remaining++
	q.c <- n
}

// Del removes a counter for an outstanding request.
// If this causes the queue to hit 0, this will halt all workers gracefully.
func (q *queue) Del() error {
	q.Lock()
	defer q.Unlock()

	if q.remaining <= 0 {
		return ErrQueueDrained
	}

	q.remaining--

	if q.remaining == 0 {
		q.cancel()
	}

	return nil
}

func New(cancelFunc func()) Queue {
	return &queue{
		c:         make(chan *node.Node, 10000),
		remaining: 0,
		cancel:    cancelFunc,
	}
}
