package loadbalancer

import (
	"container/heap"
)

type item struct {
	value any
	priority int
	index int
}

type PriorityQueue []*item

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i int, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i int, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n - 1]
	old[n - 1] = nil
	item.index = -1
	*pq = old[0 : n - 1]
	return item
}

func (pq *PriorityQueue) Push(x any) {
	item := x.(*item)
	item.index = len(*pq)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) update(item *item, priority int) {
	item.priority = priority
	heap.Fix(pq, item.index)
}
