package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func Directories() {

	// ============================================================
	// 🔹 Root Working Directory
	// ============================================================

	// temporary root directory for learning/demo

	rootDir := "temp_parent"


	// ============================================================
	// 🔹 Directory Permissions
	// ============================================================

	// 0755:
	//
	// owner:
	// read + write + execute
	//
	// others:
	// read + execute

	// execute permission on directories
	// allows traversal/access


	// ============================================================
	// 🔹 Creating Nested Directories
	// ============================================================

	// MkdirAll()
	// → recursively creates nested directories

	checkErr(
		os.MkdirAll(
			filepath.Join(
				rootDir,
				"child",
				"grandchild1",
			),
			0755,
		),
	)

	checkErr(
		os.MkdirAll(
			filepath.Join(
				rootDir,
				"child",
				"grandchild2",
			),
			0755,
		),
	)

	checkErr(
		os.MkdirAll(
			filepath.Join(
				rootDir,
				"child",
				"grandchild3",
			),
			0755,
		),
	)


	// ============================================================
	// 🔹 Creating Files
	// ============================================================

	// WriteFile()
	// → creates/truncates/writes file

	createFile(
		filepath.Join(
			rootDir,
			"child",
			"grandchild1",
			"file1.txt",
		),
	)

	createFile(
		filepath.Join(
			rootDir,
			"child",
			"grandchild2",
			"file2.txt",
		),
	)

	createFile(
		filepath.Join(
			rootDir,
			"child",
			"grandchild3",
			"file3.txt",
		),
	)

	createFile(
		filepath.Join(
			rootDir,
			"child",
			"childFile.txt",
		),
	)

	createFile(
		filepath.Join(
			rootDir,
			"parentFile.txt",
		),
	)


	// ============================================================
	// 🔹 Reading Directory Entries
	// ============================================================

	entries, err := os.ReadDir(
		filepath.Join(rootDir, "child"),
	)

	checkErr(err)

	fmt.Println("Directory Entries:")

	for _, entry := range entries {

		fmt.Println(
			"Name:",
			entry.Name(),

			"| IsDir:",
			entry.IsDir(),
		)
	}


	// ============================================================
	// 🔹 Current Working Directory
	// ============================================================

	// "."  → current directory
	// ".." → parent directory


	// ============================================================
	// 🔹 Changing Directory
	// ============================================================

	checkErr(
		os.Chdir(
			filepath.Join(
				rootDir,
				"child",
				"grandchild3",
			),
		),
	)

	fmt.Println("\nCurrent Directory Entries:")

	entries, err = os.ReadDir(".")

	checkErr(err)

	for _, entry := range entries {

		fmt.Println(entry.Name())
	}


	// ============================================================
	// 🔹 Moving Backward
	// ============================================================

	checkErr(os.Chdir("../.."))


	// ============================================================
	// 🔹 Current Working Directory Path
	// ============================================================

	currentDir, err := os.Getwd()

	checkErr(err)

	fmt.Println(
		"\nCurrent Working Directory:",
		currentDir,
	)


	// ============================================================
	// 🔹 filepath.Join()
	// ============================================================

	joinedPath := filepath.Join(
		rootDir,
		"child",
		"grandchild1",
		"file1.txt",
	)

	fmt.Println("\nJoined Path:", joinedPath)


	// ============================================================
	// 🔹 filepath.Clean()
	// ============================================================

	cleanPath := filepath.Clean(
		"./temp_parent/../temp_parent/child",
	)

	fmt.Println("Clean Path:", cleanPath)


	// ============================================================
	// 🔹 filepath.Abs()
	// ============================================================

	absolutePath, err := filepath.Abs(
		"./temp_parent",
	)

	checkErr(err)

	fmt.Println("Absolute Path:", absolutePath)


	// ============================================================
	// 🔹 filepath.Rel()
	// ============================================================

	relativePath, err := filepath.Rel(
		rootDir,
		filepath.Join(
			rootDir,
			"child",
			"grandchild1",
		),
	)

	checkErr(err)

	fmt.Println("Relative Path:", relativePath)


	// ============================================================
	// 🔹 filepath.WalkDir()
	// ============================================================

	// WalkDir():
	// → recursive traversal using DirEntry
	// → newer + more efficient

	fmt.Println("\nWalkDir Traversal:")

	err = filepath.WalkDir(

		"../temp_parent",

		func(
			path string,
			d os.DirEntry,
			err error,
		) error {

			checkErr(err)

			fmt.Println(path)

			return nil
		},
	)

	checkErr(err)


	// ============================================================
	// 🔹 filepath.Walk()
	// ============================================================

	// Walk():
	// → older traversal using FileInfo

	fmt.Println("\nWalk Traversal:")

	err = filepath.Walk(

		"../temp_parent",

		func(
			path string,
			info os.FileInfo,
			err error,
		) error {

			checkErr(err)

			fmt.Println(
				path,
				"| IsDir:",
				info.IsDir(),
			)

			return nil
		},
	)

	checkErr(err)


	// ============================================================
	// 🔹 Cleanup
	// ============================================================

	// RemoveAll()
	// → recursively deletes everything

	checkErr(
		os.RemoveAll("../temp_parent"),
	)

	fmt.Println("\nTemporary Directory Removed")
}


// ============================================================
// 🔹 Helper Function
// ============================================================

func createFile(path string) {

	err := os.WriteFile(
		path,
		[]byte(""),
		0755,
	)

	checkErr(err)
}


// ============================================================
// 🔹 Error Helper
// ============================================================

func checkErr(err error) {

	if err != nil {
		panic(err)
	}
}


// ============================================================
// 🔹 QUICK REVISION
// ============================================================

// MkdirAll()
// → recursive directory creation

// WriteFile()
// → create/write/truncate file

// ReadDir()
// → read directory entries

// Chdir()
// → change current directory

// Getwd()
// → current working directory

// filepath.Join()
// → safely join paths

// filepath.Clean()
// → normalize path

// filepath.Abs()
// → absolute path

// filepath.Rel()
// → relative path

// filepath.WalkDir()
// → newer recursive traversal

// filepath.Walk()
// → older traversal using FileInfo

// RemoveAll()
// → recursive delete

// "."  → current directory
// ".." → parent directory