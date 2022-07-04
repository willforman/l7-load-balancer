package loadbalancer

type roundRobin struct {
	servers []*server
	curr int
	len int
}

func newRoundRobin(servers []*server) *roundRobin {
	return &roundRobin{
		servers,
		0,
		len(servers),
	}
}

func (rr *roundRobin) choose() *server {
	for i := 0; i < rr.len; i++ {
		server := rr.servers[rr.curr]
		if rr.curr == rr.len - 1 {
			rr.curr = 0
		} else {
			rr.curr += 1
		}

		server.mu.Lock()
		alive := server.alive
		server.mu.Unlock()
		if alive {
			return server
		}
	}
	return nil
}

func (*roundRobin) after(*server) {

}

