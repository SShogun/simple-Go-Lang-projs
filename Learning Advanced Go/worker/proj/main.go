package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Job struct {
	ID       int
	Filename string
}

func Jobs(id int, jobs <-chan Job, results chan<- string) {
	for j := range jobs {
		fmt.Printf("Worker: %d | ID: %d | Filename: %s is now working", id, j.ID, j.Filename)

		duration := time.Duration(rand.Intn(1000)) * time.Millisecond
		time.Sleep(duration)
		results <- fmt.Sprintf("Worker: %d | ID: %d | Filename: %s is now completed!", id, j.ID, j.Filename)
	}
}

func main() {
	jobs := make(chan Job, 10)
	results := make(chan string, 10)

	for i := 1; i <= 4; i++ {
		go Jobs(i, jobs, results)
	}

	for i := 1; i <= 10; i++ {
		jobs <- Job{ID: i, Filename: (fmt.Sprintf("photo_%d.jpg", i))}
	}

	close(jobs)
	fmt.Println("Main: Collecting results...")
	for r := 1; r <= 10; r++ {
		fmt.Println("Result:", <-results)
	}

	fmt.Println("All images processed successfully!")
}
