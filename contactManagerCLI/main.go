package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Contact struct {
	Name  string `json:"name"`
	Phone int    `json:"phone"`
	Email string `json:"email"`
}

var Contacts []Contact

const contactsFile = "contacts.json"

func main() {
	fmt.Println("Welcome to the Contact Manager CLI")

	// Load existing contacts from file
	loadContacts()

	var choice string

Outerloop:
	for {
		fmt.Println("\nAdd | Edit| List | Search | Delete | Exit")
		fmt.Scanln(&choice)
		switch strings.ToLower(choice) {
		case "add":
			createContact(getContact())
			saveContacts()
		case "edit":
			editContact()
			saveContacts()
		case "list":
			listContact()
		case "search":
			searchContact()
		case "delete":
			deleteContact()
			saveContacts()
		case "exit":
			break Outerloop
		}
	}
}
func loadContacts() {
	data, err := os.ReadFile(contactsFile)
	if err != nil {
		// File doesn't exist yet, start with empty contacts
		if os.IsNotExist(err) {
			return
		}
		fmt.Println("Error loading contacts:", err)
		return
	}

	err = json.Unmarshal(data, &Contacts)
	if err != nil {
		fmt.Println("Error parsing contacts:", err)
	}
}

func saveContacts() {
	data, err := json.MarshalIndent(Contacts, "", "  ")
	if err != nil {
		fmt.Println("Error saving contacts:", err)
		return
	}

	err = os.WriteFile(contactsFile, data, 0644)
	if err != nil {
		fmt.Println("Error writing contacts file:", err)
	}
}

func deleteContact() {
	c := searchContact()
	Contacts = append(Contacts[:c.Phone], Contacts[c.Phone+1:]...)
}
func listContact() {
	for _, c := range Contacts {
		fmt.Printf("Name: %s, Phone: %d, Email: %s\n", c.Name, c.Phone, c.Email)
	}
}
func searchContact() Contact {
	var name string
	fmt.Println("Enter the Name of the Contact: ")
	fmt.Scanln(&name)
	for _, contact := range Contacts {
		if strings.EqualFold(name, contact.Name) {
			// EqualFold is used to ignore case sensitivity
			fmt.Printf("Name: %s, Phone: %d, Email: %s\n", contact.Name, contact.Phone, contact.Email)
		}
	}
	return Contact{}
}
func editContact() {
	c := searchContact()
	fmt.Println("Enter new details:")
	newC := getContact()
	c.Name = newC.Name
	c.Phone = newC.Phone
	c.Email = newC.Email

}
func getContact() Contact {
	var name, email string
	var phone int
	fmt.Println("Enter Name:")
	fmt.Scanln(&name)
	fmt.Println("Enter Phone:")
	fmt.Scanln(&phone)
	fmt.Println("Enter Email:")
	fmt.Scanln(&email)
	return Contact{Name: name, Phone: phone, Email: email}

}
func createContact(c Contact) {
	Contacts = append(Contacts, c)
}
