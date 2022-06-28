package loadbalancer

import (
	"net"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type server struct {
	addr string
	proxy httputil.ReverseProxy
	alive bool
	mu *sync.Mutex
}

func isAlive(addr string) bool {
	conn, _ := net.DialTimeout("tcp", addr, time.Second * 10)
	if conn != nil {
		conn.Close()
		return true
	}
	return false
}

func newServer(url *url.URL) *server {
	proxy := httputil.NewSingleHostReverseProxy(url)
	var mu sync.Mutex
	return &server{url.Host, *proxy, isAlive(url.Host), &mu}
}

