package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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

// Info logs an informational message with key-value pairs.
func Info(msg string, args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		panic("logger not initialized")
	}

	// Format key-value pairs
	var b strings.Builder
	b.WriteString("[INFO] ")
	b.WriteString(msg)
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			b.WriteString(fmt.Sprintf(" %s=%v", args[i], args[i+1]))
		}
	}

	instance.Println(b.String())
}

// Error logs an error message with key-value pairs.
func Error(msg string, args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		panic("logger not initialized")
	}

	// Format key-value pairs
	var b strings.Builder
	b.WriteString("[ERROR] ")
	b.WriteString(msg)
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			b.WriteString(fmt.Sprintf(" %s=%v", args[i], args[i+1]))
		}
	}

	instance.Println(b.String())
}
