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
	rr := roundRobin{
		nil,
		0,
		0,
		&mu,
	}
	rr.newInput(servers)
	return &rr
}

func (rr *roundRobin) newInput(servers []*server) {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	rr.curr = 0
	rr.len = len(servers)
	rr.servers = servers
}

func (rr *roundRobin) choose() *server {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	server := rr.servers[rr.curr % rr.len]
	rr.curr++
	return server
}

func (*roundRobin) after(*server) {

}
