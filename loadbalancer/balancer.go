package loadbalancer

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Algorithm int

const (
	RoundRobin Algorithm = iota
	LeastConnections
)

type LoadBalancerArgs struct {
	Urls []url.URL
	Port  int
	Algorithm Algorithm
}

type serverSelector interface {
	makeReq(http.ResponseWriter, *http.Request)
}

type LoadBalancer struct {
	servers []server
	port     int
	selector serverSelector
}

func newSelector(algo Algorithm, servers []server) serverSelector {
	switch (algo) {
	case RoundRobin:
		return newRoundRobin(servers)
	case LeastConnections:
		return newLeastConnections(servers)
	}
	panic("invalid load balancing algorithm")
}

func NewLoadBalancer(args *LoadBalancerArgs) (*LoadBalancer, error) {
	serverLen := len(args.Urls)
	servers := make([]server, serverLen)
	for i, url := range args.Urls {
		server := newServer(&url)
		servers[i] = *server
	}


	return &LoadBalancer{
		servers,
		args.Port,
		newSelector(args.Algorithm, servers),
	}, nil
}

func (lb *LoadBalancer) handler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		lb.selector.makeReq(w, r)
	}
}

func (lb *LoadBalancer) startWatchDog() func() {
	ticker := time.NewTicker(time.Second * 30)

	for {
		<-ticker.C
		for _, server := range lb.servers {
			alive := isAlive(server.addr)
			if alive != server.alive {
				server.mu.Lock()
				server.alive = alive
				server.mu.Unlock()
			}
		}
	}
}

func (lb *LoadBalancer) Start(algoStr string) {
	log.Printf("starting load balancer: port=%d algo=%s\n", lb.port, algoStr)
	http.HandleFunc("/", lb.handler())
	go lb.startWatchDog()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", lb.port), nil))
}
