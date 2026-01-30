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
				db.Put(args[1], args[2])
				conn.Write([]byte("+OK\r\n"))
			} else {
				conn.Write([]byte("-ERR wrong number of arguments for 'SET'\r\n"))
			}
		case "GET":
			if len(args) == 2 {
				val := db.Get(args[1])
				if val == "" {
					conn.Write([]byte("$-1\r\n"))
				} else {
					response := fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)
					conn.Write([]byte(response))
				}
			}
		case "PING":
			conn.Write([]byte("+PONG\r\n"))
		default:
			conn.Write([]byte("-ERR unkown command\r\n"))
		}

	}
}
