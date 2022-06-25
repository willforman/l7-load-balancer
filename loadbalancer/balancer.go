package loadbalancer

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type LoadBalancerArgs struct {
	Addrs []string
	Port  string
}

type LoadBalancer struct {
	serverRing serverRing
	port     string
}

func NewLoadBalancer(args *LoadBalancerArgs) (*LoadBalancer, error) {
	serverRing, err := newServerRing(args.Addrs)
	if err != nil {
		return nil, err
	}
	return &LoadBalancer{
		*serverRing,
		args.Port,
	}, nil
}

func (lb *LoadBalancer) handler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h := lb.serverRing.getAlive()
		if h != nil {
			h.proxy.ServeHTTP(w, r)
		} else {
			log.Println("no hosts alive")
		}
	}
}

func (lb *LoadBalancer) startWatchDog() func() {
	ticker := time.NewTicker(time.Second * 30)

	for {
		<-ticker.C
		lb.serverRing.doAll(func(server *server) {
			alive := isAlive(server.url.Host)
			if alive != server.alive {
				server.mu.Lock()
				server.alive = alive
				server.mu.Unlock()
			}
		})
	}
}

func (lb *LoadBalancer) Start() {
	log.Printf("starting load balancer on port %s\n", lb.port)
	http.HandleFunc("/", lb.handler())
	go lb.startWatchDog()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", lb.port), nil))
}
