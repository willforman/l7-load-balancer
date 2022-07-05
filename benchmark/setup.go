package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"

	. "github.com/willforman/l7-load-balancer/loadbalancer"
)

func startServer(port string, wg *sync.WaitGroup) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
	urls := make([]string, numServers) 
	for i := 0; i < numServers; i++ {
		portStr := strconv.Itoa(startPort + i)
		urls[i] = fmt.Sprintf("http://localhost:%s", portStr)

		servers[i] = startServer(portStr, serversDone)
	}
	return urls, servers
}

func cleanUp(srvrs []*http.Server, srvrsDone *sync.WaitGroup) {
	for _, srvr := range srvrs {
		srvr.Shutdown(context.Background())
	}

	srvrsDone.Wait()
}

func startLb(lbPort int, urls []string) (string, *LoadBalancer){
	lbAddr := fmt.Sprintf("http://localhost:%d", lbPort)
	lb, err := NewLoadBalancer(lbPort, "rr", urls)
	if err != nil {
		panic(err)
	}
	go lb.Start()
	return lbAddr, lb
}