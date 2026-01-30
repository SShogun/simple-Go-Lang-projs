package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func testClient() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Test SET command
	fmt.Println("\n=== Testing SET command ===")
	fmt.Fprintf(conn, "*3\r\n$3\r\nSET\r\n$4\r\nname\r\n$6\r\ngemini\r\n")
	response, _ := reader.ReadString('\n')
	fmt.Printf("Server response: %s", response)

	// Test GET command
	fmt.Println("\n=== Testing GET command ===")
	fmt.Fprintf(conn, "*2\r\n$3\r\nGET\r\n$4\r\nname\r\n")
	response, _ = reader.ReadString('\n')
	fmt.Printf("Server response: %s", response)
	if strings.Contains(response, "$") {
		value, _ := reader.ReadString('\n')
		fmt.Printf("Value: %s", value)
	}

	// Test PING command
	fmt.Println("\n=== Testing PING command ===")
	fmt.Fprintf(conn, "*1\r\n$4\r\nPING\r\n")
	response, _ = reader.ReadString('\n')
	fmt.Printf("Server response: %s", response)

	fmt.Println("\nTest completed!")
}
