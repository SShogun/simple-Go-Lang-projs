package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	file *os.File
	mu   sync.Mutex
}

func (l *Logger) LogSet(key string, value string) error {
	panic("unimplemented")
}

func (d *RWData) Delete(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.value, key)
}

func NewLuminaDB(logFile string) (*LuminaDB, error) {
	l, err := NewLogger(logFile)
	if err != nil {
		return nil, err
	}

	return &LuminaDB{store: &RWData{value: make(map[string]string)}, logger: l}, nil
}

func (db *LuminaDB) Put(key, value string) error {
	if err := db.logger.LogSet(key, value); err != nil {
		return fmt.Errorf("Failed to log to disk: %w", err)
	}
	db.store.Set(key, value)
	return nil
}

func (db *LuminaDB) Get(key string) string {
	return db.store.Get(key)
}

// Delete removes from Disk then Memory
func (db *LuminaDB) Delete(key string) error {
	if err := db.logger.LogDelete(key); err != nil {
		return fmt.Errorf("failed to log delete: %w", err)
	}

	db.store.Delete(key)
	return nil
}

func (db *LuminaDB) Close() error {
	return db.logger.Close()
}
func NewLogger(filename string) (*Logger, error) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &Logger{file: f}, nil
}

// func (l *Logger) LogSet(key, value string) error {
// 	l.mu.Lock()
// 	defer l.mu.Unlock()

//		_, err := l.file.WriteString("SET " + key + " " + value + "\n")
//		return err
//	}
func (l *Logger) LogSetBinary(key, value string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 1. Build the frame (Action 1 for SET)
	timestamp := time.Now().Unix()
	keyBytes := []byte(key)
	valBytes := []byte(value)

	// Frame:
	size := 1 + 8 + 4 + 4 + len(keyBytes) + len(valBytes)
	buf := make([]byte, size)

	buf[0] = 1 // SET
	binary.BigEndian.PutUint64(buf[1:9], uint64(timestamp))
	binary.BigEndian.PutUint32(buf[9:13], uint32(len(keyBytes)))
	binary.BigEndian.PutUint32(buf[13:17], uint32(len(valBytes)))
	copy(buf[17:], keyBytes)
	copy(buf[17+len(keyBytes):], valBytes)

	// 2. Write the binary blob
	_, err := l.file.Write(buf)
	return err
}
func (l *Logger) EncodeFrame(action byte, timestamp int64, key, value string) ([]byte, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	keyLen := uint32(len(key))
	valueLen := uint32(len(value))
	size := 1 + 8 + 4 + 4 + len(key) + len(value)
	data := make([]byte, size)

	data[0] = action
	binary.BigEndian.PutUint64(data[1:9], uint64(timestamp))
	binary.BigEndian.PutUint32(data[9:13], keyLen)
	binary.BigEndian.PutUint32(data[13:17], valueLen)
	copy(data[17:17+len(key)], []byte(key))
	copy(data[17+len(key):], []byte(value))
	return data, nil
}

func (l *Logger) LogDelete(key string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	_, err := l.file.WriteString("DEL " + key + "\n")
	return err
}

func (l *Logger) Close() error {
	return l.file.Close()
}

// Recovery function to be implemented
func (db *LuminaDB) Recover() error {
	file, err := os.Open("lumina.log")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")

		if len(parts) < 2 {
			continue
		}

		command := parts[0]
		key := parts[1]

		switch command {
		case "SET":
			if len(parts) == 3 {
				value := parts[2]

				db.store.Set(key, value)
			}
		case "DEL":
			db.store.Delete(key)
		}
	}
	return scanner.Err()
}

// CLI loop
func startShell(db *LuminaDB) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("\n--- LuminaDB Shell ---")
	fmt.Println("Commands: SET <key> <val> | GET <key> | DEL <key> | EXIT")
	fmt.Print("> ")

	for scanner.Scan() {
		input := scanner.Text()
		parts := strings.Fields(input) // Splits by any whitespace

		if len(parts) == 0 {
			fmt.Print("> ")
			continue
		}

		command := strings.ToUpper(parts[0])

		switch command {
		case "SET":
			if len(parts) != 3 {
				fmt.Println("Error: SET requires key and value")
			} else {
				db.Put(parts[1], parts[2])
				fmt.Println("OK")
			}

		case "GET":
			if len(parts) != 2 {
				fmt.Println("Error: GET requires key")
			} else {
				val := db.Get(parts[1])
				if val == "" {
					fmt.Println("(nil)")
				} else {
					fmt.Println(val)
				}
			}

		case "DEL":
			if len(parts) != 2 {
				fmt.Println("Error: DEL requires key")
			} else {
				db.Delete(parts[1])
				fmt.Println("OK")
			}

		case "EXIT":
			fmt.Println("Shutting down LuminaDB...")
			return

		default:
			fmt.Printf("Unknown command: %s\n", command)
		}

		fmt.Print("> ")
	}
}

func main() {
	clientMode := flag.Bool("client", false, "run as client")
	flag.Parse()

	if *clientMode {
		testClient()
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

func handle(conn net.Conn, db *LuminaDB) {
	defer conn.Close()
	parser := RespParser{reader: bufio.NewReader(conn)}
	for {
		args, err := parser.Parse()
		if err != nil {
			return
		}

		cmd := strings.ToUpper(args[0])
		switch cmd {
		case "SET":
			db.logger.LogSet(args[1], args[2])
			db.store.Set(args[1], args[2])
			conn.Write([]byte("+OK\r\n"))
		case "GET":
			val := db.store.Get(args[1])
			if val == "" {
				conn.Write([]byte("$-1\r\n"))
			} else {
				conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)))
			}
		}
	}
}

// (Supporting structs RWData and LuminaDB omitted for brevity but remain the same)
type RWData struct {
	mu    sync.RWMutex
	value map[string]string
}

func (d *RWData) Get(k string) string { d.mu.RLock(); defer d.mu.RUnlock(); return d.value[k] }
func (d *RWData) Set(k, v string)     { d.mu.Lock(); defer d.mu.Unlock(); d.value[k] = v }

type LuminaDB struct {
	store  *RWData
	logger *Logger
}
