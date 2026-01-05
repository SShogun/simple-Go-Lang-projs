package main

import (
	"fmt"
	"time"
)

func taxi(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Printf("Taxi %d: Started the job %d\n", id, j)
		time.Sleep(time.Second) // Simulate driving time
		fmt.Printf("Taxi %d: Finished job %d\n", id, j)
		results <- j * 2 // Send some "result" back
	}
}

func main() {
	const numJobs = 5
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// 2. Start 3 Workers (Taxis)
	for w := 1; w <= 3; w++ {
		go taxi(w, jobs, results)
	}

	// 3. Send 5 Jobs to the conveyor belt
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs) // Crucial: Tell workers no more jobs are coming!

	// 4. Collect results
	for a := 1; a <= numJobs; a++ {
		<-results
	}
	fmt.Println("All rides completed.")
}
