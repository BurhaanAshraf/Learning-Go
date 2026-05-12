package main

import (
	"fmt"
	"os"
)

func tempFilesAndDirs() {

	// ============================================================
	// 🔹 Temporary Files
	// ============================================================

	// CreateTemp(dir, pattern)
	//
	// "" → system temporary directory
	//
	// pattern:
	// filename prefix

	tempFile, err := os.CreateTemp(
		"",
		"temporaryFile",
	)

	checkError(err)

	fmt.Println(
		"Temporary File Created:",
		tempFile.Name(),
	)


	// ============================================================
	// 🔹 defer Execution Order
	// ============================================================

	// defer works in:
	// LIFO order
	// (Last In First Out)

	// So:
	// last defer added
	// runs first


	// ============================================================
	// 🔹 Cleanup Order
	// ============================================================

	// We want:
	// 1. close file
	// 2. remove file

	// So Remove() is deferred first
	// and Close() second

	defer func() {

		checkError(
			os.Remove(tempFile.Name()),
		)

		fmt.Println(
			"Temporary File Removed",
		)

	}()

	defer func() {

		checkError(tempFile.Close())

		fmt.Println(
			"Temporary File Closed",
		)

	}()


	// ============================================================
	// 🔹 Temporary Directory
	// ============================================================

	// MkdirTemp()
	// → creates unique temporary directory

	tempDir, err := os.MkdirTemp(
		"",
		"temporaryFolder",
	)

	checkError(err)

	fmt.Println(
		"Temporary Directory Created:",
		tempDir,
	)


	// ============================================================
	// 🔹 RemoveAll()
	// ============================================================

	// RemoveAll()
	// → recursively deletes directory
	// and everything inside it

	defer func() {

		checkError(
			os.RemoveAll(tempDir),
		)

		fmt.Println(
			"Temporary Directory Removed",
		)

	}()
}


// ============================================================
// 🔹 Error Helper
// ============================================================

func checkError(err error) {

	if err != nil {
		panic(err)
	}
}


// ============================================================
// 🔹 QUICK REVISION
// ============================================================

// CreateTemp()
// → create temporary file

// MkdirTemp()
// → create temporary directory

// Remove()
// → remove single file

// RemoveAll()
// → recursive delete

// defer
// → delayed execution

// defer order
// → LIFO (Last In First Out)

// Better cleanup order:
// close resource → then remove resource