package loadbalancer

import (
	"net/http/httputil"
	"net/url"
)

type host struct {
	addr string
	proxy *httputil.ReverseProxy
	alive bool
}

type hostRing struct {
	hosts []host
	curr int
	len int
}

func newHost(addr string) (*host, error) {
	url, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	return &host{ addr, proxy, true }, nil
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

	return &hostRing{ hosts, 0, hostsLen }, nil
}

func (ring *hostRing) get() host {
	host := ring.hosts[ring.curr]
	if ring.curr == ring.len - 1 {
		ring.curr = 0
	} else {
		ring.curr++
	}
	return host
}
