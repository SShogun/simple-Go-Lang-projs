package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"sync"
)

type Job struct {
	Filepath string
}

type Result struct {
	Filepath   string
	HashString string
}

func worker(jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		f, _ := os.Open(job.Filepath)
		h := sha256.New()
		io.Copy(h, f)
		hashSum := fmt.Sprintf("%x", h.Sum(nil))
		results <- Result{Filepath: job.Filepath, HashString: hashSum}
	}
}

func main() {
	jobs := make(chan Job, 100)
	results := make(chan Result, 100)
	var wg sync.WaitGroup

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}
	go func() {
		wg.Wait()
		close(jobs)
	}()

	fmt.Println("The outputs are:")
	for r := range results {
		fmt.Printf("Path: %s | HashSum: %s\n", r.Filepath, r.HashString)
	}
}
