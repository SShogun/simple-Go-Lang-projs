package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var websites = []string{"google.com", "github.com", "stackoverflow.com", "go.dev", "medium.com"}

func worker(id int, wg *sync.WaitGroup, wb string, ch chan string) {
	defer wg.Done()

	fmt.Printf("Checking %s....\n", wb)
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	ch <- fmt.Sprintf("%s website is up!!", wb)
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan string)
	for i, gr := range websites {
		wg.Add(1)
		go worker(i, &wg, gr, ch)
		received := <-ch
		fmt.Println("Message: ", received)
	}

	fmt.Println("Waiting for workers to finish....")
	wg.Wait()
	fmt.Printf("All workers finished. Exiting")
}
