package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Logger struct {
	file *os.File
	mu   sync.Mutex
}

func (l *Logger) LogSet(key string, value string) error {
	return l.LogSetBinary(key, value)
}

func (d *RWData) Delete(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.value, key)
}
func NewLogger(filename string) (*Logger, error) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &Logger{file: f}, nil
}
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
	header := make([]byte, 17)

	for {
		_, err := io.ReadFull(file, header)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading header: %w", err)
		}

		action := header[0]

		keyLen := binary.BigEndian.Uint32(header[9:13])
		valLen := binary.BigEndian.Uint32(header[13:17])

		keyBytes := make([]byte, keyLen)
		_, err = io.ReadFull(file, keyBytes)
		if err != nil {
			return fmt.Errorf("error reading key: %w", err)
		}
		valBytes := make([]byte, valLen)
		_, err = io.ReadFull(file, valBytes)
		if err != nil {
			return fmt.Errorf("error reading value: %w", err)
		}
		key := string(keyBytes)
		value := string(valBytes)

		// Apply the action
		switch action {
		case 1: // SET
			db.store.Set(key, value)
		case 2: // DEL
			db.store.Delete(key)
		}
	}
	return nil
}
