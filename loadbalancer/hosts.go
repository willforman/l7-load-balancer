package loadbalancer

import (
	"net"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type server struct {
	url *url.URL
	proxy *httputil.ReverseProxy
	alive bool
	mu sync.Mutex
}

type serverRing struct {
	servers []server
	curr  int
	len   int
}

func isAlive(addr string) bool {
	conn, _ := net.DialTimeout("tcp", addr, time.Second * 10)
	if conn != nil {
		conn.Close()
		return true
	}
	return false
}

func newServer(addr string) (*server, error) {
	url, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	var mu sync.Mutex
	return &server{url, proxy, isAlive(url.Host), mu}, nil
}

func newServerRing(addrs []string) (*serverRing, error) {
	serversLen := len(addrs)
	servers := make([]server, serversLen)

	for i, addr := range addrs {
		server, err := newServer(addr)
		if err != nil {
			return nil, err
		}
		servers[i] = *server
	}

	return &serverRing{servers, 0, serversLen}, nil
}

func (ring *serverRing) get() *server {
	server := &ring.servers[ring.curr]
	if ring.curr == ring.len-1 {
		ring.curr = 0
	} else {
		ring.curr++
	}
	return server
}

func (ring *serverRing) getAlive() *server {
	for i := 0; i < ring.len; i++ {
		server := ring.get()
		server.mu.Lock()
		alive := server.alive
		server.mu.Unlock()
		if alive {
			return server
		}
	}
	return nil
}

func (ring *serverRing) doAll(fn func(*server)) {
	for i := 0; i < ring.len; i++ {
		fn(&ring.servers[i])
	}
}
