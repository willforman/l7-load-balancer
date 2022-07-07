package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"
)

func mean(nums []int64) float64 {
	total := int64(0)
	for _, num := range nums {
		total += num
	}
	return float64(total) / float64(len(nums))
}

func std(nums []int64, mean float64) float64 {
	total := float64(0)
	for _, num := range nums {
		total += math.Pow(float64(num) - mean, 2)
	}
	return total / float64(len(nums))
}

func printPendingReqs(pendingCalls []int64, ports []string, ticker *time.Ticker, done <-chan bool) {
	maxStd := -1.0
	outputLines := len(ports) + 1

	for {
		select {
		case <-done:
			fmt.Printf("\033[%dB", outputLines)
			fmt.Printf("max std = %f\n", maxStd)
			return
		case _ = <-ticker.C:
			for i, port := range ports {
				calls := pendingCalls[i]
				fmt.Printf("\033[K") // Clear the current line
				fmt.Printf("%s %s\n", port, strings.Repeat("#", int(calls)))
			}
			callsMean := mean(pendingCalls)
			callsStd := std(pendingCalls, callsMean)
			if callsStd > maxStd {
				maxStd = callsStd
			}
			fmt.Printf("mean = %f, std = %f\n", callsMean, callsStd)
			// Move cursor back to top so it can overwrite to show progress
			fmt.Printf("\033[%dA", outputLines)
		}
	}
}

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


func handleResults(numReqs int, results <-chan BenchmarkRequest, pendingReqs []int64, ports []string) {
	for i := 1; i <= numReqs; i++ {
		<-results
	}
}

