package logger

import (
	"fmt"
	"os"
	"time"
)

// Logger is a simple logging struct
type Logger struct {
	// You can customize the logger by adding more fields here
	Prefix string
	Output *os.File
}

// NewLogger creates a new Logger instance with the given prefix and output file
func NewLogger(prefix string, output *os.File) *Logger {
	return &Logger{
		Prefix: prefix,
		Output: output,
	}
}

// Log prints a log message with a timestamp and the specified prefix
func (l *Logger) Info(message string) {
	logEntry := fmt.Sprintf("[%s] %s: %s\n", time.Now().Format("2006-01-02 15:04:05"), l.Prefix, message)

	if l.Output != nil {
		_, _ = l.Output.WriteString(logEntry)
	} else {
		fmt.Print(logEntry)
	}
}
