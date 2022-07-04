package loadbalancer

import "container/heap"

type leastConnections struct {
	pq PriorityQueue
	srvrToItem map[string]*item
}

func newLeastConnections(srvrs []*server) *leastConnections {
	numSrvrs := len(srvrs)
	pq := make(PriorityQueue, numSrvrs)
	srvrToItem := make(map[string]*item, numSrvrs)
	for i, server := range srvrs {
		item := &item{
			server,
			0,
			i,
		}
		pq[i] = item
		srvrToItem[server.host] = item
	}
	heap.Init(&pq)
	return &leastConnections{
		pq,
		srvrToItem,
	}
}

func (lc *leastConnections) choose() *server {
	item := heap.Pop(&lc.pq).(*item)
	server := item.value.(server)

	item.priority += 1
	heap.Push(&lc.pq, item)
	return &server
}

func (lc *leastConnections) after(srvr *server) {
	item := lc.srvrToItem[srvr.host]
	lc.pq.update(item, item.priority - 1)
}
