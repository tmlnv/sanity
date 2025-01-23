package logger

import (
	"io"
	"log"
	"os"
	"sync"
)

var (
	instance *log.Logger
	mu       sync.Mutex
)

// Init initializes the logger with the specified log file.
func Init(logFile string) {
	mu.Lock()
	defer mu.Unlock()

	writers := []io.Writer{os.Stdout}
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Could not open log file: %v", err)
		}
		writers = append(writers, file)
	}

	instance = log.New(io.MultiWriter(writers...), "[sanity] ", log.LstdFlags|log.Lmsgprefix)
}

// Info logs an informational message.
func Info(msg string, args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		panic("logger not initialized")
	}
	instance.Printf("[INFO] "+msg, args...)
}

// Error logs an error message.
func Error(msg string, args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		panic("logger not initialized")
	}
	instance.Printf("[ERROR] "+msg, args...)
}
