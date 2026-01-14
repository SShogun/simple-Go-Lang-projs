package main

import (
	"fmt"
	"sync"
	"time"
)

type Result struct {
	URL     string
	Status  string
	Latency time.Duration
}

func worker(id int, jobs <-chan string, results chan<- Result, wg *sync.WaitGroup, latency time.Duration) {
	defer wg.Done()
	for url := range jobs {
		fmt.Printf("Worker %d processing %s\n", id, url)
		time.Sleep(latency)
		results <- Result{URL: url, Status: "200 OK", Latency: latency}
	}
}

func main() {
	var wg sync.WaitGroup
	jobs := make(chan string, 100)
	results := make(chan Result, 100)

	for w := 1; w <= 5; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg, 100*time.Millisecond)
	}

	for i := 0; i < 20; i++ {
		jobs <- fmt.Sprintf("http://site-%d.com", i)
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		fmt.Println("Result received: ", res.URL, res.Status)
	}
}
