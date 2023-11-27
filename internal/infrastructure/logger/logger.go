package logger

import (
	"fmt"
	"os"

	"golang.org/x/exp/slog"
)

type Attribute int

const (
	FgBlack Attribute = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

// Logger is a simple logging struct
type Logger struct {
	// You can customize the logger by adding more fields here
	Prefix string
	Output *os.File
  log *slog.Logger
}

// NewLogger creates a new Logger instance with the given prefix and output file
func NewLogger(prefix string, output *os.File) *Logger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return &Logger{
		Prefix: prefix,
		Output: output,
    log: logger,
	}
}

// Log prints a log message with a timestamp and the specified prefix
func (l *Logger) Info(message string) {
	// logEntry := fmt.Sprintf("[%s] %s: %s\n", time.Now().Format("2006-01-02 15:04:05"), l.Prefix, message)
	colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", 34, message)
	slog.Info(colored)
}

func (l *Logger) Error(err string) {
	// errorMessage := fmt.Sprintf("[%s] %s ERROR: %v\n", time.Now().Format("2006-01-02 15:04:05"), l.Prefix, err)
	colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", 31, err)
	slog.Error(colored)
}

func (l *Logger) Warning(message string) {
	// message := fmt.Sprintf("[%s] %s ERROR: %v\n", time.Now().Format("2006-01-02 15:04:05"), l.Prefix, err)
	colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", FgYellow, message)
	slog.Warn(colored)
}
