package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Logger struct {
	file *os.File
	mu   sync.Mutex
}

// Memory Store
type RWData struct {
	mu    sync.RWMutex
	value map[string]string
}

func (d *RWData) Get(key string) string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.value[key]
}

func (d *RWData) Set(key, value string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.value[key] = value
}

func (d *RWData) Delete(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.value, key)
}

// Main Engine
type LuminaDB struct {
	store  *RWData
	logger *Logger
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

func (l *Logger) LogSet(key, value string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	_, err := l.file.WriteString("SET " + key + " " + value + "\n")
	return err
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
	db, err := NewLuminaDB("lumina.log")
	if err != nil {
		fmt.Printf("Initialization failed: %v\n", err)
		return
	}
	defer db.Close()

	fmt.Println("LuminaDB Engine Started.")
	fmt.Print("Checking for recovery data... ")
	if err := db.Recover(); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Done.")
	}

	startShell(db)
}
