package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type LoadBalancer struct {
	hostRing HostRing
	port int
}

func newLoadBalancer(args *appArgs) (*LoadBalancer, error) {
	hr, err := newHostRing(args.hosts)
	if err != nil {
		return nil, err
	}
	return &LoadBalancer{
		*hr,
		args.port,
	}, nil
}

func newProxy(host string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	return httputil.NewSingleHostReverseProxy(url), nil
}

func (lb *LoadBalancer) Handler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h := lb.hostRing.get()
		println("Host chosen: " + h.addr)
        h.proxy.ServeHTTP(w, r)
    }
}

func startServer(lb *LoadBalancer) {
	http.HandleFunc("/", lb.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
