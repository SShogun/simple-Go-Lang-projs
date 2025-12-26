package main

import (
	"fmt"
	"sync"
)

var websites = []string{"google.com", "github.com", "stackoverflow.com", "go.dev", "medium.com"}

func worker(id int, wg *sync.WaitGroup, wb string, ch chan string) {
	defer wg.Done()
	ch <- fmt.Sprintf("%s website is up!!", wb)

}

func main() {
	var wg sync.WaitGroup
	ch := make(chan string)

	for i, gr := range websites {
		wg.Add(1)
		go worker(i, &wg, gr, ch)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for msg := range ch {
		println("Message: ", msg)
	}
}
