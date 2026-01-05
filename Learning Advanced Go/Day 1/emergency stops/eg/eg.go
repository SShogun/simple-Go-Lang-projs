package main

import (
	"fmt"
	"time"
)

func main() {
	resultChan := make(chan string)

	// Simulate a database query
	go func() {
		// Try changing this to 1 second or 3 seconds to see the difference
		time.Sleep(2 * time.Second)
		resultChan <- "Credit Score: 750"
	}()

	fmt.Println("Main: Fetching credit score...")

	// The select block picks whichever channel "speaks" first
	select {
	case res := <-resultChan:
		fmt.Println("Success:", res)
	case <-time.After(1500 * time.Millisecond): // 1.5 second deadline
		fmt.Println("Error: The bureau took too long. Aborting. 🚨")
	}
}
