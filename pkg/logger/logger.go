//go:generate go run github.com/dmarkham/enumer -type=Level -output level_enumer___generated.go
package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

// Level represents the severity of a log message.
// It can be one of DEBUG, INFO, WARN, ERROR, ALWAYS, or SILENT.
type Level int

const (
	// SILENT represents no logging, effectively muting all log messages.
	SILENT Level = iota - 1
	// DEBUG represents debug-level messages, useful for development and troubleshooting.
	DEBUG
	// INFO represents informational messages, typically used for normal operation.
	INFO
	// WARN represents warning messages, which indicate potential issues but not failures.
	WARN
	// ERROR represents error messages, indicating failure in operation.
	ERROR
	// ALWAYS represents messages that should always be logged, regardless of the current log level.
	ALWAYS
)

// Logger holds the configuration for logging.
// It includes the current logging level, the output writer, and color mappings for each log level.
type Logger struct {
	level  Level                  // level is the current logging level. Messages with lower severity will be ignored.
	output io.Writer              // output is the writer where log messages will be written (e.g., stdout, a file).
	colors map[Level]*color.Color // colors holds color settings for each log Level to make log messages more readable.
}

// NewCustom creates a new Logger instance with the specified log level and output writer.
// If an invalid log level is provided, it defaults to INFO.
func NewCustom(level Level, output io.Writer) *Logger {
	if !level.IsALevel() {
		fmt.Fprintf(os.Stderr, "Invalid log level: %q, setting to %q\n", level, INFO)
		level = INFO
	}

	return &Logger{
		level:  level,
		output: output,
		colors: map[Level]*color.Color{
			DEBUG:  color.New(color.FgBlue),
			INFO:   color.New(color.FgGreen),
			WARN:   color.New(color.FgYellow),
			ERROR:  color.New(color.FgRed),
			ALWAYS: color.New(color.FgGreen),
		},
	}
}

// New creates a new Logger instance with the specified log level and writes to stdout.
// If an invalid log level is provided, it defaults to INFO.
func New(level Level) *Logger {
	if !level.IsALevel() {
		fmt.Fprintf(os.Stderr, "Invalid log level: %q, setting to %q\n", level, INFO)
		level = INFO
	}

	return &Logger{
		level:  level,
		output: os.Stdout,
		colors: map[Level]*color.Color{
			DEBUG:  color.New(color.FgBlue),
			INFO:   color.New(color.FgGreen),
			WARN:   color.New(color.FgYellow),
			ERROR:  color.New(color.FgRed),
			ALWAYS: color.New(color.FgGreen),
		},
	}
}

// log prints a log message with the specified log level, applying colors based on the level if available.
// The message will only be logged if the level's severity is equal to or higher than the Logger's current level.
func (l *Logger) log(level Level, format string, args ...any) {
	if l.level == SILENT {
		return
	}

	if level >= l.level {
		message := fmt.Sprintf(format, args...)
		if c, ok := l.colors[level]; ok {
			c.Fprintln(l.output, message)
		} else {
			fmt.Fprintln(l.output, message)
		}
	}
}

// Always logs a message at the INFO level, regardless of the current log level.
func (l *Logger) Always(format string, args ...any) {
	l.log(ALWAYS, format, args...)
}

// Debug logs a debug-level message if the current log level allows it.
func (l *Logger) Debug(format string, args ...any) {
	l.log(DEBUG, format, args...)
}

// Info logs an informational message if the current log level allows it.
func (l *Logger) Info(format string, args ...any) {
	l.log(INFO, format, args...)
}

// Warn logs a warning message if the current log level allows it.
func (l *Logger) Warn(format string, args ...any) {
	l.log(WARN, format, args...)
}

// Error logs an error message if the current log level allows it.
func (l *Logger) Error(format string, args ...any) {
	l.log(ERROR, format, args...)
}
