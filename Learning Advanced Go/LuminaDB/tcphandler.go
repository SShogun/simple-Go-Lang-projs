package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func handleClient(conn net.Conn, db *LuminaDB) {
	defer conn.Close()

	parser := NewRespParser(conn)

	for {
		args, err := parser.Parse()
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Client error: %v\n", err)
			}
			return
		}

		if len(args) == 0 {
			continue
		}

		command := strings.ToUpper(args[0])

		switch command {
		case "SET":
			if len(args) == 3 {
				err := db.Put(args[1], args[2])
				if err != nil {
					fmt.Printf("Error setting the key: %v\n", err)
					return
				}
				_, err = conn.Write([]byte("+OK\r\n"))
				if err != nil {
					fmt.Printf("Error writing to client: %v\n", err)
					return
				}
			} else {
				conn.Write([]byte("-ERR wrong number of arguments for 'SET'\r\n"))
			}
		case "GET":
			if len(args) == 2 {
				val, err := db.Get(args[1])
				if err != nil {
					fmt.Printf("Error getting the key: %v\n", err)
					return
				}
				if val == "" {
					conn.Write([]byte("$-1\r\n"))
				} else {
					response := fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)
					_, err = conn.Write([]byte(response))
					if err != nil {
						fmt.Printf("Error writing to client: %v\n", err)
						return
					}
				}
			} else {
				conn.Write([]byte("-ERR wrong number of arguments for 'GET'\r\n"))
			}
		case "DEL":
			if len(args) == 2 {
				err := db.Delete(args[1])
				if err != nil {
					fmt.Printf("Error deleting the key: %v\n", err)
					return
				}
				_, err = conn.Write([]byte(":1\r\n"))
				if err != nil {
					fmt.Printf("Error writing to client: %v\n", err)
					return
				}
			} else {
				conn.Write([]byte("-ERR wrong number of arguments for 'DEL'\r\n"))
			}
		case "PING":
			_, err = conn.Write([]byte("+PONG\r\n"))
			if err != nil {
				fmt.Printf("Error writing to client: %v\n", err)
				return
			}
		case "EXISTS":
			if len(args) == 2 {
				if db.Exists(args[1]) {
					_, err = conn.Write([]byte(":1\r\n"))
					if err != nil {
						fmt.Printf("Error writing to client: %v\n", err)
						return
					}
				} else {
					conn.Write([]byte(":0\r\n"))
				}
			} else {
				conn.Write([]byte("-ERR wrong number of arguments for 'EXISTS'\r\n"))
			}
		case "DBSIZE":
			if len(args) == 1 {
				size := db.Size()
				response := fmt.Sprintf(":%d\r\n", size)
				conn.Write([]byte(response))
			}
		case "FLUSHALL":
			if len(args) == 1 {
				db.FLUSHALL()
				conn.Write([]byte("+OK\r\n"))
			} else {
				conn.Write([]byte("-ERR wrong number of arguments for 'FLUSHALL'\r\n"))
			}
		default:
			conn.Write([]byte("-ERR unknown command\r\n"))
		}

	}
}
