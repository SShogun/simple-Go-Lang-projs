package main

import (
	"fmt"
	"os"
	"sync"
)

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

func (db *LuminaDB) Size() int {
	db.store.mu.RLock()
	defer db.store.mu.RUnlock()

	return len(db.store.value)
}

func (db *LuminaDB) Exists(s string) bool {
	db.store.mu.RLock()
	defer db.store.mu.RUnlock()

	_, exists := db.store.value[s]
	return exists
}

func (db *LuminaDB) FLUSHALL() {
	db.store.mu.Lock()
	defer db.store.mu.Unlock()
	db.store.value = make(map[string]string)

	os.Truncate(db.logger.file.Name(), 0)
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

func (db *LuminaDB) Get(key string) (string, error) {
	return db.store.Get(key), nil
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

// func (l *Logger) LogSet(key, value string) error {
// 	l.mu.Lock()
// 	defer l.mu.Unlock()

//		_, err := l.file.WriteString("SET " + key + " " + value + "\n")
//		return err
//	}
