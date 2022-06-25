package loadbalancer

import (
	"net"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type host struct {
	url *url.URL
	proxy *httputil.ReverseProxy
	alive bool
	mu sync.Mutex
}

type hostRing struct {
	hosts []host
	curr  int
	len   int
}

func isAlive(addr string) bool {
	conn, _ := net.DialTimeout("tcp", addr, time.Second * 10)
	println(addr, ":", conn != nil)
	if conn != nil {
		defer conn.Close()
	}
	return conn != nil
}

func newHost(addr string) (*host, error) {
	url, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	var mu sync.Mutex
	return &host{url, proxy, isAlive(url.Host), mu}, nil
}

func newHostRing(addrs []string) (*hostRing, error) {
	hostsLen := len(addrs)
	hosts := make([]host, hostsLen)

	for i, addr := range addrs {
		host, err := newHost(addr)
		if err != nil {
			return nil, err
		}
		hosts[i] = *host
	}

	return &hostRing{hosts, 0, hostsLen}, nil
}

func (ring *hostRing) get() *host {
	host := &ring.hosts[ring.curr]
	if ring.curr == ring.len-1 {
		ring.curr = 0
	} else {
		ring.curr++
	}
	return host
}

func (ring *hostRing) getAlive() *host {
	for i := 0; i < ring.len; i++ {
		host := ring.get()
		host.mu.Lock()
		alive := host.alive
		host.mu.Unlock()
		if alive {
			return host
		}
	}
	return nil
}

func (ring *hostRing) doAll(fn func(*host)) {
	for i := 0; i < ring.len; i++ {
		fn(&ring.hosts[i])
	}
}
