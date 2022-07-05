package main

import (
	"flag"
	"sync"
	"time"
)


func main() {
	startPort := flag.Int("startPort", 8080, "port to start using")
	numServers := flag.Int("numServers", 3, "number of servers to create")
	numReqs := flag.Int("numReqs", 10, "number of requests to make")
	reqPeriodMs := flag.Int("reqPeriod", 500, "milliseconds to wait between requests")
	flag.Parse()

	var serversDone sync.WaitGroup
	serversDone.Add(*numServers)

	urls, servers := startServers(*startPort, *numServers, &serversDone)

	lbUrl, lb := startLb(*startPort + *numServers, urls)

	out := make(chan BenchmarkRequest, *numReqs)
	reqPeriod := time.Millisecond * time.Duration(*reqPeriodMs)

	go runBenchmark(lbUrl, *numReqs, reqPeriod, out)

	for i := 0; i < *numReqs; i++ {
		println("recv:", (<-out).port)
	}

	close(out)

	cleanUp(servers, &serversDone)
	err := lb.Stop()
	if err != nil {
		panic(err)
	}

	println("gracefully shutdown. exiting...")
}
