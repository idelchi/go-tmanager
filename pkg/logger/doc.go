// Package logger provides a simple logging framework that supports
// different log levels (DEBUG, INFO, WARN, ERROR, ALWAYS, SILENT) and color-coded output.
//
// The logger is configurable to print messages at a specified level or higher,
// and it allows for custom output writers. Each log level has a corresponding
// color to make messages easier to distinguish in the output.
//
// Example usage:
//
//	Create a new logger that logs to stdout and has an INFO level.
//	l := logger.New(logger.INFO)
//
//	// Log messages at various levels.
//	l.Debug("This is a debug message")  // Won't be shown since log level is INFO.
//	l.Info("This is an info message")   // Will be shown.
//	l.Warn("This is a warning message") // Will be shown.
//	l.Error("This is an error message") // Will be shown.
//	l.Always("This is an always message") // Will be shown.
//
// You can also create a custom logger with a different output:
//
//	// Create a logger that writes to a file and logs at the DEBUG level.
//	file, _ := os.Create("app.log")
//	l := logger.NewCustom(logger.DEBUG, file)
package logger
