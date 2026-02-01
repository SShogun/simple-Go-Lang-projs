package main

import (
	"flag"
	"fmt"
	"net"
)

func main() {
	clientMode := flag.Bool("client", false, "run as client")
	flag.Parse()

	if *clientMode {
		runInteractiveClient()
		return
	}

	// Create a new LuminaDB instance
	db, err := NewLuminaDB("lumina.log")
	if err != nil {
		fmt.Println("Error creating database:", err)
		return
	}
	defer db.Close()

	// Recover from log
	if err := db.Recover(); err != nil {
		fmt.Println("Error recovering database:", err)
	}

	// Start TCP server
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("LuminaDB server listening on localhost:8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleClient(conn, db)
	}
}
