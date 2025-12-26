package main

import "fmt"

func greet(ch chan string) {
	ch <- "Hello, World!"
}

func main() {
	ch := make(chan string)
	go greet(ch)
	message := <-ch
	fmt.Println(message)
}
