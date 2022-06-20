package main

import (
	"net/http/httputil"
	"net/url"
)

type Host struct {
	addr string
	proxy *httputil.ReverseProxy
	alive bool
}

type HostRing struct {
	hosts []Host
	curr int
	len int
}

func newHost(addr string) (*Host, error) {
	url, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	return &Host{ addr, proxy, true }, nil
}

func newHostRing(addrs []string) (*HostRing, error) {
	hostsLen := len(addrs)
	hosts := make([]Host, hostsLen)

	for i, addr := range addrs {
		host, err := newHost(addr)
		if err != nil {
			return nil, err
		}
		hosts[i] = *host
	}

	return &HostRing{ hosts, 0, hostsLen }, nil
}

func (ring *HostRing) get() Host {
	host := ring.hosts[ring.curr]
	if ring.curr == ring.len - 1 {
		ring.curr = 0
	} else {
		ring.curr++
	}
	return host
}
