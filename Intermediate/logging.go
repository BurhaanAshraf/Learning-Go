package main

import (
	"log"
	"os"
)
func logging() {

	// ============================================================
	// 🔹 Default Logger
	// ============================================================

	// log.Println()
	// → uses Go's default logger

	log.Println("This is a log message")

	log.SetPrefix("INFO: ")
	log.Println("This is an info message")


	// ============================================================
	// 🔹 Logger Flags
	// ============================================================

	// Flags add metadata to logs

	// log.Ldate
	// → current date

	// log.Ltime
	// → current time

	// log.Lshortfile
	// → filename and line number

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("This is a log message with date and time.")


	// ============================================================
	// 🔹 Custom Loggers
	// ============================================================

	// Custom loggers allow different:
	// - prefixes
	// - flags
	// - output destinations

	infoLogger.Println("This is an INFO message")
	warnLogger.Println("This is a Warning message")
	errorLogger.Println("This is an Error message")


	// ============================================================
	// 🔹 File Logging
	// ============================================================

	// Logs can also be written to files

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Ensure file is closed when main exits

	defer file.Close()


	// ============================================================
	// 🔹 Loggers With File Output
	// ============================================================

	infoLogger1 := log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	warnLogger1 := log.New(file, "WARN: ", log.Ldate|log.Ltime)

	errorLogger1 := log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	debugLogger := log.New(file, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

	debugLogger.Println("This is a debug message")
	warnLogger1.Println("This is a warning message")
	infoLogger1.Println("This is an info message")
	errorLogger1.Println("This is an error message")
}

var (

	// log.New()
	// → returns *log.Logger

	// Each logger can have:
	// - separate prefix
	// - separate flags
	// - separate output destination

	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	warnLogger = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime)

	errorLogger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)


// ============================================================
// 🔹 QUICK REVISION
// ============================================================

// log.Println()
// → default logger

// log.SetPrefix()
// → add custom prefix

// log.SetFlags()
// → configure log metadata

// log.New()
// → create custom logger

// log.Ldate
// → current date

// log.Ltime
// → current time

// log.Lshortfile
// → filename + line number

// os.Stdout
// → terminal output

// os.OpenFile()
// → open/create log file

// os.O_APPEND
// → append logs instead of overwrite

// defer file.Close()
// → cleanup file resource

// INFO
// → normal events

// WARN
// → recoverable issues

// ERROR
// → failures

// DEBUG
// → development details