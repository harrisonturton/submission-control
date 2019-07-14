package ci

import "container/list"

// Queue is a queue
type Queue interface {
	Front() *list.Element
	Len() int
	Add(interface{})
	Remove()
}

type queueImpl struct {
	*list.List
}

func (q *queueImpl) Add(v interface{}) {
	q.PushBack(v)
}

func (q *queueImpl) Remove() {
	e := q.Front()
	q.List.Remove(e)
}

// NewQueue is a new instance of a Queue
func NewQueue() Queue {
	return &queueImpl{list.New()}
}
