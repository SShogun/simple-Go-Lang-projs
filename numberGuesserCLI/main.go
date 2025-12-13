package main

import (
	"fmt"
	"math/rand"
)

func randomNumberGen(min, max int) int {
	rnd := rand.Intn(max-min) + min
	return rnd
}

func checkGuess(guess, targent int) bool {
	return guess == targent
}

func main() {
	fmt.Println("Welcome to the Number Guesser CLI")
	fmt.Println("Enter 2 numbers to define the range")
	min, max := 0, 0
	fmt.Scanln(&min)
	fmt.Scanln(&max)
	g := randomNumberGen(min, max)

	var usInp int
	var count int = 0
	fmt.Println("Input a number")
	fmt.Scanln(&usInp)

	if checkGuess(usInp, g) {
		fmt.Println("You guessed it!")
	} else {
		for usInp != g {
			if usInp < g {
				fmt.Println("Too low, try again")
				count++
			} else {
				fmt.Println("Too high, try again")
				count++
			}
			fmt.Scanln(&usInp)
		}
		fmt.Println("You guessed it! in ", count, "tries")
	}
}
