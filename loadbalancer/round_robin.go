package loadbalancer

import "sync"

type roundRobin struct {
	servers []*server
	curr int
	len int
	mu *sync.Mutex
}

func newRoundRobin(servers []*server) *roundRobin {
	var mu sync.Mutex
	return &roundRobin{
		servers,
		0,
		len(servers),
		&mu,
	}
}

func (rr *roundRobin) choose() *server {
	for i := 0; i < rr.len; i++ {
		rr.mu.Lock()
		server := rr.servers[rr.curr]
		rr.mu.Unlock()
		if rr.curr == rr.len - 1 {
			rr.curr = 0
		} else {
			rr.curr += 1
		}
		return server
	}
	return nil
}

func (*roundRobin) after(*server) {

}

func (rr *roundRobin) passAliveServers(newSrvrs []*server) {
	rr.mu.Lock()
	rr.servers = newSrvrs
	rr.curr = 0
	rr.len = len(newSrvrs)
	rr.mu.Unlock()
}
