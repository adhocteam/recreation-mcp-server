package util

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// LogLevel represents the severity level of a log message
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// String returns the string representation of a log level
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger provides structured logging with configurable levels
type Logger struct {
	level  LogLevel
	output io.Writer
	format string
}

// NewLogger creates a new Logger instance
func NewLogger(levelStr, format string) *Logger {
	level := parseLogLevel(levelStr)
	return &Logger{
		level:  level,
		output: os.Stdout,
		format: format,
	}
}

// parseLogLevel converts a string to a LogLevel
func parseLogLevel(levelStr string) LogLevel {
	switch strings.ToLower(levelStr) {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn":
		return WARN
	case "error":
		return ERROR
	default:
		return INFO
	}
}

// SetOutput sets the output destination for the logger
func (l *Logger) SetOutput(w io.Writer) {
	l.output = w
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(levelStr string) {
	l.level = parseLogLevel(levelStr)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.log(DEBUG, msg, fields...)
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...interface{}) {
	l.log(INFO, msg, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.log(WARN, msg, fields...)
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...interface{}) {
	l.log(ERROR, msg, fields...)
}

// log writes a log message if it meets the minimum level
func (l *Logger) log(level LogLevel, msg string, fields ...interface{}) {
	if level < l.level {
		return
	}

	timestamp := time.Now().Format(time.RFC3339)

	if l.format == "json" {
		l.logJSON(level, timestamp, msg, fields...)
	} else {
		l.logText(level, timestamp, msg, fields...)
	}
}

// logJSON outputs log messages in JSON format
func (l *Logger) logJSON(level LogLevel, timestamp, msg string, fields ...interface{}) {
	output := fmt.Sprintf(`{"timestamp":"%s","level":"%s","message":"%s"`,
		timestamp, level.String(), msg)

	if len(fields) > 0 {
		output += `,"fields":{`
		for i := 0; i < len(fields); i += 2 {
			if i > 0 {
				output += ","
			}
			key := fmt.Sprintf("%v", fields[i])
			value := fmt.Sprintf("%v", fields[i+1])
			output += fmt.Sprintf(`"%s":"%s"`, key, value)
		}
		output += "}"
	}

	output += "}\n"
	log.Print(output)
}

// logText outputs log messages in plain text format
func (l *Logger) logText(level LogLevel, timestamp, msg string, fields ...interface{}) {
	output := fmt.Sprintf("[%s] %s: %s", timestamp, level.String(), msg)

	if len(fields) > 0 {
		output += " |"
		for i := 0; i < len(fields); i += 2 {
			key := fmt.Sprintf("%v", fields[i])
			value := fmt.Sprintf("%v", fields[i+1])
			output += fmt.Sprintf(" %s=%s", key, value)
		}
	}

	output += "\n"
	log.Print(output)
}

// Debugf logs a debug message with formatting
func (l *Logger) Debugf(format string, args ...interface{}) {
	if DEBUG >= l.level {
		l.Debug(fmt.Sprintf(format, args...))
	}
}

// Infof logs an info message with formatting
func (l *Logger) Infof(format string, args ...interface{}) {
	if INFO >= l.level {
		l.Info(fmt.Sprintf(format, args...))
	}
}

// Warnf logs a warning message with formatting
func (l *Logger) Warnf(format string, args ...interface{}) {
	if WARN >= l.level {
		l.Warn(fmt.Sprintf(format, args...))
	}
}

// Errorf logs an error message with formatting
func (l *Logger) Errorf(format string, args ...interface{}) {
	if ERROR >= l.level {
		l.Error(fmt.Sprintf(format, args...))
	}
}
