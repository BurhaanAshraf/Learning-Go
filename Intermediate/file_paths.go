package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func filePath() {

	// ============================================================
	// 🔹 Relative vs Absolute Paths
	// ============================================================

	// Relative Path:
	// → path relative to current working directory
	// → depends on where program is running

	relativePath := "./data/file.txt"

	// Absolute Path:
	// → complete exact path from system root

	absolutePath := "/home/user/docs/file.txt"

	fmt.Println("Relative Path:", relativePath)
	fmt.Println("Absolute Path:", absolutePath)


	// ============================================================
	// 🔹 filepath.Join()
	// ============================================================

	// Join():
	// → safely combines multiple path parts
	// → automatically uses OS-specific separator

	// Linux/macOS:
	// /

	// Windows:
	// \

	// Better than manually doing:
	// "home/user/downloads/file.zip"

	joinedPath := filepath.Join(
		"home",
		"user",
		"downloads",
		"file.zip",
	)

	fmt.Println("Joined Path:", joinedPath)


	// ============================================================
	// 🔹 filepath.Clean()
	// ============================================================

	// Clean():
	// → normalizes/simplifies path

	// Removes:
	// - duplicate separators
	// - unnecessary "."
	// - resolves ".."

	// "."  → current directory
	// ".." → parent directory

	normalizedPath := filepath.Clean(
		"./data/../data/file.txt",
	)

	fmt.Println(
		"Normalized Path:",
		normalizedPath,
	)


	// ============================================================
	// 🔹 filepath.Split()
	// ============================================================

	// Split():
	// → splits path into:
	// directory + filename

	dir, file := filepath.Split(
		"/home/user/docs/downloads/file.txt",
	)

	fmt.Println("Directory:", dir)
	fmt.Println("File:", file)


	// ============================================================
	// 🔹 filepath.Base()
	// ============================================================

	// Base():
	// → returns last part of path

	// Useful when:
	// - extracting filename
	// - ignoring directories

	base := filepath.Base(
		"/home/user/docs/downloads/file.txt",
	)

	fmt.Println("Base:", base)


	// ============================================================
	// 🔹 filepath.IsAbs()
	// ============================================================

	// IsAbs():
	// → checks whether path is absolute

	// Relative paths:
	// depend on current working directory

	// Absolute paths:
	// always point to same location

	fmt.Println(
		"Is Relative Path Absolute?",
		filepath.IsAbs(relativePath),
	)

	fmt.Println(
		"Is Absolute Path Absolute?",
		filepath.IsAbs(absolutePath),
	)


	// ============================================================
	// 🔹 filepath.Ext()
	// ============================================================

	// Ext():
	// → extracts file extension

	extension := filepath.Ext(file)

	fmt.Println("Extension:", extension)


	// ============================================================
	// 🔹 Removing File Extension
	// ============================================================

	// TrimSuffix():
	// → removes matching suffix from string

	fileNameWithoutExt := strings.TrimSuffix(
		file,
		extension,
	)

	fmt.Println(
		"File Without Extension:",
		fileNameWithoutExt,
	)


	// ============================================================
	// 🔹 filepath.Rel()
	// ============================================================

	// Rel(base,target)
	// → calculates relative path
	// from base → target

	relativeResult, err := filepath.Rel(
		"a/c",
		"a/b/t/file",
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(
		"Relative Path Result:",
		relativeResult,
	)

	// Means:
	// "how do I go from a/c to a/b/t/file?"


	// ============================================================
	// 🔹 filepath.Abs()
	// ============================================================

	// Abs():
	// → converts relative path
	// → into absolute path

	absoluteResult, err := filepath.Abs(
		relativePath,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(
		"Absolute Path:",
		absoluteResult,
	)
}


// ============================================================
// 🔹 QUICK REVISION
// ============================================================

// filepath.Join()
// → safely join path parts

// filepath.Clean()
// → normalize/simplify path

// filepath.Split()
// → directory + filename

// filepath.Base()
// → last path element

// filepath.IsAbs()
// → check if path is absolute

// filepath.Ext()
// → get file extension

// strings.TrimSuffix()
// → remove extension

// filepath.Rel()
// → relative path between locations

// filepath.Abs()
// → relative → absolute path