package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var wg sync.WaitGroup

	var totalRequests atomic.Uint64

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			totalRequests.Add(1)
		}()
	}
	wg.Wait()
	fmt.Println("Total Requests:", totalRequests.Load())
}
