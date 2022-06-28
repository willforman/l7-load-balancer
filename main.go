package main

import (
	"flag"
	"fmt"
	"net/url"

	. "github.com/willforman/l7-load-balancer/loadbalancer"
)

func parseArgs(port int, algoStr string, addrs []string) (*LoadBalancerArgs, error) {
	if port < 1024 || port > 65535 {
		return nil, fmt.Errorf("port out of range 1024 < p < 65535 [%d]", port)
	}

	var algo Algorithm
	switch (algoStr) {
	case "lc":
		algo = LeastConnections
	case "rr":
		algo = RoundRobin
	default:
		return nil, fmt.Errorf("algo choice not lc or rr [%s]", algoStr)
	}

	numAddrs := len(addrs)
	if numAddrs == 0 {
		return nil, fmt.Errorf("no provided addrs as arguments")
	}
	urls := make([]url.URL, numAddrs)
	for _, addr := range addrs {
		url, err := url.Parse(addr)
		if err != nil {
			return nil, err
		}
		if url.Scheme != "http" {
			return nil, fmt.Errorf("scheme not http: [%s -> %s]", addr, url.Scheme)
		}
		urls = append(urls, *url)
	}
		

	return &LoadBalancerArgs{
		Port:  port,
		Urls: urls,
		Algorithm: algo,
	}, nil
}

func main() {
	port := flag.Int("port", 8080, "Port num")
	algo := flag.String("algo", "lc", "Load balancing algorithm (either LeastConnections or RoundRobin)")
	flag.Parse()

	args, err := parseArgs(*port, *algo, flag.Args())
	if err != nil {
		panic(err)
	}
	lb, err := NewLoadBalancer(args)
	if err != nil {
		panic(err)
	}
	lb.Start(*algo)
}
