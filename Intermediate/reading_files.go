package main

import (
	"bufio"
	"fmt"
	"os"
)

func readingFiles() {

	// ============================================================
	// 🔹 Opening File
	// ============================================================

	// os.Open()
	// → opens existing file in read-only mode
	// → returns *os.File

	file, err := os.Open(
		"/media/burhaan/Shared/Learning-GO/Intermediate/Output.txt",
	)

	if err != nil {
		panic("Cannot Open File")
	}

	fmt.Println("File Opened Successfully")


	// ============================================================
	// 🔹 defer Cleanup
	// ============================================================

	// defer executes when surrounding function finishes

	// Used for:
	// - cleanup
	// - closing files
	// - releasing resources

	defer func() {

		fmt.Println("Closing Open File")

		file.Close()

	}()


	// ============================================================
	// 🔹 Reading Whole File Using Read()
	// ============================================================

	// Read():
	// → reads bytes into slice

	// Create byte buffer of size 1024
	data := make([]byte, 1024)

	// n = actual bytes read
	n, err := file.Read(data)

	if err != nil {
		panic("Cannot Read File")
	}

	// data[:n]
	// → only valid read bytes

	fmt.Println(
		"File Content:\n",
		string(data[:n]),
	)


	// ============================================================
	// 🔹 IMPORTANT FILE POINTER CONCEPT
	// ============================================================

	// Read() moves the internal file pointer forward

	// So after:
	// file.Read(data)

	// the pointer may already be at EOF
	// (End Of File)

	// That is why scanner may not read anything
	// after full file read


	// ============================================================
	// 🔹 Seek()
	// ============================================================

	// Seek(offset, whence)
	// → moves internal file pointer

	// offset:
	// → how many bytes to move

	// whence:
	// 0 → beginning of file
	// 1 → current position
	// 2 → end of file

	// Seek(0,0)
	// → move 0 bytes from beginning
	// → resets pointer to start of file

	_, err = file.Seek(0, 0)

	if err != nil {
		panic("Cannot Reset File Pointer")
	}


	// ============================================================
	// 🔹 Scanner (Line By Line Reading)
	// ============================================================

	// Scanner is useful for:
	// - line-by-line reading
	// - text processing
	// - CLI/file input

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		// Text() returns current line
		line := scanner.Text()

		fmt.Println("Line:", line)
	}

	// Check scanner errors
	err = scanner.Err()

	if err != nil {
		panic("Error Reading File")
	}
}


// ============================================================
// 🔹 QUICK REVISION
// ============================================================

// os.Open()
// → open file in read-only mode

// Read()
// → read bytes into slice

// data[:n]
// → valid bytes only

// Scanner
// → line-by-line reading

// scanner.Scan()
// → move to next line

// scanner.Text()
// → current line

// Seek(offset, whence)
// → move file pointer

// whence:
// 0 → beginning
// 1 → current position
// 2 → end

// Seek(0,0)
// → reset pointer to beginning

// defer
// → cleanup when function finishes