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

func runBenchmark(url string, numReqs int, period time.Duration, output chan<- BenchmarkRequest) {
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

func handleResults(numReqs int, ports []string, results <-chan BenchmarkRequest) {
	counts := make(map[string]int, len(ports))
	for _, port := range ports {
		counts[port] = 0
	}
	total := int64(0)

	for i := 1; i <= numReqs; i++ {
		res := <-results
		total += res.dur.Microseconds()
		counts[res.port]++

		for port, count := range counts {
			fmt.Printf("%s %s\n", port, strings.Repeat("#", count))
		}
		// Move cursor up 
		if i != numReqs {
			fmt.Print(strings.Repeat("\033[A", len(ports)))
		}
	}
	mean := total / int64(numReqs)
	fmt.Printf("mean = %d.%dms\n", mean / 1000, mean % 1000)
}
