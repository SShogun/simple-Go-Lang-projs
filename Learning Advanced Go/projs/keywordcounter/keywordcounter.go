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

var (
	results = make(map[string]int)
	mu      sync.Mutex
)

func worker(jobs <-chan Job, wg *sync.WaitGroup) {
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
		mu.Lock()
		results[job.Filepath] = count
		mu.Unlock()
		fmt.Printf("Finished Processing %s\n", job.Filepath)
	}
}

func main() {
	files := []string{"app.log", "server.log", "db.log", "auth.log"}
	keyword := "ERROR"

	var wg sync.WaitGroup
	jobs := make(chan Job, len(files))

	for _, file := range files {
		jobs <- Job{Filepath: file, Keyword: keyword}
	}

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go worker(jobs, &wg)
	}
	close(jobs)

	wg.Wait()
	fmt.Println("\n--- Final Report ---")
	for file, count := range results {
		fmt.Printf("%s: %d occurrences\n", file, count)
	}
}
