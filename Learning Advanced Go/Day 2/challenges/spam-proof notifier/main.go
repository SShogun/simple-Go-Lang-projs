package main

import (
	"fmt"
	"time"
)

func worker(emailChan chan int) {
	for emailID := range emailChan {
		fmt.Printf("Processing email %d...\n", emailID)
		time.Sleep(2 * time.Second)
		fmt.Printf("Email %d processed.\n", emailID)
	}
}

func main() {
	emailChan := make(chan int, 3)

	// Start the worker goroutine
	go worker(emailChan)

	// Try to send 10 email requests
	for i := 1; i <= 10; i++ {
		select {
		case emailChan <- i:
			fmt.Printf("Email %d queued.\n", i)
		default:
			fmt.Printf("System overloaded. Email %d skipped to prevent lag.\n", i)
		}
		time.Sleep(500 * time.Millisecond) // Small delay between attempts
	}

	// Close the channel and wait for worker to finish
	close(emailChan)
	time.Sleep(15 * time.Second) // Give worker time to process remaining emails
}
