package main

import (
	"fmt"
	"time"
)

func PrimaryAPI(resp chan string) {
	time.Sleep(1500 * time.Millisecond)
	resp <- "Primary API Request Successful!"
}

func BackupAPI(resp chan string) {
	time.Sleep(500 * time.Millisecond)
	resp <- "Backup API Request Successful!"
}

func main() {
	resp := make(chan string)
	go PrimaryAPI(resp)
	go BackupAPI(resp)
	select {
	case msg := <-resp:
		fmt.Println("Received response from API: " + msg)
	case <-time.After(800 * time.Millisecond):
		fmt.Println("Error: Both API requests timed out. Please try again later.")
	}

}
