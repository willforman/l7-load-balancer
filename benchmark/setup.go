package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	. "github.com/willforman/l7-load-balancer/loadbalancer"
)

func startServer(port string, wg *sync.WaitGroup, sleepDuration time.Duration) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(sleepDuration)
		io.WriteString(w, port)
	})
	srvr := &http.Server{ 
		Addr: ":" + port,
		Handler: mux,
	}
	go func() {
		err := srvr.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
		wg.Done()
	}()
	return srvr
}

func startServers(startPort int, numServers int, serversDone *sync.WaitGroup) ([]string, []*http.Server) {
	servers := make([]*http.Server, numServers)
	ports := make([]string, numServers) 
	for i := 0; i < numServers; i++ {
		ports[i] = strconv.Itoa(startPort + i)
		sleepDuration := time.Millisecond * time.Duration((i * 100))
		servers[i] = startServer(ports[i], serversDone, sleepDuration)
	}
	return ports, servers
}

func cleanUp(srvrs []*http.Server, srvrsDone *sync.WaitGroup) {
	for _, srvr := range srvrs {
		srvr.Shutdown(context.Background())
	}

	srvrsDone.Wait()
}

func startLb(lbPort int, urls []string, algoStr string) (string, *LoadBalancer){
	lbAddr := fmt.Sprintf("http://localhost:%d", lbPort)
	lb, err := NewLoadBalancer(lbPort, algoStr, urls)
	if err != nil {
		panic(err)
	}
	go lb.Start()
	return lbAddr, lb
}
