package main

import (
	"os"
	"sync"
)

type Logger struct {
	file *os.File
	mu   sync.Mutex
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
