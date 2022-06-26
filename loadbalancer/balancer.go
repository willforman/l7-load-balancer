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
	Port  string
	Algorithm Algorithm
}

type serverSelector interface {
	get([]server) *server
}

type LoadBalancer struct {
	servers []server
	port     string
	selector serverSelector
}

func NewLoadBalancer(args *LoadBalancerArgs) (*LoadBalancer, error) {
	serverLen := len(args.Addrs)
	servers := make([]server, serverLen)
	for i, addr := range args.Addrs {
		server, err := newServer(addr)
		if err != nil {
			panic(err)
		}
		servers[i] = *server
	}

	var selector serverSelector

	if args.Algorithm == RoundRobin {
		selector = &roundRobin{0, serverLen};
	}

	return &LoadBalancer{
		servers,
		args.Port,
		selector,
	}, nil
}

func (lb *LoadBalancer) handler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		server := lb.selector.get(lb.servers)
		if server != nil {
			server.proxy.ServeHTTP(w, r)
		} else {
			log.Println("no hosts alive")
		}
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

func (lb *LoadBalancer) Start() {
	log.Printf("starting load balancer on port %s\n", lb.port)
	http.HandleFunc("/", lb.handler())
	go lb.startWatchDog()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", lb.port), nil))
}
