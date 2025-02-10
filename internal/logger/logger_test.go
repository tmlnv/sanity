package logger

import (
	"os"
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name    string
		logFile string
		wantErr bool
	}{
		{
			name:    "init with stdout only",
			logFile: "",
			wantErr: false,
		},
		{
			name:    "init with valid log file",
			logFile: "test.log",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				Close()
				if tt.logFile != "" {
					os.Remove(tt.logFile)
				}
			}()

			Init(tt.logFile)

			// Verify logger is initialized
			if instance == nil {
				t.Error("logger instance is nil after initialization")
			}

			// Test logging
			Info("test message")
		})
	}
}

func TestClose(t *testing.T) {
	tests := []struct {
		name     string
		setupLog bool
		wantErr  bool
	}{
		{
			name:     "close initialized logger",
			setupLog: true,
			wantErr:  false,
		},
		{
			name:     "close uninitialized logger",
			setupLog: false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupLog {
				Init("test.log")
			}

			err := Close()
			if (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Cleanup
			if tt.setupLog {
				os.Remove("test.log")
			}
		})
	}
}

func Test_getGoroutineID(t *testing.T) {
	tests := []struct {
		name string
		want uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getGoroutineID(); got != tt.want {
				t.Errorf("getGoroutineID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInfo(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		args    []interface{}
		wantErr bool
	}{
		{
			name:    "log info message without args",
			msg:     "test message",
			args:    nil,
			wantErr: false,
		},
		{
			name:    "log info message with args",
			msg:     "test message",
			args:    []interface{}{"key1", "value1", "key2", 42},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			tmpFile := "test.log"
			Init(tmpFile)
			defer func() {
				Close()
				os.Remove(tmpFile)
			}()

			// Test logging
			Info(tt.msg, tt.args...)

			// Verify log file contains the message
			content, err := os.ReadFile(tmpFile)
			if err != nil {
				t.Fatalf("Failed to read log file: %v", err)
			}
			if !strings.Contains(string(content), tt.msg) {
				t.Errorf("Log file does not contain message: %s", tt.msg)
			}
		})
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		args    []interface{}
		wantErr bool
	}{
		{
			name:    "log error message without args",
			msg:     "error message",
			args:    nil,
			wantErr: false,
		},
		{
			name:    "log error message with args",
			msg:     "error message",
			args:    []interface{}{"error", "test error"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			tmpFile := "test.log"
			Init(tmpFile)
			defer func() {
				Close()
				os.Remove(tmpFile)
			}()

			// Test logging
			Error(tt.msg, tt.args...)

			// Verify log file contains the message
			content, err := os.ReadFile(tmpFile)
			if err != nil {
				t.Fatalf("Failed to read log file: %v", err)
			}
			if !strings.Contains(string(content), tt.msg) {
				t.Errorf("Log file does not contain message: %s", tt.msg)
			}
		})
	}
}

func TestDebug(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		args    []interface{}
		wantErr bool
	}{
		{
			name:    "log debug message without args",
			msg:     "debug message",
			args:    nil,
			wantErr: false,
		},
		{
			name:    "log debug message with args",
			msg:     "debug message",
			args:    []interface{}{"debug", "test debug", "level", "verbose"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			tmpFile := "test.log"
			Init(tmpFile)
			defer func() {
				Close()
				os.Remove(tmpFile)
			}()

			// Test logging
			Debug(tt.msg, tt.args...)

			// Verify log file contains the message and goroutine ID
			content, err := os.ReadFile(tmpFile)
			if err != nil {
				t.Fatalf("Failed to read log file: %v", err)
			}
			if !strings.Contains(string(content), tt.msg) {
				t.Errorf("Log file does not contain message: %s", tt.msg)
			}
			if !strings.Contains(string(content), "[GoroutineID:") {
				t.Error("Log file does not contain goroutine ID")
			}
		})
	}
}
