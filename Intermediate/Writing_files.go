package main

import (
	"fmt"
	"os"
)

func WritingFiles() {

	// ============================================================
	// 🔹 Creating File
	// ============================================================

	// os.Create():
	// → creates new file
	// → truncates file if already exists
	// → returns *os.File

	file, err := os.Create("Output.txt")

	if err != nil {
		panic("Error Creating File")
	}


	// ============================================================
	// 🔹 defer Keyword
	// ============================================================

	// defer delays function execution
	// until surrounding function finishes

	// Here:
	// file.Close() will execute automatically
	// when main() ends

	// Used mainly for cleanup:
	// - closing files
	// - DB cleanup
	// - network cleanup
	// - flushing buffers

	// IMPORTANT:
	// defer executes even if:
	// - return occurs
	// - panic occurs

	defer file.Close()


	// ============================================================
	// 🔹 Writing Bytes to File
	// ============================================================

	data := []byte("Hello World\n")

	// Write():
	// → writes byte slice to file

	// returns:
	// n   → bytes written
	// err → possible error

	n, err := file.Write(data)

	if err != nil {
		panic("Cannot Write To File")
	}

	fmt.Println("Bytes Written:", n)

	fmt.Println(
		"Data Written Successfully...",
	)


	// ============================================================
	// 🔹 WriteString()
	// ============================================================

	newFile, err := os.Create("WriteString.txt")

	if err != nil {
		panic("Error Creating File")
	}

	// deferred function for second file
	defer newFile.Close()


	// WriteString():
	// → directly writes string to file
	// → internally handles byte conversion

	m, err := newFile.WriteString(
		"Hello World\n",
	)

	if err != nil {
		panic("Error Writing String")
	}

	fmt.Println(
		"Bytes Written Using WriteString:",
		m,
	)

	fmt.Println(
		"String Written Successfully...",
	)
}


// ============================================================
// 🔹 DIFFERENCE: Write vs WriteString
// ============================================================

// Write()
// → writes []byte

// Example:
// file.Write([]byte("Hello"))


// WriteString()
// → writes string directly

// Example:
// file.WriteString("Hello")


// ============================================================
// 🔹 QUICK REVISION
// ============================================================

// os.Create()
// → create/truncate file

// *os.File
// → file handler object

// defer
// → executes when surrounding function finishes

// Close()
// → releases file resource

// Write()
// → write byte slice

// WriteString()
// → write string directly

// Always close files
// to avoid resource leaks