package main

import (
	"flag"
	"fmt"
	"log"

	. "github.com/willforman/l7-load-balancer/loadbalancer"
)

func main() {
	port := flag.Int("port", 8080, "Port num")
	algoStr := flag.String("algo", "lc", "Load balancing algorithm (either LeastConnections or RoundRobin)")
	flag.Parse()

	lb, err := NewLoadBalancer(*port, *algoStr, flag.Args())
	if err != nil {
		panic(fmt.Errorf("NewLoadBalancer: %w", err))
	}
	log.Printf("starting load balancer: port=%d algo=%s\n", *port, *algoStr)
	lb.Start()
}
