package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Client struct {
	conn   net.Conn
	reader *bufio.Reader
}

func NewClient(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:   conn,
		reader: bufio.NewReader(conn),
	}, nil
}

func (c *Client) SendCommand(args []string) error {
	// Build RESP array
	fmt.Fprintf(c.conn, "*%d\r\n", len(args))
	for _, arg := range args {
		fmt.Fprintf(c.conn, "$%d\r\n%s\r\n", len(arg), arg)
	}
	return nil
}

func (c *Client) ReadResponse() string {
	line, err := c.reader.ReadString('\n')
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	line = strings.TrimSpace(line)

	switch line[0] {
	case '+':
		return line[1:]
	case '-':
		return "(error) " + line[1:]
	case ':':
		return line[1:]
	case '$':
		if line == "$-1" {
			return "(nil)"
		}
		value, _ := c.reader.ReadString('\n')
		return strings.TrimSpace(value)
	default:
		return line
	}
}

func (c *Client) Close() {
	c.conn.Close()
}

func runInteractiveClient() {
	fmt.Println("╔═══════════════════════════════════════╗")
	fmt.Println("║      LuminaDB Interactive Client      ║")
	fmt.Println("╚═══════════════════════════════════════╝")
	fmt.Println()

	client, err := NewClient("localhost:8080")
	if err != nil {
		fmt.Printf("❌ Could not connect to server: %v\n", err)
		fmt.Println("\nMake sure the server is running with: go run .")
		return
	}
	defer client.Close()

	fmt.Println("✓ Connected to localhost:8080")
	fmt.Println("\nCommands: SET, GET, DEL, PING, EXISTS, DBSIZE, EXIT")
	fmt.Println("Example: SET mykey myvalue")
	fmt.Println("─────────────────────────────────────────")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("luminadb> ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		args := strings.Fields(input)
		cmd := strings.ToUpper(args[0])

		if cmd == "EXIT" || cmd == "QUIT" {
			fmt.Println("Goodbye!")
			break
		}

		err := client.SendCommand(args)
		if err != nil {
			fmt.Printf("Error sending command: %v\n", err)
			continue
		}

		response := client.ReadResponse()
		fmt.Println(response)
	}
}
