package loadbalancer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type algo int

const (
	RoundRobin algo = iota
	LeastConnections
)

type serverSelector interface {
	newInput([]*server)
	choose() *server
	after(*server)
}

type LoadBalancer struct {
	servers []*server
	port     int
	selector serverSelector
	lbServer *http.Server
}

func newSelector(algo algo, servers []*server) serverSelector {
	switch (algo) {
	case RoundRobin:
		return newRoundRobin(servers)
	case LeastConnections:
		return newLeastConnections(servers)
	}
	panic("invalid load balancing algorithm")
}

func proxyHandler(selector serverSelector, servers []*server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		srvr := selector.choose()
		for srvr != nil {
			ok := true
			srvr.proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
				ok = false
			}
			srvr.proxy.ServeHTTP(w, r)
			if ok {
				selector.after(srvr)
				return
			}
			aliveSrvrs := healthCheck(servers)
			selector.newInput(aliveSrvrs)
			srvr = selector.choose()
		}
		w.WriteHeader(503) 
	}
}

func NewLoadBalancer(port int, algoStr string, urls []string) (*LoadBalancer, error) {
	if port < 1024 || port > 65535 {
		return nil, fmt.Errorf("port out of range 1024 < p < 65535 [%d]", port)
	}

	var algo algo
	switch (algoStr) {
	case "lc":
		algo = LeastConnections
	case "rr":
		algo = RoundRobin
	default:
		return nil, fmt.Errorf("algo choice not lc or rr : %s", algoStr)
	}

	serverLen := len(urls)
	if serverLen == 0 {
		return nil, fmt.Errorf("must provide at least one server")
	}
	servers := make([]*server, serverLen)
	for i, addr := range urls {
		server, err := newServer(addr)
		if err != nil {
			return nil, fmt.Errorf("newServer: %w", err)
		}
		servers[i] = server
	}

	selector := newSelector(algo, servers)

	handler := http.NewServeMux()
	handler.HandleFunc("/", proxyHandler(selector, servers))
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	return &LoadBalancer{
		servers,
		port,
		selector,
		server,
	}, nil
}

func healthCheck(allSrvrs []*server) []*server {
	var aliveSrvrs []*server
	for _, server := range allSrvrs {
		if isAlive(server.host) {
			aliveSrvrs = append(aliveSrvrs, server)
		}
	}
	return aliveSrvrs
}

func (lb *LoadBalancer) periodicHealthCheck() func() {
	ticker := time.NewTicker(time.Second * 20)

	for {
		<-ticker.C
		aliveSrvrs := healthCheck(lb.servers)
		lb.selector.newInput(aliveSrvrs)
	}
}

func (lb *LoadBalancer) Start() {
	go lb.periodicHealthCheck()
	err := lb.lbServer.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %v", err)
	}
}

func (lb *LoadBalancer) Stop() error {
	return lb.lbServer.Shutdown(context.Background())
}

