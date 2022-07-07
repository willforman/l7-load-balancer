package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type BenchmarkRequest struct {
	port string
	dur time.Duration
}

func runBenchmark(
	url string, 
	numReqs int, 
	period time.Duration, 
	output chan<- BenchmarkRequest,
) {
	for i := 1; i <= numReqs; i++ {
		go func(i int) {
			start := time.Now()
			resp, err := http.Get(url)
			delta := time.Now().Sub(start)
			if err != nil {
				return
			}
			defer resp.Body.Close()
			b, _ := io.ReadAll(resp.Body)
			output <- BenchmarkRequest{
				string(b),
				delta,
			}
		}(i)
		time.Sleep(period)
	}
}

func printPendingReqs(pendingCalls []int64, ports []string, ticker *time.Ticker, done <-chan bool) {
	for {
		select {
		case <-done:
			return
		case _ = <-ticker.C:
			for i, port := range ports {
				calls := pendingCalls[i]
				fmt.Printf("\033[K")
				fmt.Printf("%s %s\n", port, strings.Repeat("#", int(calls)))
			}
			fmt.Print(strings.Repeat("\033[A", len(ports)))
		}
	}
}

func handleResults(numReqs int, results <-chan BenchmarkRequest, pendingReqs []int64, ports []string) {
	total := int64(0)
	for i := 1; i <= numReqs; i++ {
		res := <-results
		total += res.dur.Microseconds()
	}
	mean := total / int64(numReqs)
	fmt.Printf("mean = %d.%dms\n", mean / 1000, mean % 1000)
}

