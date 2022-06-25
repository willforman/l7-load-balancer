package loadbalancer

import (
	"log"
	"net/http"
	"time"
)

type LoadBalancerArgs struct {
	Hosts []string
	Port  string
}

type LoadBalancer struct {
	hostRing hostRing
	port     string
}

func NewLoadBalancer(args *LoadBalancerArgs) (*LoadBalancer, error) {
	hr, err := newHostRing(args.Hosts)
	if err != nil {
		return nil, err
	}
	return &LoadBalancer{
		*hr,
		args.Port,
	}, nil
}

func (lb *LoadBalancer) handler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h := lb.hostRing.getAlive()
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
		lb.hostRing.doAll(func(host *host) {
			alive := isAlive(host.url.Host)
			if alive != host.alive {
				host.mu.Lock()
				host.alive = alive
				host.mu.Unlock()
			}
		})
	}
}

func (lb *LoadBalancer) Start() {
	log.Printf("starting load balancer on port %s\n", lb.port)
	http.HandleFunc("/", lb.handler())
	go lb.startWatchDog()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
