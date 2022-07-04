package loadbalancer

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Algorithm int

const (
	RoundRobin Algorithm = iota
	LeastConnections
)

type LoadBalancerArgs struct {
	Addrs []string
	Port  int
	Algorithm Algorithm
}

type serverSelector interface {
	choose() *server
	after(*server)
}

type LoadBalancer struct {
	servers []*server
	port     int
	selector serverSelector
}

func newSelector(algo Algorithm, servers []*server) serverSelector {
	switch (algo) {
	case RoundRobin:
		return newRoundRobin(servers)
	case LeastConnections:
		return newLeastConnections(servers)
	}
	panic("invalid load balancing algorithm")
}

func NewLoadBalancer(args *LoadBalancerArgs) (*LoadBalancer, error) {
	serverLen := len(args.Addrs)
	servers := make([]*server, serverLen)
	for i, addr := range args.Addrs {
		server := newServer(addr)
		servers[i] = server
	}

	return &LoadBalancer{
		servers,
		args.Port,
		newSelector(args.Algorithm, servers),
	}, nil
}

func (lb *LoadBalancer) handler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		srvr := lb.selector.choose()
		if srvr == nil {
			w.WriteHeader(503)
		} else {
			srvr.proxy.ServeHTTP(w, r)
			lb.selector.after(srvr)
		}
	}
}

func (lb *LoadBalancer) printStatus() {
	for _, server := range lb.servers {
		println(server.host, ":",  server.alive)
	}
}

func (lb *LoadBalancer) startHealthCheck() func() {
	ticker := time.NewTicker(time.Second * 10)

	for {
		<-ticker.C
		for _, server := range lb.servers {
			alive := isAlive(server.host)
			if alive != server.alive {
				println("diff:", alive)
				println("server.alive before:", server.alive)
				server.mu.Lock()
				server.alive = alive
				println("server.alive before lock:", server.alive)
				server.mu.Unlock()
				println("server.alive after:", server.alive)
			}
		}
		lb.printStatus()
	}
}

func (lb *LoadBalancer) Start(algoStr string) {
	log.Printf("starting load balancer: port=%d algo=%s\n", lb.port, algoStr)
	lb.printStatus()
	http.HandleFunc("/", lb.handler())
	go lb.startHealthCheck()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", lb.port), nil))
}
