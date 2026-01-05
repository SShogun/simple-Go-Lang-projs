package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)

	go func() {
		time.Sleep(3 * time.Second) // Simulate a slow response
		ch <- "Data Received"
	}()

	select {
	case res := <-ch:
		fmt.Println(res)
	case <-time.After(1 * time.Second): // "I'm giving you 1 second"
		fmt.Println("Error: Network timeout. Please try again later.")
	}
}
