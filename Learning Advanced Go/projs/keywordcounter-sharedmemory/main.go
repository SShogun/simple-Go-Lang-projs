package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Job struct {
	Filepath string
	Keyword  string
}

type Result struct {
	Filepath string
	Count    int
}

func worker(jobs <-chan Job, wg *sync.WaitGroup, results chan<- Result) {
	defer wg.Done()
	for job := range jobs {
		count := 0
		fp, err := os.Open(job.Filepath)

		if err != nil {
			continue
		}
		defer fp.Close()
		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			line := scanner.Text()
			count += strings.Count(line, job.Keyword)
		}
		results <- Result{Filepath: job.Filepath, Count: count}
		fmt.Printf("Finished Processing %s\n", job.Filepath)
	}
}

func main() {
	files := []string{"app.log", "server.log", "db.log", "auth.log"}
	keyword := "ERROR"

	var wg sync.WaitGroup
	jobs := make(chan Job, len(files))
	results := make(chan Result, len(files))

	for _, file := range files {
		jobs <- Job{Filepath: file, Keyword: keyword}
	}

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go worker(jobs, &wg, results)
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	fmt.Println("\n--- Final Report ---")
	for r := range results {
		fmt.Printf("%s: %d occurrences\n", r.Filepath, r.Count)
	}
}
