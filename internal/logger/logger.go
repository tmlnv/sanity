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

func Init(logFile string) {
	mu.Lock()
	defer mu.Unlock()

	writers := []io.Writer{os.Stdout}
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			writers = append(writers, file)
		}
	}

	instance = log.New(io.MultiWriter(writers...), "[sanity] ", log.LstdFlags|log.Lmsgprefix)
}

func Info(msg string, args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	instance.Printf("[INFO] "+msg, args...)
}

func Error(msg string, args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	instance.Printf("[ERROR] "+msg, args...)
}
