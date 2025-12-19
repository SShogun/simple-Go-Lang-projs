package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// ValidateCategory checks if category contains only letters and spaces
func ValidateCategory(category string) error {
	for _, char := range category {
		if !unicode.IsLetter(char) && !unicode.IsSpace(char) {
			return fmt.Errorf("invalid category: category must contain only letters, not numbers or special characters")
		}
	}
	if len(category) == 0 {
		return fmt.Errorf("invalid category: category cannot be empty")
	}
	return nil
}

// ValidateAmount checks if amount is a positive number
func ValidateAmount(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("invalid amount: amount must be a positive number")
	}
	return nil
}

// readLine reads a full line from stdin with a prompt.
func readLine(prompt string, reader *bufio.Reader) (string, error) {
	fmt.Print(prompt)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

type expense struct {
	id          int
	amount      float64
	category    string
	description string
}

var curr_id int = 1
var expenses []expense

func main() {
	fmt.Println("Welcome to the Expense Tracker!\nLets begin!")

outerLoop:
	for {
		var choice string
		fmt.Println("\nAdd Expense\nView Expenses\nEdit Expense\nTotal Expenses\nDelete Expense\nExit")
		fmt.Print("Choose an option: ")
		fmt.Scanln(&choice)

		switch strings.ToLower(choice) {
		case "add":
			var exp1 expense
			fmt.Print("Enter amount: ")
			fmt.Scanln(&exp1.amount)

			// Validate amount
			if err := ValidateAmount(exp1.amount); err != nil {
				fmt.Println("Error:", err)
				continue
			}

			fmt.Print("Enter category: ")
			fmt.Scanln(&exp1.category)

			// Validate category
			if err := ValidateCategory(exp1.category); err != nil {
				fmt.Println("Error:", err)
				continue
			}

			fmt.Print("Enter description: ")
			fmt.Scanln(&exp1.description)
			exp1.id = curr_id
			curr_id++
			expenses = append(expenses, exp1)
			fmt.Println("Expense added successfully!")

		case "view":
			if len(expenses) != 0 {
				fmt.Println("Your previous expense was: ")
				for _, data := range expenses {
					fmt.Printf("%d. You spent: %.2f on %s (%s)\n", data.id, data.amount, data.category, data.description)
				}
			}
		case "total":
			if len(expenses) != 0 {

				var total float64
				for _, e := range expenses {
					total += e.amount
				}
				fmt.Printf("Total Expenses: %f\n", total)
			}
		case "exit":
			break outerLoop

		case "edit":
			if len(expenses) == 0 {
				fmt.Println("No expenses to edit")
				continue
			}

			reader := bufio.NewReader(os.Stdin)

			fmt.Println("Expenses:")
			for _, data := range expenses {
				fmt.Printf("%d. You spent: %.2f on %s (%s)\n", data.id, data.amount, data.category, data.description)
			}

			idStr, err := readLine("Enter the expense id: ", reader)
			if err != nil {
				fmt.Println("Error reading id")
				continue
			}
			ind, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Println("Invalid expense id")
				continue
			}

			idx := -1
			var current expense
			for i, e := range expenses {
				if e.id == ind {
					idx = i
					current = e
					break
				}
			}
			if idx == -1 {
				fmt.Println("Invalid expense id")
				continue
			}

			newexp := current

			amountStr, err := readLine(fmt.Sprintf("Enter amount (%.2f to keep): ", current.amount), reader)
			if err != nil {
				fmt.Println("Error reading amount")
				continue
			}
			if amountStr != "" {
				parsedAmt, err := strconv.ParseFloat(amountStr, 64)
				if err != nil {
					fmt.Println("Error: amount must be a number")
					continue
				}
				if err := ValidateAmount(parsedAmt); err != nil {
					fmt.Println("Error:", err)
					continue
				}
				newexp.amount = parsedAmt
			}

			categoryStr, err := readLine(fmt.Sprintf("Enter category (%s to keep): ", current.category), reader)
			if err != nil {
				fmt.Println("Error reading category")
				continue
			}
			if categoryStr != "" {
				if err := ValidateCategory(categoryStr); err != nil {
					fmt.Println("Error:", err)
					continue
				}
				newexp.category = categoryStr
			}

			descStr, err := readLine(fmt.Sprintf("Enter description (%s to keep): ", current.description), reader)
			if err != nil {
				fmt.Println("Error reading description")
				continue
			}
			if descStr != "" {
				newexp.description = descStr
			}

			expenses[idx] = newexp
			fmt.Println("Expense updated successfully!")
		case "delete":
			if len(expenses) == 0 {
				fmt.Println("No expenses to delete")
				continue
			}
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Expenses:")
			for _, data := range expenses {
				fmt.Printf("%d. You spent: %.2f on %s (%s)\n", data.id, data.amount, data.category, data.description)
			}

			idStr, err := readLine("Enter the expense id: ", reader)
			if err != nil {
				fmt.Println("Error reading id")
				continue
			}
			ind, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Println("Invalid expense id")
				continue
			}

			idx := -1
			for i, e := range expenses {
				if e.id == ind {
					idx = i
					break
				}
			}
			if idx == -1 {
				fmt.Println("Invalid expense id")
				continue
			}

			expenses = append(expenses[:idx], expenses[idx+1:]...)
			fmt.Println("Expense deleted successfully!")
		default:
			fmt.Println("Invalid choice. Please choose: add, view, edit, total, exit.")
		}

	}
}
