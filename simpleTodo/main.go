package main

import (
	"fmt"
)

type Todo struct {
	ID          int
	Description string
	Status      bool
}

var tasks []Todo
var nextID = 1

func main() {
	fmt.Println("Simple TODO app CLI")

	for {
		fmt.Println("\nCommands: add | list | done | delete | exit")
		fmt.Print("Enter command: ")
		var input string
		fmt.Scanln(&input)

		if input == "exit" {
			fmt.Println("Goodbye!")
			break
		} else if input == "list" {
			listTasks()

		} else if input == "done" {
			fmt.Print("Enter task ID: ")
			var id int
			fmt.Scanln(&id)
			markDone(id)

		} else if input == "delete" {
			fmt.Print("Enter task ID: ")
			var id int
			fmt.Scanln(&id)
			deleteTask(id)

		} else if input == "add" {
			fmt.Print("Enter task description: ")
			var desc string
			fmt.Scanln(&desc)
			addTask(desc)

		} else {
			fmt.Println("Unknown command")
		}
	}
}

func addTask(desc string) {
	task := Todo{ID: nextID, Description: desc, Status: false}
	tasks = append(tasks, task)
	fmt.Printf("Added task #%d: %s\n", nextID, desc)
	nextID++
}

func listTasks() {
	if len(tasks) == 0 {
		fmt.Println("No tasks yet!")
		return
	}
	fmt.Println("\n=== Your Tasks ===")
	for _, task := range tasks {
		status := "⬜"
		if task.Status {
			status = "✅"
		}
		fmt.Printf("[%d] %s %s\n", task.ID, status, task.Description)
	}
}

func markDone(id int) {
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Status = true
			fmt.Printf("Marked task #%d as done!\n", id)
			return
		}
	}
	fmt.Println("Task not found")
}

func deleteTask(id int) {
	for i := range tasks {
		if tasks[i].ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Printf("Deleted task #%d\n", id)
			return
		}
	}
	fmt.Println("Task not found")
}
