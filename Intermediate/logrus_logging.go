package main

import "github.com/sirupsen/logrus"

func Logrus() {

// ============================================================
// 🔹 Logrus
// ============================================================

// Logrus
// → popular third-party logging library for Go

// Provides:
// - log levels
// - structured logging
// - JSON output
// - fields/metadata

// More powerful than standard log package


log := logrus.New()


// ============================================================
// 🔹 Log Levels
// ============================================================

// Log levels help categorize logs

// Trace
// → most detailed logs

// Debug
// → development information

// Info
// → normal application events

// Warn
// → recoverable issues

// Error
// → failures requiring attention

// Fatal
// → log + os.Exit(1)

// Panic
// → log + panic

log.SetLevel(logrus.InfoLevel)


// Only logs with level >= Info
// will be displayed


// ============================================================
// 🔹 Formatters
// ============================================================

// Formatter controls log output format

// TextFormatter
// → human readable output

// JSONFormatter
// → structured machine-readable output

// JSON logs are commonly used in:
// - APIs
// - microservices
// - production systems

log.SetFormatter(&logrus.JSONFormatter{})


// ============================================================
// 🔹 Basic Logging
// ============================================================

// Info()
// → normal events

// Warn()
// → warnings

// Error()
// → errors

log.Info("This is an info message")
log.Warn("This is a warning message")
log.Error("This is an error message")


// ============================================================
// 🔹 Structured Logging
// ============================================================

// WithFields()
// → attaches metadata to log entry

// Metadata makes logs easier to:
// - search
// - filter
// - analyze

// Fields become part of JSON output

log.WithFields(logrus.Fields{
	"username": "John Doe",
	"method": "GET",
}).Info("User logged in")

}

// WithFields()
// → attaches metadata to log entry

// Info()/Warn()/Error()
// → actually write the log

// WithFields() alone does not produce output


// ============================================================
// 🔹 Why Structured Logs?
// ============================================================

// Instead of:
//
// User John logged in using GET

// We store:
//
// username = John
// method = GET

// This allows log systems to
// filter and aggregate data easily.