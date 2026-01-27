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

func main() {
	// 1. Initialize
	db, _ := NewLuminaDB("lumina.log")
	defer db.Close()

	// 2. RECOVER: This is where the magic happens
	fmt.Println("Recovering data from disk...")
	if err := db.Recover(); err != nil {
		fmt.Printf("Recovery failed: %v\n", err)
	}

	// 3. Check if old data exists
	val := db.Get("my_key")
	if val != "" {
		fmt.Printf("Found existing data: %s\n", val)
	} else {
		fmt.Println("No existing data. Creating new entry...")
		db.Put("my_key", "Hello_From_The_Past")
	}
}
