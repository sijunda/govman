// Package logger provides a centralized logging system for user-facing messages
package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	viper "github.com/spf13/viper"
)

// LogLevel represents the verbosity level of log messages
type LogLevel int

const (
	// QuietLevel shows only errors
	QuietLevel LogLevel = iota
	// NormalLevel shows essential information
	NormalLevel
	// VerboseLevel shows detailed information
	VerboseLevel
)

// Logger handles all user-facing messages
type Logger struct {
	level  LogLevel
	writer io.Writer
	mutex  sync.Mutex
}

// Timer represents a timer for measuring operation duration
type Timer struct {
	start time.Time
	name  string
}

// New creates a new Logger instance
func New() *Logger {
	l := &Logger{
		writer: os.Stderr,
	}

	// Set log level based on flags
	if viper.GetBool("quiet") {
		l.level = QuietLevel
	} else if viper.GetBool("verbose") {
		l.level = VerboseLevel
	} else {
		l.level = NormalLevel
	}

	return l
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level LogLevel) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.level = level
}

// Error logs an error message (always shown)
func (l *Logger) Error(format string, args ...interface{}) {
	if l.level >= QuietLevel {
		fmt.Fprintf(l.writer, "❌ Error: "+format+"\n", args...)
	}
}

// ErrorWithHelp logs an error message with additional help text
func (l *Logger) ErrorWithHelp(errorMsg, helpMsg string, args ...interface{}) {
	if l.level >= QuietLevel {
		fmt.Fprintf(l.writer, "❌ Error: "+errorMsg+"\n", args...)
		if helpMsg != "" {
			fmt.Fprintf(l.writer, "💡 Help: %s\n", helpMsg)
		}
	}
}

// StartTimer starts a new timer for measuring operation duration
func (l *Logger) StartTimer(name string) *Timer {
	if l.level >= VerboseLevel {
		l.Verbose("Starting %s...", name)
	}
	return &Timer{
		start: time.Now(),
		name:  name,
	}
}

// StopTimer stops the timer and logs the duration
func (l *Logger) StopTimer(t *Timer) {
	if l.level >= VerboseLevel && t != nil {
		duration := time.Since(t.start)
		l.Verbose("Completed %s in %v", t.name, duration)
	}
}

// Info logs an informational message
func (l *Logger) Info(format string, args ...interface{}) {
	if l.level >= NormalLevel {
		fmt.Fprintf(l.writer, "ℹ️  "+format+"\n", args...)
	}
}

// Success logs a success message
func (l *Logger) Success(format string, args ...interface{}) {
	if l.level >= NormalLevel {
		fmt.Fprintf(l.writer, "✅ "+format+"\n", args...)
	}
}

// Warning logs a warning message
func (l *Logger) Warning(format string, args ...interface{}) {
	if l.level >= NormalLevel {
		fmt.Fprintf(l.writer, "⚠️  "+format+"\n", args...)
	}
}

// Verbose logs a verbose message (only shown in verbose mode)
func (l *Logger) Verbose(format string, args ...interface{}) {
	if l.level >= VerboseLevel {
		fmt.Fprintf(l.writer, "🔍 [VERBOSE] "+format+"\n", args...)
	}
}

// Debug logs a debug message (only shown in verbose mode)
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.level >= VerboseLevel {
		fmt.Fprintf(l.writer, "🐛 [DEBUG] "+format+"\n", args...)
	}
}

// Progress logs a progress message
func (l *Logger) Progress(format string, args ...interface{}) {
	if l.level >= NormalLevel {
		fmt.Fprintf(l.writer, "🔄 "+format+"\n", args...)
	}
}

// Step logs a step in a process
func (l *Logger) Step(format string, args ...interface{}) {
	if l.level >= VerboseLevel {
		fmt.Fprintf(l.writer, "📋 "+format+"\n", args...)
	}
}

// Download logs a download message
func (l *Logger) Download(format string, args ...interface{}) {
	if l.level >= NormalLevel {
		fmt.Fprintf(l.writer, "📦 "+format+"\n", args...)
	}
}

// Extract logs an extraction message
func (l *Logger) Extract(format string, args ...interface{}) {
	if l.level >= NormalLevel {
		fmt.Fprintf(l.writer, "📂 "+format+"\n", args...)
	}
}

// Verify logs a verification message
func (l *Logger) Verify(format string, args ...interface{}) {
	if l.level >= NormalLevel {
		fmt.Fprintf(l.writer, "🔍 "+format+"\n", args...)
	}
}

// globalLogger is the shared logger instance
var globalLogger *Logger
var once sync.Once

// Get returns the global logger instance
func Get() *Logger {
	once.Do(func() {
		globalLogger = New()
	})
	return globalLogger
}

// Error logs an error using the global logger
func Error(format string, args ...interface{}) {
	Get().Error(format, args...)
}

// Info logs an info message using the global logger
func Info(format string, args ...interface{}) {
	Get().Info(format, args...)
}

// Success logs a success message using the global logger
func Success(format string, args ...interface{}) {
	Get().Success(format, args...)
}

// Warning logs a warning message using the global logger
func Warning(format string, args ...interface{}) {
	Get().Warning(format, args...)
}

// Verbose logs a verbose message using the global logger
func Verbose(format string, args ...interface{}) {
	Get().Verbose(format, args...)
}

// Debug logs a debug message using the global logger
func Debug(format string, args ...interface{}) {
	Get().Debug(format, args...)
}

// Progress logs a progress message using the global logger
func Progress(format string, args ...interface{}) {
	Get().Progress(format, args...)
}

// Download logs a download message using the global logger
func Download(format string, args ...interface{}) {
	Get().Download(format, args...)
}

// Extract logs an extraction message using the global logger
func Extract(format string, args ...interface{}) {
	Get().Extract(format, args...)
}

// Verify logs a verification message using the global logger
func Verify(format string, args ...interface{}) {
	Get().Verify(format, args...)
}

// StartTimer starts a new timer using the global logger
func StartTimer(name string) *Timer {
	return Get().StartTimer(name)
}

// StopTimer stops a timer using the global logger
func StopTimer(t *Timer) {
	Get().StopTimer(t)
}

// ErrorWithHelp logs an error with help text using the global logger
func ErrorWithHelp(errorMsg, helpMsg string, args ...interface{}) {
	Get().ErrorWithHelp(errorMsg, helpMsg, args...)
}

// Step logs a step using the global logger
func Step(format string, args ...interface{}) {
	Get().Step(format, args...)
}
