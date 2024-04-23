package process

import (
	"errors"
	"sync"

	data "github.com/multiversx/mx-chain-ws-connector-template-go/data/hyperOutportBlocks"
)

type hyperOutportBlocksQueue struct {
	mu       sync.RWMutex
	Elements []*data.HyperOutportBlock
	Size     int
}

func NewHyperOutportBlocksQueue() *hyperOutportBlocksQueue {
	return &hyperOutportBlocksQueue{
		Elements: make([]*data.HyperOutportBlock, 0),
	}
}

func (q *hyperOutportBlocksQueue) Enqueue(elem *data.HyperOutportBlock) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	// TODO: add queue size to config
	//if q.GetLength() == q.Size {
	//	return errors.New("hyperOutportBlocksQueue is full")
	//}
	q.Elements = append(q.Elements, elem)
	return nil
}

func (q *hyperOutportBlocksQueue) Dequeue() (*data.HyperOutportBlock, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.IsEmpty() {
		return nil, errors.New("hyperOutportBlocksQueue is empty")
	}
	element := q.Elements[0]
	if q.GetLength() == 1 {
		q.Elements = nil
		return element, nil
	}
	q.Elements = q.Elements[1:]
	return element, nil // Slice off the element once it is dequeued.
}

func (q *hyperOutportBlocksQueue) GetLength() int {
	return len(q.Elements)
}

func (q *hyperOutportBlocksQueue) IsEmpty() bool {
	return len(q.Elements) == 0
}

func (q *hyperOutportBlocksQueue) Peek() (*data.HyperOutportBlock, error) {
	if q.IsEmpty() {
		return nil, errors.New("hyperOutportBlocksQueue is empty")
	}
	return q.Elements[0], nil
}
