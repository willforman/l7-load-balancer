package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
)


func main() {
	startPort := flag.Int("startPort", 8080, "port to start using")
	numServers := flag.Int("numServers", 3, "number of servers to create")
	numReqs := flag.Int("numReqs", 10, "number of requests to make")
	reqPeriodMs := flag.Int("reqPeriod", 500, "milliseconds to wait between requests")
	algoStr := flag.String("algo", "lc", "Load balancing algorithm (either LeastConnections or RoundRobin)")
	flag.Parse()

	var serversDone sync.WaitGroup
	serversDone.Add(*numServers)

	ports, servers := startServers(*startPort, *numServers, &serversDone)
	urls := make([]string, *numServers)
	for i, port := range ports {
		urls[i] = fmt.Sprintf("http://localhost:%s", port)
	}

	lbUrl, lb := startLb(*startPort + *numServers, urls, *algoStr)

	out := make(chan BenchmarkRequest, *numReqs)
	reqPeriod := time.Millisecond * time.Duration(*reqPeriodMs)

	go runBenchmark(lbUrl, *numReqs, reqPeriod, out)

	handleResults(*numReqs, ports, out)

	close(out)

	cleanUp(servers, &serversDone)
	err := lb.Stop()
	if err != nil {
		panic(err)
	}

	println("gracefully shutdown. exiting...")
}
