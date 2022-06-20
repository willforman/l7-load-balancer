package main

import (
	"fmt"
	"os"
	"strconv"
)

type appArgs struct {
	port int
	hosts []string
}

func parseArgs(args []string) (*appArgs, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("num args provided < 2 [%d]", len(args))
	}
	port, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, fmt.Errorf("port cannot be parsed [%d]", port)
	}
	if port < 1024 || port > 65535 {
		return nil, fmt.Errorf("port out of range 1024 < p < 65535 [%d]", port)
	}
	hosts := args[1:]
	return &appArgs{
		port,
		hosts,
	}, nil
}

func main() {
	args, err := parseArgs(os.Args[1:])
	if err != nil {
		panic(err)
	}
	lb, err := newLoadBalancer(args)
	if err != nil {
		panic(err)
	}
	startServer(lb)
}
