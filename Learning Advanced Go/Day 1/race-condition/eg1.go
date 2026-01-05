package main

import (
	"fmt"
	"sync"
)

var mu sync.Mutex

func main() {
	var count = 0
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// PROBLEM: 1000 workers are grabbing 'count' at the same time
			mu.Lock()
			count = count + 1
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Println("Final Count is:", count)
	// Prediction: It won't be 1000! It might be 942 or 980.
}
