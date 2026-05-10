package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func lineFilters() {

	// ============================================================
	// 🔹 Opening File
	// ============================================================

	// os.Open()
	// → opens file in read-only mode

	file, err := os.Open(
		"Intermediate/example.txt",
	)

	if err != nil {
		panic("Cannot Open File")
	}


	// ============================================================
	// 🔹 defer Cleanup
	// ============================================================

	// defer executes when surrounding function ends

	defer func() {

		fmt.Println("Closing Open File")

		file.Close()

	}()


	// ============================================================
	// 🔹 Scanner
	// ============================================================

	// Scanner:
	// → reads file line by line

	// Useful for:
	// - text processing
	// - log parsing
	// - searching/filtering lines

	scanner := bufio.NewScanner(file)


	// ============================================================
	// 🔹 Keyword Searching
	// ============================================================

	keyword := "important"

	for scanner.Scan() {

		// Text()
		// → current scanned line

		line := scanner.Text()


		// ========================================================
		// 🔹 Searching Text
		// ========================================================

		// strings.Contains()
		// → checks substring presence

		if strings.Contains(line, keyword) {


			// ====================================================
			// 🔹 Replacing Text
			// ====================================================

			// ReplaceAll()
			// → replaces all matching substrings

			updatedLine := strings.ReplaceAll(
				line,
				keyword,
				"necessary",
			)

			fmt.Println(
				"Filtered Line:",
				line,
			)

			fmt.Println(
				"Updated Line:",
				updatedLine,
			)
		}
	}


	// ============================================================
	// 🔹 Scanner Errors
	// ============================================================

	// Err()
	// → returns scanner error if occurred

	err = scanner.Err()

	if err != nil {
		panic("Error Scanning File")
	}


	// ============================================================
	// 🔹 Alternative Approach (Read Whole File)
	// ============================================================

	// os.ReadFile()
	// → loads entire file into memory

	// Better for:
	// - small files
	// - simple replacements

	data, err := os.ReadFile(
		"Intermediate/example.txt",
	)

	if err != nil {
		panic("Cannot Read Entire File")
	}

	content := string(data)

	updatedContent := strings.ReplaceAll(
		content,
		"important",
		"necessary",
	)

	fmt.Println("\nWhole File Replacement:")
	fmt.Println(updatedContent)


	// Scanner approach is better for:
	// - large files
	// - streaming
	// - memory efficiency
	// - line-by-line processing
}


// ============================================================
// 🔹 QUICK REVISION
// ============================================================

// os.Open()
// → open file

// Scanner
// → line-by-line reading

// scanner.Scan()
// → move to next line

// scanner.Text()
// → current line

// strings.Contains()
// → check substring existence

// strings.ReplaceAll()
// → replace all matches

// os.ReadFile()
// → reads entire file into memory

// defer
// → cleanup when function ends