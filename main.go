package main

import (
	"fmt"
	"os"
	"strconv"

	. "github.com/willforman/l7-load-balancer/loadbalancer"
)

func parseArgs(args []string) (*LoadBalancerArgs, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("num args provided < 2 [%d]", len(args))
	}
	portNum, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, fmt.Errorf("port cannot be parsed [%d]", portNum)
	}
	if portNum < 1024 || portNum > 65535 {
		return nil, fmt.Errorf("port out of range 1024 < p < 65535 [%d]", portNum)
	}
	hosts := args[1:]
	return &LoadBalancerArgs{
		Port:  args[0],
		Hosts: hosts,
	}, nil
}

func main() {
	args, err := parseArgs(os.Args[1:])
	if err != nil {
		panic(err)
	}
	lb, err := NewLoadBalancer(args)
	if err != nil {
		panic(err)
	}
	lb.Start()
}
