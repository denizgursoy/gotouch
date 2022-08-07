package model

type Queue struct {
	items []*QueueElement
}

func (q *Queue) push(element *QueueElement) {
	q.items = append(q.items, element)
}

func (q *Queue) pop() *QueueElement {
	element := q.items[0]
	q.items = append(q.items[:0], q.items[1:]...)
	return element
}

type QueueElement struct {
}
