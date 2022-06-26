package loadbalancer

type roundRobin struct {
	curr int
	len int
}

func (rr *roundRobin) get(servers []server) *server {
	for i := 0; i < rr.len; i++ {
		server := servers[rr.curr]
		if rr.curr == rr.len - 1 {
			rr.curr = 0
		} else {
			rr.curr += 1
		}

		server.mu.Lock()
		alive := server.alive
		server.mu.Unlock()
		if alive {
			return &server
		}
	}
	return nil
}
