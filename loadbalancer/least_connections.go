package loadbalancer

import (
	"container/heap"
	"sync"
)

type leastConnections struct {
	pq PriorityQueue
	serverToItem map[string]*item
	mu *sync.Mutex
}

func newLeastConnections(servers []*server) *leastConnections {
	var mu sync.Mutex
	lc := leastConnections{
		nil,
		nil,
		&mu,
	}
	lc.newInput(servers)
	return &lc
}

func (lc *leastConnections) newInput(servers []*server) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	numSrvrs := len(servers)
	lc.pq = make(PriorityQueue, numSrvrs)
	lc.serverToItem = make(map[string]*item, numSrvrs)
	for i, server := range servers {
		item := &item{
			server,
			0,
			i,
		}
		lc.pq[i] = item
		lc.serverToItem[server.host] = item
	}
	heap.Init(&lc.pq)
}

func (lc *leastConnections) choose() *server {
	item := heap.Pop(&lc.pq).(*item)
	server := item.value.(*server)

	item.priority += 1
	heap.Push(&lc.pq, item)
	return server
}

func (lc *leastConnections) after(srvr *server) {
	item := lc.serverToItem[srvr.host]
	lc.pq.update(item, item.priority - 1)
}
