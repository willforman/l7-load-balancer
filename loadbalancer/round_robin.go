package loadbalancer

import "net/http"

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

func (rr *roundRobin) makeReq(servers []server, w http.ResponseWriter, r *http.Request) {
	server := rr.get(servers)
	server.proxy.ServeHTTP(w, r)
}
