package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Job represents a file to be searched
type Job struct {
	Filepath string
	Keyword  string
}

// Result represents the findings from a single file
type Result struct {
	Filepath string
	Count    int
	Error    error
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		count, err := searchInFile(job.Filepath, job.Keyword)
		results <- Result{
			Filepath: job.Filepath,
			Count:    count,
			Error:    err,
		}
	}
}

// Helper function to handle File I/O safely
func searchInFile(path string, keyword string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close() // Properly closes once this file is done

	count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		count += strings.Count(scanner.Text(), keyword)
	}
	return count, scanner.Err()
}

func main() {
	// 1. Capture CLI Arguments
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <keyword> <path>")
		return
	}
	keyword := os.Args[1]
	rootPath := os.Args[2]

	start := time.Now()

	// 2. Setup Channels and WaitGroup
	jobs := make(chan Job, 100)
	results := make(chan Result, 100)
	var wg sync.WaitGroup

	// 3. Start Workers (Consumers)
	workerCount := 5
	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}

	// 4. File Discovery (Producer)
	go func() {
		err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Skip files we can't access
			}
			if !info.IsDir() {
				jobs <- Job{Filepath: path, Keyword: keyword}
			}
			return nil
		})
		if err != nil {
			fmt.Printf("Error walking path: %v\n", err)
		}
		close(jobs) // Signal workers that discovery is finished
	}()

	// 5. Orchestration (The Watcher)
	go func() {
		wg.Wait()
		close(results) // Signal the report loop that all work is done
	}()

	// 6. Report (The Collector)
	totalMatches := 0
	filesProcessed := 0

	fmt.Printf("Searching for '%s' in %s...\n\n", keyword, rootPath)

	for res := range results {
		filesProcessed++
		if res.Error != nil {
			continue // Silently skip error files for a clean CLI output
		}
		if res.Count > 0 {
			totalMatches += res.Count
			fmt.Printf("[MATCH] %d found in: %s\n", res.Count, res.Filepath)
		}
	}

	// Final Summary
	fmt.Println("---------------------------------------")
	fmt.Printf("Files Scanned:  %d\n", filesProcessed)
	fmt.Printf("Total Matches:  %d\n", totalMatches)
	fmt.Printf("Execution Time: %v\n", time.Since(start))
}
