package main

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Zap() {

	// ============================================================
	// 🔹 WHAT IS ZAP?
	// ============================================================

	// Zap is a high-performance structured logging library
	// developed by Uber.
	//
	// It is designed for:
	// - Production applications
	// - Microservices
	// - High-throughput systems
	// - Distributed systems
	//
	// Advantages:
	// - Extremely fast
	// - Low memory allocations
	// - Structured logging
	// - JSON output
	// - Production ready


	// ============================================================
	// 🔹 PRODUCTION LOGGER
	// ============================================================

	// NewProduction() creates a production-ready logger.
	//
	// Default features:
	// - JSON output
	// - Info level and above
	// - Unix timestamp
	// - Stack traces for severe errors
	//
	// Example output:
	//
	// {
	//   "level":"info",
	//   "ts":1748580000.123,
	//   "msg":"User Logged In"
	// }

	logger, err := zap.NewProduction()


	// ============================================================
	// 🔹 ERROR HANDLING
	// ============================================================

	// Logger creation can fail.
	//
	// Therefore error handling is required.

	if err != nil {
		log.Println("Error initializing zap logger")
		return
	}


	// ============================================================
	// 🔹 CLEANUP
	// ============================================================

	// Sync() flushes buffered log entries.
	//
	// Some logs may still be waiting in memory.
	// Sync ensures everything gets written.
	//
	// defer ensures Sync runs when main exits.

	defer logger.Sync()


	// ============================================================
	// 🔹 BASIC LOGGING
	// ============================================================

	// Info()
	// Used for normal application events.

	logger.Info("This is an info message")


	// ============================================================
	// 🔹 STRUCTURED LOGGING
	// ============================================================

	// Structured logging stores data as fields.
	//
	// Instead of:
	//
	// "User John Doe logged in via GET"
	//
	// We store:
	//
	// username = John Doe
	// method   = GET
	//
	// This makes logs:
	// - Searchable
	// - Filterable
	// - Easy to aggregate

	logger.Info(
		"User Logged In",

		// String field
		zap.String("username", "John Doe"),

		// Another string field
		zap.String("method", "GET"),
	)


	// ============================================================
	// 🔹 COMMON FIELD TYPES
	// ============================================================

	// zap.String(key, value)
	// zap.Int(key, value)
	// zap.Bool(key, value)
	// zap.Float64(key, value)
	// zap.Duration(key, value)
	// zap.Error(err)

	// Example:
	//
	// logger.Info(
	//     "Example",
	//     zap.String("name", "Burhaan"),
	//     zap.Int("age", 19),
	//     zap.Bool("active", true),
	// )


	// ============================================================
	// 🔹 LOG LEVELS
	// ============================================================

	// Debug()
	// → Development information

	// Info()
	// → Normal application events

	// Warn()
	// → Recoverable issues

	// Error()
	// → Failures

	// Panic()
	// → Logs message and calls panic()

	// Fatal()
	// → Logs message and exits application


	// ============================================================
	// 🔹 CUSTOM PRODUCTION CONFIGURATION
	// ============================================================

	// NewProductionConfig()
	//
	// Unlike NewProduction(),
	// this returns a configurable blueprint.
	//
	// We can customize:
	// - Timestamp format
	// - Timestamp key name
	// - Log level
	// - Output paths
	// - Encoders

	config := zap.NewProductionConfig()


	// ============================================================
	// 🔹 CUSTOM TIMESTAMP KEY
	// ============================================================

	// Default timestamp key:
	//
	// "ts"
	//
	// We change it to:
	//
	// "timestamp"

	config.EncoderConfig.TimeKey = "timestamp"


	// ============================================================
	// 🔹 CUSTOM TIMESTAMP FORMAT
	// ============================================================

	// Default Production Logger:
	//
	// "ts":1748580000.123
	//
	// Which is Unix Epoch time.
	//
	// ISO8601 produces:
	//
	// "2026-05-30T15:30:45.123+0530"
	//
	// Human-readable and commonly used.

	config.EncoderConfig.EncodeTime =
		zapcore.ISO8601TimeEncoder


	// ============================================================
	// 🔹 BUILD CUSTOM LOGGER
	// ============================================================

	// Build() creates a logger
	// from the configuration.

	newLogger, err := config.Build()

	if err != nil {
		log.Println("Error building custom logger")
		return
	}


	// Flush custom logger before exit.

	defer newLogger.Sync()


	// ============================================================
	// 🔹 CUSTOM LOGGER OUTPUT
	// ============================================================

	newLogger.Info(
		"New Log Created",

		// Username field
		zap.String("username", "Burhaan"),

		// Age field
		zap.Int("Age", 19),
	)

	// Example output:
	//
	// {
	//   "level":"info",
	//   "timestamp":"2026-05-30T15:30:45.123+0530",
	//   "msg":"New Log Created",
	//   "username":"Burhaan",
	//   "Age":19
	// }
}


// ============================================================
// 🔹 QUICK REVISION
// ============================================================

// zap.NewProduction()
// → Creates a ready-made production logger

// zap.NewProductionConfig()
// → Returns configurable logger settings

// config.Build()
// → Creates logger from configuration

// logger.Sync()
// → Flushes buffered logs

// logger.Info()
// → Information log

// logger.Warn()
// → Warning log

// logger.Error()
// → Error log

// logger.Fatal()
// → Log + Exit application

// logger.Panic()
// → Log + Panic

// zap.String()
// → Creates string field

// zap.Int()
// → Creates integer field

// zap.Bool()
// → Creates boolean field

// zap.Float64()
// → Creates float field

// zap.Error()
// → Creates error field

// TimeKey
// → Name of timestamp field

// EncodeTime
// → Controls timestamp format

// ISO8601TimeEncoder
// → Human-readable timestamp format

// Structured Logging
// → Logs contain searchable key-value fields

// Zap
// → Fast, low-allocation, production-grade logger


// ============================================================
// 🔹 INTERVIEW NOTES
// ============================================================

// Q: Why use Zap instead of fmt.Println()?
//
// A:
// - Structured logging
// - Log levels
// - JSON output
// - Better performance
// - Production ready


// Q: Difference between NewProduction()
// and NewProductionConfig()?
//
// NewProduction()
// → Ready-made logger
//
// NewProductionConfig()
// → Configurable blueprint


// Q: Why structured logging?
//
// Because logs become:
// - Machine-readable
// - Searchable
// - Filterable
// - Easy to aggregate


// Q: Why call Sync()?
//
// To ensure buffered logs are written
// before the application exits.