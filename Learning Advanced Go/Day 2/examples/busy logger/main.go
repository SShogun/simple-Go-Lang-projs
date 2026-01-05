package main

import "fmt"

func main() {
	// A channel with a buffer of only 1
	logChan := make(chan string, 1)

	// 1. Fill the buffer
	logChan <- "Log 1: System started"

	// 2. Try to send another log without waiting
	select {
	case logChan <- "Log 2: System running":
		fmt.Println("Log 2 sent successfully")
	default:
		// This runs instantly because logChan is full
		fmt.Println("Warning: Logger busy. Log 2 dropped to maintain performance.")
	}

	// 3. Try to receive without waiting
	select {
	case msg := <-logChan:
		fmt.Println("Received:", msg)
	default:
		fmt.Println("No logs found.")
	}
}

// question: what happens to the 2nd log message when the logger is busy?
// answer: The 2nd log message is dropped and a warning is printed indicating that the logger is busy.
// so when does the log message finally get logged?
// The log message does not get logged at all; it is discarded to maintain performance.
