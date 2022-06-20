package loadbalancer

import (
	"log"
	"net/http"
)

type LoadBalancerArgs struct {
	Hosts []string
	Port  int
}

type LoadBalancer struct {
	hostRing hostRing
	port     int
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
		h := lb.hostRing.get()
		h.proxy.ServeHTTP(w, r)
	}
}

func (lb *LoadBalancer) Start() {
	http.HandleFunc("/", lb.handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
