package main

import (
	"io"
	"net/http"
	"time"
)

type BenchmarkRequest struct {
	port string
	dur time.Duration
}

func runBenchmark(url string, numWorkers int, period time.Duration, output chan<- BenchmarkRequest) {
	for i := 1; i <= numWorkers; i++ {
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
