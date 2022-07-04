package loadbalancer

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type server struct {
	host string
	proxy httputil.ReverseProxy
	alive bool
	mu *sync.Mutex
}

// Host must be of form host+port (no scheme)
func isAlive(host string) bool {
	conn, err := net.DialTimeout("tcp", host, time.Second * 5)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func (srvr *server) proxyErrHandler() func(w http.ResponseWriter, r *http.Request, err error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		srvr.mu.Lock()
		srvr.alive = false
		srvr.mu.Unlock()
	}
}

// SrvrUrl must be of form scheme+host+port
func newServer(srvrUrl string) (*server, error) {
	url, err := url.Parse(srvrUrl)
	if err != nil {
		return nil, fmt.Errorf("url parse err: %w", err)
	}
	if url.Scheme != "http" {
		return nil, fmt.Errorf("scheme given not http: %s", url.Scheme)
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	var mu sync.Mutex
	srvr := &server{url.Host, *proxy, isAlive(url.Host), &mu}
	srvr.proxy.ErrorHandler = srvr.proxyErrHandler()
	return srvr, nil
}

