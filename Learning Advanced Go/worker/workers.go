package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var websites = []string{"google.com", "github.com", "stackoverflow.com", "go.dev", "medium.com"}

func worker(id int, wg *sync.WaitGroup, wb string) {
	defer wg.Done()

	fmt.Printf("Checking %s....\n", wb)
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	fmt.Printf("%s website is up!!\n", wb)
}

func main() {
	var wg sync.WaitGroup

	for i, gr := range websites {
		wg.Add(1)
		go worker(i, &wg, gr)

	}

	fmt.Println("Waiting for workers to finish....")
	wg.Wait()
	fmt.Printf("All workers finished. Exiting")
}
