package main

import (
	"fmt"
	"strings"
	"unicode"
)

func analyzePassword(password string) int {
	score := 0

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	const specialChars = "!@#$%^&*()-_=+[]{}|;:'\",.<>?/`~"

	if len(password) >= 8 {
		score++
		fmt.Println("The password has a length of greater than 8")
	} else {
		fmt.Println("The password does not have a length of greater than 8")
	}

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if unicode.IsLower(char) {
			hasLower = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
		if strings.ContainsRune(specialChars, char) {
			hasSpecial = true
		}
	}
	if hasUpper {
		score++
		fmt.Println("Has upper case character")
	} else {
		fmt.Println("Doesnt have upper case character")
	}

	if hasLower {
		score++
		fmt.Println("Has lower case character")
	} else {
		fmt.Println("Doesnt have lower case character")
	}

	if hasDigit {
		score++
		fmt.Println("Has digit character")
	} else {
		fmt.Println("Doesnt have digit character")
	}

	if hasSpecial {
		score++
		fmt.Println("Has special character")
	} else {
		fmt.Println("Doesnt have special character")
	}

	return score
}

func main() {
	var password string
	fmt.Print("Enter a password to analyze: ")

	_, err := fmt.Scanln(&password)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	strength := analyzePassword(password)
	fmt.Printf("Password strength score: %d out of 5\n", strength)
}
