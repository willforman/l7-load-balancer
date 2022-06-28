package loadbalancer

import "net/http"

type roundRobin struct {
	servers []server
	curr int
	len int
}

func newRoundRobin(servers []server) *roundRobin {
	return &roundRobin{
		servers,
		0,
		len(servers),
	}
}

func (rr *roundRobin) get() *server {
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
			return &server
		}
	}
	return nil
}

func (rr *roundRobin) makeReq(w http.ResponseWriter, r *http.Request) {
	server := rr.get()
	server.proxy.ServeHTTP(w, r)
}
