package loadbalancer

import (
	"container/heap"
	"net/http"
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

type leastConnections struct {
	pq PriorityQueue
}

func newLeastConnections(servers []server) *leastConnections {
	pq := make(PriorityQueue, len(servers))
	for i, server := range servers {
		pq[i] = &item{
			server,
			0,
			i,
		}
	}
	heap.Init(&pq)
	return &leastConnections{
		pq,
	}
}

func (lc *leastConnections) makeReq(w http.ResponseWriter, r *http.Request) {
	item := heap.Pop(&lc.pq).(*item)
	server := item.value.(server)

	item.priority += 1
	heap.Push(&lc.pq, item)
	server.proxy.ServeHTTP(w, r)
	lc.pq.update(item, item.priority - 1)
}
