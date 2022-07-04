package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/willforman/l7-load-balancer/loadbalancer"
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

func cleanUp(srvrs []*http.Server, srvrsDone *sync.WaitGroup) {
	for _, srvr := range srvrs {
		srvr.Shutdown(context.Background())
	}

	srvrsDone.Wait()
	println("servers done")
}

func main() {
	srvrPorts := [3]string{"8081", "8082", "8083"}
	numSrvrs := len(srvrPorts)
	srvrs := make([]*http.Server, numSrvrs)
	serversDone := sync.WaitGroup{}
	serversDone.Add(numSrvrs)

	for i, port := range srvrPorts {
		srvr := startServer(port, &serversDone)
		srvrs[i] = srvr
	}

	urls := make([]url.URL, numSrvrs)
	for i, port := range srvrPorts {
		url, err := url.Parse(fmt.Sprintf("http://localhost:%s", port))
		if err != nil {
			panic(err)
		}
		urls[i] = *url
	}
	
	lbPort := 8080
	lbAddr := fmt.Sprintf("http://localhost:%d", lbPort)
	lbArgs := loadbalancer.LoadBalancerArgs{
		Port: lbPort,
		Urls: urls,
		Algorithm: loadbalancer.RoundRobin,
	}
	lb, err := loadbalancer.NewLoadBalancer(&lbArgs)
	if err != nil {
		panic(err)
	}
	go lb.Start("rr")

	numWorkers := 3
	out := make(chan string, numWorkers)

	go func() {
		for i := 1; i <= numWorkers; i++ {
			go func(i int) {
				resp, err := http.Get(lbAddr)
				if err != nil {
					return
				}
				defer resp.Body.Close()
				b, _ := io.ReadAll(resp.Body)
				out <- string(b)
			}(i)
			time.Sleep(time.Millisecond * time.Duration(500 * i))
		}
	}()

	for i := 0; i < numWorkers; i++ {
		println("recv: ", <-out)
	}

	close(out)

	cleanUp(srvrs, &serversDone)

	println("main: done. exiting...")
}