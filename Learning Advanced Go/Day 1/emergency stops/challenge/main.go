package main

import (
	"fmt"
	"time"
)

// Fast Bank: sends "Success" after 500ms
func fastBank(ch chan string) {
	time.Sleep(500 * time.Millisecond)
	ch <- "Success"
}

// Slow Bank: sends "Success" after 5 seconds
func slowBank(ch chan string) {
	time.Sleep(5 * time.Second)
	ch <- "Success"
}

func main() {
	fastCh := make(chan string)
	slowCh := make(chan string)

	// Launch both banks as goroutines
	go fastBank(fastCh)
	go slowBank(slowCh)

	// Use select to capture the first bank that finishes
	select {
	case msg := <-fastCh:
		fmt.Println("Fast Bank:", msg)
	case msg := <-slowCh:
		fmt.Println("Slow Bank:", msg)
	case <-time.After(2 * time.Second):
		fmt.Println("Payment Gateway Timeout: Reverting Transaction.")
	}
}